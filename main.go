package main

import (
	"fmt"

	"github.com/jericho-yu/gorm-pool/gormPool"
)

func main() {
	var e error
	//  创建mysql连接池
	mysqlPool := gormPool.NewMySqlPool(gormPool.NewDbSetting("./settings/db.yaml"))

	// 获取mysql链接
	mysqlConn := mysqlPool.GetConn()
	fmt.Println("mysqlConn:", mysqlConn)

	// 创建postgres链接
	postgresPool := gormPool.NewPostgresPool(gormPool.NewDbSetting("./settings/db.yaml"))

	// 获取postgres链接
	postgresConn := postgresPool.GetConn()
	fmt.Println("postgresConn:", postgresConn)

	// 关闭数据库链接
	defer func() {
		e = mysqlPool.Close()
		if e != nil {
			panic(e)
		}
		e = postgresPool.Close()
		if e != nil {
			panic(e)
		}
	}()
}
