package router

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/auroraLZDF/gin-api-demo/app/Controllers/store"
	"os"
	"io"
	"github.com/auroraLZDF/gin-api-demo/app/Controllers/goods"
	"github.com/auroraLZDF/gin-api-demo/app/Controllers/moldata"
	"github.com/auroraLZDF/gin-api-demo/app/Controllers/purchase"
	"github.com/auroraLZDF/gin-api-demo/app/Controllers/supply"
)

func InitRouter() *gin.Engine {
	// 写入日志的文件
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "member api",
		})
	})

	api := router.Group("/collection")
	{
		// 店铺
		api.GET("store", store.Index)
		api.GET("store/types", store.Types)

		// 商品
		api.GET("/goods", goods.Index)
		api.GET("/goods/brands", goods.Brands)


		// 百科
		api.GET("/moldata", moldata.Index)
		api.GET("/moldata/categories", moldata.Categories)
		//api.GET("/moldata/:userId/status", moldata.HasCollection)	// todo: 这个路由报错

		// 采购单（询盘）
		api.GET("/inquiry", purchase.Index)
		api.GET("/inquiry/:userId/status", purchase.HasCollection)

		// 供应单
		api.GET("/supply", supply.Index)
		api.GET("/supply/:userId/status", supply.HasCollection)

	}

	// middleware
	api = router.Group("/collection" /*, middleware.Handler()*/)
	{

		// 店铺
		api.POST("store/:storeId", store.Create)
		api.POST("store/:storeId/cancel", store.Cancel)
		api.DELETE("store/batch_cancel", store.BatchCancel)

		// 商品
		api.POST("goods/:goodsId", goods.Create)
		api.POST("goods/:goodsId/cancel", goods.Cancel)
		api.DELETE("goods/empty", goods.EmptyGoods)
		api.DELETE("goods/batch_cancel", goods.BatchCancel)

		// 百科
		api.POST("moldata/:molId", moldata.Create)
		api.POST("moldata/:molId/cancel", moldata.Cancel)
		api.DELETE("moldata/batch_cancel", moldata.BatchCancel)

		// 采购（询盘）
		api.POST("inquiry/:code", purchase.Create)
		api.DELETE("inquiry/:userId/batchCancel", purchase.BranchCancel)

		// 供应
		api.POST("supply/:supplyId", supply.Create)
		api.DELETE("supply/:userId/batchCancel", supply.BranchCancel)

	}

	return router
}
