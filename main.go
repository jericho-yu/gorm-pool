package main

import (
	"fmt"

	"github.com/jericho-yu/gorm-pool/gormPool"
)

func main() {
	dbSetting := NewDbSetting("./settings")

	//  创建mysql连接池
	mysqlPool := gormPool.NewMySqlPool(
		dbSetting.MySql.Main.Username,
		dbSetting.MySql.Main.Password,
		dbSetting.MySql.Main.Host,
		dbSetting.MySql.Main.Port,
		dbSetting.MySql.Main.Database,
		dbSetting.MySql.Main.Charset,
		dbSetting.MySql.MaxIdleTime,
		dbSetting.MySql.MaxLifetime,
		dbSetting.MySql.MaxIdleConns,
		dbSetting.MySql.MaxOpenConns,
	)

	// 创建单数据库链接
	mysqlSingle := mysqlPool.GetMain()
	fmt.Println(mysqlSingle)

	// 创建读写分离数据库
	mysqlRws := mysqlPool.GetRws(
		dbSetting.MySql.Sources,
		dbSetting.MySql.Replicas,
	)
	fmt.Println(mysqlRws)

	// 关闭数据库链接
	defer func() {
		e := mysqlPool.Close()
		if e != nil {
			panic(e)
		}
	}()
}
