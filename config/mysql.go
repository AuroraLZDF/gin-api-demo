package config

import (
	_ "github.com/go-sql-driver/mysql" //加载mysql
	"github.com/jinzhu/gorm"
	"fmt"
	"time"
)

var Db *gorm.DB

//初始化方法
func SetDb() {
	var err error
	//var dbConfig = DbConfig()

	format := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local&timeout=10ms",
		Config.DbUser, Config.DbPass, Config.DbHost, Config.DbPort, Config.DbName) + "&loc=Asia%2FShanghai"

	Db, err = gorm.Open("mysql", format)
	if err != nil {
		panic("mysql connect error: " + err.Error())
	}

	/*if env := os.Getenv("APP_ENV"); env != "production" {

	}*/

	Db.DB().SetMaxIdleConns(20)  // 数据库的空闲连接
	Db.DB().SetMaxOpenConns(20) // 数据库的最大连接
	Db.DB().SetConnMaxLifetime(60 * time.Second)	// 防止出现 (invalid connection) 数据库连接超时错误
}
