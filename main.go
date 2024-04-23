package main

import (
	"fmt"

	"github.com/jericho-yu/gorm-pool/gormPool"
)

func main() {
	dbSetting := gormPool.NewDbSetting("./settings")

	//  创建mysql连接池
	mysqlPool := gormPool.NewMySqlPool(dbSetting)

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
