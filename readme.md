```yaml
common:
  driver: "mysql"

mysql:
  maxOpenConns: 100
  maxIdleConns: 20
  maxLifetime: 100
  maxIdleTime: 10
  rws: true
  main:
    username: "root"
    password: "root"
    host: 127.0.0.1
    port: 3308
    database: "abc_passport"
    charset: "utf8mb4"
    collation: "utf8mb4_general_ci"
  sources:
    conn1:
      username: "root"
      password: "root"
      host: 127.0.0.1
      port: 3308
      database: "abc_passport"
      charset: "utf8mb4"
      collation: "utf8mb4_general_ci"
    conn2:
      username: "root"
      password: "root"
      host: 127.0.0.1
      port: 3308
      database: "abc_passport"
      charset: "utf8mb4"
      collation: "utf8mb4_general_ci"
  replicas:
    conn3:
      username: "root"
      password: "root"
      host: 127.0.0.1
      port: 3308
      database: "abc_passport"
      charset: "utf8mb4"
      collation: "utf8mb4_general_ci"
    conn4:
      username: "root"
      password: "root"
      host: 127.0.0.1
      port: 3308
      database: "abc_passport"
      charset: "utf8mb4"
      collation: "utf8mb4_general_ci"

postgres:
  maxOpenConns: 100
  maxIdleConns: 20
  maxLifetime: 100
  maxIdleTime: 10
  main:
    username: "postgres"
    password: "postgres"
    host: 127.0.0.1
    port: 5432
    database: "abc_passport"
    sslmode: "disable"
    timezone: "Asia/Shanghai"

sqlServer:
  maxOpenConns: 100
  maxIdleConns: 20
  maxLifetime: 100
  maxIdleTime: 10
  main:
    username: "admin"
    password: "Admin@1234"
    host: 127.0.0.1
    port: 9930
    database: "abc_passport"
```
```go
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

	// 创建postgres链接池
	postgresPool := gormPool.PostgresPoolApp.New(gormPool.NewDbSetting("./settings/db.yaml"))

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

```