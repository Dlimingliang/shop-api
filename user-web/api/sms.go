package api

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"

	"github.com/Dlimingliang/shop-api/user-web/forms"
	"github.com/Dlimingliang/shop-api/user-web/global"
)

func GenerateSmsCode(width int) string {
	numbers := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	n := len(numbers)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numbers[rand.Intn(n)])
	}
	return sb.String()
}

func SendSms(ctx *gin.Context) {

	sendSmsForm := forms.SendSmsForm{}
	if err := ctx.ShouldBind(&sendSmsForm); err != nil {
		HandlerValidatorErr(err, ctx)
		return
	}

	//模拟验证码发送
	smsCode := GenerateSmsCode(4)
	zap.S().Info("验证码已发送,验证码为: ", smsCode)

	//将发送的验证码保存起来
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.ServerConfig.RedisConfig.Host, global.ServerConfig.RedisConfig.Port),
		Password: global.ServerConfig.RedisConfig.Password,
	})

	rdb.Set(context.Background(), sendSmsForm.Mobile, smsCode, 300*time.Second)
	ctx.JSON(http.StatusOK, gin.H{"msg": "验证码发送成功"})
}
