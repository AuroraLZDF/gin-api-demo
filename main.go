package main

import "github.com/auroraLZDF/gin-api-demo/config"
import "github.com/auroraLZDF/gin-api-demo/router"

func init()  {
	config.SetConfig()
	config.SetDb()
}

func main()  {
	//当整个程序完成之后关闭数据库连接
	defer config.Db.Close()

	router := router.InitRouter()
	router.Run(":8080")
}


