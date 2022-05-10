package router

import (
	"github.com/Dlimingliang/shop-api/goods-web/api"
	"github.com/gin-gonic/gin"
)

func InitGoodsRouter(group *gin.RouterGroup) {
	goodsGroup := group.Group("goods")
	{
		goodsGroup.GET("page", api.GetGoodsPage)
	}
}
