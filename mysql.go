package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var myDb *sql.DB
var hasDbInit bool

const (
	username = "root"
	password = "root"
	network  = "tcp"
	server   = "localhost"
	port     = 3306
	database = "tuliao"
)

func initDb() {
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", username, password, network, server, port, database)
	var err error
	myDb, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("open mysql failed ,err:%v\n", err)
		if myDb != nil {
			myDb.Close()
			myDb = nil
		}
		hasDbInit = false
		return
	}
	myDb.SetConnMaxLifetime(100 * time.Second) //最大连接周期，超过时间连接就关闭
	myDb.SetMaxOpenConns(100)                  //最大连接数
	myDb.SetMaxIdleConns(16)                   //设置闲置连接数

	hasDbInit = true

}
