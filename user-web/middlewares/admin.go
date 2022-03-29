package middlewares

import (
	"github.com/Dlimingliang/shop-api/user-web/global"
	"github.com/Dlimingliang/shop-api/user-web/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func IsAdmin() gin.HandlerFunc {
	return func(context *gin.Context) {
		claims, _ := context.Get(global.JWTGinContextKey)
		currentUser := claims.(*models.CustomClaims)
		if currentUser.RoleId != 1 {
			context.JSON(http.StatusForbidden, gin.H{"msg": "无权限"})
			context.Abort()
			return
		}
		context.Next()
	}
}
