package router

import (
	"github.com/gin-gonic/gin"

	"github.com/Dlimingliang/shop-api/user-web/api"
)

func InitUserRouter(group *gin.RouterGroup) {
	userGroup := group.Group("user")
	{
		userGroup.GET("list", api.GetUserList)
	}
}
