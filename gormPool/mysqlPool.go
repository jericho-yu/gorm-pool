package gormPool

import (
	"fmt"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

type MySqlPool struct {
	username     string
	password     string
	host         string
	port         uint16
	database     string
	charset      string
	sources      map[string]*MySqlConnection
	replicas     map[string]*MySqlConnection
	rws          bool
	maxIdleTime  int
	maxLifetime  int
	maxIdleConns int
	maxOpenConns int
	mainDsn      *Dsn
	mainConn     *gorm.DB
}

var (
	mysqlPoolIns  *MySqlPool
	mysqlPoolOnce sync.Once
)

// GetMain 获取主数据库链接
func (receiver *MySqlPool) GetMain() *gorm.DB {
	return receiver.mainConn
}

// NewMySqlPool 创建mysql链接池对象
func NewMySqlPool(
	Username,
	Password,
	Host string,
	Port uint16,
	Database,
	Charset string,
	MaxIdleTime,
	MaxLifetime,
	MaxIdleConns,
	MaxOpenConns int,
) *MySqlPool {
	mysqlPoolOnce.Do(func() {
		mysqlPoolIns = &MySqlPool{
			username:     Username,
			password:     Password,
			host:         Host,
			port:         Port,
			database:     Database,
			charset:      Charset,
			sources:      make(map[string]*MySqlConnection),
			replicas:     make(map[string]*MySqlConnection),
			rws:          false,
			maxIdleTime:  MaxIdleTime,
			maxLifetime:  MaxLifetime,
			maxIdleConns: MaxIdleConns,
			maxOpenConns: MaxOpenConns,
		}
	})

	var (
		err      error
		dbConfig *gorm.Config
	)

	// 配置主库
	mysqlPoolIns.mainDsn = &Dsn{
		Name: "main",
		Content: fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			Username,
			Password,
			Host,
			Port,
			Database,
			Charset,
		),
	}

	// 数据库配置
	dbConfig = &gorm.Config{
		PrepareStmt:                              true,  // 预编译
		CreateBatchSize:                          500,   // 批量操作
		DisableForeignKeyConstraintWhenMigrating: true,  // 禁止自动创建外键
		SkipDefaultTransaction:                   false, // 开启自动事务
		QueryFields:                              true,  // 查询字段
		AllowGlobalUpdate:                        false, // 不允许全局修改,必须带有条件
	}

	// 配置主库
	mysqlPoolIns.mainConn, err = gorm.Open(mysql.Open(mysqlPoolIns.mainDsn.Content), dbConfig)
	if err != nil {
		panic(fmt.Sprintf("配置主库失败：%s", err.Error()))
	}

	mysqlPoolIns.mainConn = mysqlPoolIns.mainConn.Session(&gorm.Session{})

	sqlDb, _ := mysqlPoolIns.mainConn.DB()
	sqlDb.SetConnMaxIdleTime(time.Duration(MaxIdleTime) * time.Hour)
	sqlDb.SetConnMaxLifetime(time.Duration(MaxLifetime) * time.Hour)
	sqlDb.SetMaxIdleConns(MaxIdleConns)
	sqlDb.SetMaxOpenConns(MaxOpenConns)

	return mysqlPoolIns
}

// GetRws 获取带有读写分离的数据库链接
func (receiver *MySqlPool) GetRws(Sources map[string]*MySqlConnection, Replicas map[string]*MySqlConnection) *gorm.DB {
	var (
		err                                 error
		sourceDialectors, replicaDialectors []gorm.Dialector
		sources                             []*Dsn
		replicas                            []*Dsn
	)
	// 配置写库
	if len(Sources) > 0 {
		sources = make([]*Dsn, 0)
		for idx, item := range Sources {
			sources = append(sources, &Dsn{
				Name: idx,
				Content: fmt.Sprintf(
					"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
					item.Username,
					item.Password,
					item.Host,
					item.Port,
					item.Database,
					item.Charset,
				),
			})
		}
	}

	// 配置读库
	if len(Replicas) > 0 {
		replicas = make([]*Dsn, 0)
		for idx, item := range Replicas {
			replicas = append(replicas, &Dsn{
				Name: idx,
				Content: fmt.Sprintf(
					"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
					item.Username,
					item.Password,
					item.Host,
					item.Port,
					item.Database,
					item.Charset,
				),
			})
		}
	}

	if len(sources) > 0 {
		sourceDialectors = make([]gorm.Dialector, len(sources))
		for i := 0; i < len(sources); i++ {
			sourceDialectors[i] = mysql.Open(sources[i].Content)
		}
	}

	if len(replicas) > 0 {
		replicaDialectors = make([]gorm.Dialector, len(replicas))
		for i := 0; i < len(replicas); i++ {
			replicaDialectors[i] = mysql.Open(replicas[i].Content)
		}
	}

	err = receiver.mainConn.Use(
		dbresolver.Register(dbresolver.Config{
			Sources:           sourceDialectors,          // 写库
			Replicas:          replicaDialectors,         // 读库
			Policy:            dbresolver.RandomPolicy{}, // 策略
			TraceResolverMode: true,
		}).
			SetConnMaxIdleTime(time.Duration(receiver.maxIdleTime) * time.Hour).
			SetConnMaxLifetime(time.Duration(receiver.maxLifetime) * time.Hour).
			SetMaxIdleConns(receiver.maxIdleConns).
			SetMaxOpenConns(receiver.maxOpenConns),
	)
	if err != nil {
		panic(fmt.Errorf("数据库链接错误：%s", err.Error()))
	}

	return receiver.mainConn
}

// Close 关闭数据库链接
func (receiver *MySqlPool) Close() error {
	if receiver.mainConn != nil {
		db, err := receiver.mainConn.DB()
		if err != nil {
			return fmt.Errorf("关闭数据库链接失败：获取数据库链接失败 %s", err.Error())
		}
		err = db.Close()
		if err != nil {
			return fmt.Errorf("关闭数据库连接失败 %s", err.Error())
		}
	}
	return nil
}
