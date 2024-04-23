package main

import (
	"fmt"

	"github.com/jericho-yu/gorm-pool/gormPool"
)

func main() {
	//  创建mysql连接池
	mysqlPool := gormPool.NewMySqlPool(gormPool.NewDbSetting("./settings"))

	// 获取数据库链接
	mysqlSingle := mysqlPool.GetConn()
	fmt.Println(mysqlSingle)

	// 关闭数据库链接
	defer func() {
		e := mysqlPool.Close()
		if e != nil {
			panic(e)
		}
	}()
}
