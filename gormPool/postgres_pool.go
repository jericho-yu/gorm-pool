package gormPool

import (
	"fmt"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresPool struct {
	username     string
	password     string
	host         string
	port         uint16
	database     string
	maxIdleTime  int
	maxLifetime  int
	maxIdleConns int
	maxOpenConns int
	mainDsn      *Dsn
	mainConn     *gorm.DB
}

var (
	postgresPoolIns   *PostgresPool
	postgresPoolOnce  sync.Once
	PostgresDsnFormat = "host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s"
)

// NewPostgresPool 创建mysql链接池对象
func NewPostgresPool(dbSetting *DbSetting) *PostgresPool {
	postgresPoolOnce.Do(func() {
		postgresPoolIns = &PostgresPool{
			username:     dbSetting.Postgres.Main.Username,
			password:     dbSetting.Postgres.Main.Password,
			host:         dbSetting.Postgres.Main.Host,
			port:         dbSetting.Postgres.Main.Port,
			database:     dbSetting.Postgres.Main.Database,
			maxIdleTime:  dbSetting.Postgres.MaxIdleTime,
			maxLifetime:  dbSetting.Postgres.MaxLifetime,
			maxIdleConns: dbSetting.Postgres.MaxIdleConns,
			maxOpenConns: dbSetting.Postgres.MaxOpenConns,
		}
	})

	var (
		err      error
		dbConfig *gorm.Config
	)

	// 配置主库
	postgresPoolIns.mainDsn = &Dsn{
		Name: "main",
		Content: fmt.Sprintf(
			PostgresDsnFormat,
			dbSetting.Postgres.Main.Host,
			dbSetting.Postgres.Main.Username,
			dbSetting.Postgres.Main.Password,
			dbSetting.Postgres.Main.Database,
			dbSetting.Postgres.Main.Port,
			dbSetting.Postgres.Main.SslMode,
			dbSetting.Postgres.Main.TimeZone,
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
	postgresPoolIns.mainConn, err = gorm.Open(postgres.Open(postgresPoolIns.mainDsn.Content), dbConfig)
	if err != nil {
		panic(fmt.Sprintf("配置数据库失败：%s", err.Error()))
	}

	postgresPoolIns.mainConn = postgresPoolIns.mainConn.Session(&gorm.Session{})

	sqlDb, _ := postgresPoolIns.mainConn.DB()
	sqlDb.SetConnMaxIdleTime(time.Duration(postgresPoolIns.maxIdleTime) * time.Hour)
	sqlDb.SetConnMaxLifetime(time.Duration(postgresPoolIns.maxLifetime) * time.Hour)
	sqlDb.SetMaxIdleConns(postgresPoolIns.maxIdleConns)
	sqlDb.SetMaxOpenConns(postgresPoolIns.maxOpenConns)

	return postgresPoolIns
}

// GetMain 获取主数据库链接
func (receiver *PostgresPool) GetConn() *gorm.DB {
	return receiver.mainConn
}

// Close 关闭数据库链接
func (receiver *PostgresPool) Close() error {
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
