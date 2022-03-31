package api

import (
	"fmt"
	"go.uber.org/zap"
	"math/rand"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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

	//模拟验证码发送
	smsCode := GenerateSmsCode(4)
	zap.S().Info("验证码已发送,验证码为: ", smsCode)

	//将发送的验证码保存起来
}
