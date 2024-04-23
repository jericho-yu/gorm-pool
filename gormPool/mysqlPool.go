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

// NewMySqlPool 创建mysql链接池对象
func NewMySqlPool(dbSetting *DbSetting) *MySqlPool {
	mysqlPoolOnce.Do(func() {
		mysqlPoolIns = &MySqlPool{
			username:     dbSetting.MySql.Main.Username,
			password:     dbSetting.MySql.Main.Password,
			host:         dbSetting.MySql.Main.Host,
			port:         dbSetting.MySql.Main.Port,
			database:     dbSetting.MySql.Main.Database,
			charset:      dbSetting.MySql.Main.Charset,
			sources:      dbSetting.MySql.Sources,
			replicas:     dbSetting.MySql.Replicas,
			rws:          false,
			maxIdleTime:  dbSetting.MySql.MaxIdleTime,
			maxLifetime:  dbSetting.MySql.MaxLifetime,
			maxIdleConns: dbSetting.MySql.MaxIdleConns,
			maxOpenConns: dbSetting.MySql.MaxOpenConns,
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
			dbSetting.MySql.Main.Username,
			dbSetting.MySql.Main.Password,
			dbSetting.MySql.Main.Host,
			dbSetting.MySql.Main.Port,
			dbSetting.MySql.Main.Database,
			dbSetting.MySql.Main.Charset,
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
	sqlDb.SetConnMaxIdleTime(time.Duration(mysqlPoolIns.maxIdleTime) * time.Hour)
	sqlDb.SetConnMaxLifetime(time.Duration(mysqlPoolIns.maxLifetime) * time.Hour)
	sqlDb.SetMaxIdleConns(mysqlPoolIns.maxIdleConns)
	sqlDb.SetMaxOpenConns(mysqlPoolIns.maxOpenConns)

	return mysqlPoolIns
}

// GetMain 获取主数据库链接
func (receiver *MySqlPool) GetConn() *gorm.DB {
	receiver.getRws()
	return receiver.mainConn
}

// getRws 获取带有读写分离的数据库链接
func (receiver *MySqlPool) getRws() *gorm.DB {
	var (
		err                                 error
		sourceDialectors, replicaDialectors []gorm.Dialector
		sources                             []*Dsn
		replicas                            []*Dsn
	)
	// 配置写库
	if len(receiver.sources) > 0 {
		sources = make([]*Dsn, 0)
		for idx, item := range receiver.sources {
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
	if len(receiver.replicas) > 0 {
		replicas = make([]*Dsn, 0)
		for idx, item := range receiver.replicas {
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
