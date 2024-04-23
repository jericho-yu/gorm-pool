package gormPool

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"time"
)

type MySqlPool struct{}

func NewMySqlPool(
	Username,
	Password,
	Host string,
	Port uint16,
	Database,
	Charset string,
	Sources map[string]*MySqlConnection,
	Replicas map[string]*MySqlConnection,
	Rws bool,
	MaxIdleTime,
	MaxLifetime,
	MaxIdleConns,
	MaxOpenConns int,
) *gorm.DB {
	var (
		err                                 error
		dbConfig                            *gorm.Config
		db                                  *gorm.DB
		sourceDialectors, replicaDialectors []gorm.Dialector
		main                                *Dsn
		sources                             []*Dsn
		replicas                            []*Dsn
	)

	// 配置主库
	main = &Dsn{
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

	// 数据库配置
	dbConfig = &gorm.Config{
		PrepareStmt:                              true,  // 预编译
		CreateBatchSize:                          500,   // 批量操作
		DisableForeignKeyConstraintWhenMigrating: true,  // 禁止自动创建外键
		SkipDefaultTransaction:                   false, // 开启自动事务
		QueryFields:                              true,  // 查询字段
		AllowGlobalUpdate:                        false, // 不允许全局修改,必须带有条件
	}

	if Rws {
		// 配置主库
		db, err = gorm.Open(mysql.Open(main.Content), dbConfig)
		if err != nil {
			panic(fmt.Errorf("配置主库失败：%s", err.Error()))
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

		err = db.Use(
			dbresolver.Register(dbresolver.Config{
				Sources:           sourceDialectors,          // 写库
				Replicas:          replicaDialectors,         // 读库
				Policy:            dbresolver.RandomPolicy{}, // 策略
				TraceResolverMode: true,
			}).
				SetConnMaxIdleTime(time.Duration(MaxIdleTime) * time.Hour).
				SetConnMaxLifetime(time.Duration(MaxLifetime) * time.Hour).
				SetMaxIdleConns(MaxIdleConns).
				SetMaxOpenConns(MaxOpenConns),
		)
		if err != nil {
			panic(fmt.Errorf("数据库链接错误：%s", err.Error()))
		}
	} else {
		// 配置主库
		db, err = gorm.Open(mysql.Open(main.Content), dbConfig)
		if err != nil {
			panic(fmt.Sprintf("配置主库失败：%s", err.Error()))
		}

		db = db.Session(&gorm.Session{})

		sqlDb, _ := db.DB()
		sqlDb.SetConnMaxIdleTime(time.Duration(MaxIdleTime) * time.Hour)
		sqlDb.SetConnMaxLifetime(time.Duration(MaxLifetime) * time.Hour)
		sqlDb.SetMaxIdleConns(MaxIdleConns)
		sqlDb.SetMaxOpenConns(MaxOpenConns)
	}

	return db
}
