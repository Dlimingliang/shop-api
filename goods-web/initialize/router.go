package initialize

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/Dlimingliang/shop-api/goods-web/router"
)

func Routers() *gin.Engine {
	ginRouter := gin.Default()
	//设置跨域
	ginRouter.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
	}))
	//设置路由
	//ginRouter.POST("register", api.Register)
	apiGroup := ginRouter.Group("v1")
	//测试接口，暂时不限制登录
	//apiGroup.Use(middlewares.JWTAuth()).Use(middlewares.IsAdmin())
	router.InitGoodsRouter(apiGroup)
	return ginRouter
}
