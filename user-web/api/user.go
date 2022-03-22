package api

import (
	"context"
	"fmt"
	"github.com/Dlimingliang/shop-api/user-web/global"
	"github.com/Dlimingliang/shop-api/user-web/global/response"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Dlimingliang/shop-api/user-web/proto"
)

func HandlerGrpcErrToHttpErr(err error, ctx *gin.Context) {
	if err != nil {
		if s, ok := status.FromError(err); ok {
			switch s.Code() {
			case codes.NotFound:
				ctx.JSON(http.StatusNotFound, gin.H{
					"msg": s.Message(),
				})
			case codes.Internal:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			case codes.InvalidArgument:
				ctx.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"msg": "系统内部错误",
				})
			}
			return
		}
	}
}

func GetUserList(ctx *gin.Context) {

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserServiceConfig.Host,
		global.ServerConfig.UserServiceConfig.Port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 [用户服务失败]", "msg", err.Error())
	}

	userClient := proto.NewUserClient(conn)
	rsp, err := userClient.GetUserPage(context.Background(), &proto.UserPageRequest{
		Page:     0,
		PageSize: 0,
	})
	if err != nil {
		zap.S().Errorw("[GetUserList] 查询 [用户列表]失败", "msg", err.Error())
		HandlerGrpcErrToHttpErr(err, ctx)
		return
	}

	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		userResp := response.UserResponse{
			Id:       value.Id,
			UserName: value.UserName,
			Birthday: response.JsonTime(time.Unix(int64(value.Birthday), 0)),
			Gender:   strconv.Itoa(int(value.Gender)),
			Mobile:   value.Mobile,
		}
		result = append(result, userResp)
	}
	ctx.JSON(http.StatusOK, result)
}
