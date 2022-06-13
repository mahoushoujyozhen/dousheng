package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var dbConn *sql.DB

func initDB() (err error) {
	dsn := "root:huangzhen123@tcp(175.178.106.176:3306)/Publish?charset=utf8mb4"
	// open函数只是验证格式是否正确，并不是创建数据库连接
	dbConn, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	//dbConn.SetMaxOpenConns(10)
	//dbConn.SetMaxIdleConns(5)
	//dbConn.SetConnMaxLifetime(time.Minute * 60)
	// 与数据库建立连接
	err = dbConn.Ping()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	//var err error
	r := gin.Default()

	//mysqlDB connect mysql database
	//err = initDB()
	//if err != nil {
	//	fmt.Printf("err:%v\n", err)
	//	fmt.Println("mysql DB connect failed!!")
	//	return
	//}

	fmt.Println("mysql DB connect success!")

	initRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
