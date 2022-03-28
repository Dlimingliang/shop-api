package initialize

import (
	"github.com/Dlimingliang/shop-api/user-web/middlewares"
	"github.com/Dlimingliang/shop-api/user-web/router"
	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	ginRouter := gin.Default()
	apiGroup := ginRouter.Group("v1")
	apiGroup.Use(middlewares.JWTAuth())
	router.InitUserRouter(apiGroup)
	return ginRouter
}
