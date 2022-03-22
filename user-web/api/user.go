package api

import (
	"context"
	"fmt"
	"github.com/Dlimingliang/shop-api/user-web/forms"
	"github.com/Dlimingliang/shop-api/user-web/global"
	"github.com/Dlimingliang/shop-api/user-web/global/response"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Dlimingliang/shop-api/user-web/proto"
)

func HandlerValidatorErr(err error, ctx *gin.Context) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	ctx.JSON(http.StatusBadRequest, gin.H{
		"error": removeTopSruct(errs.Translate(global.ValidatorTrans)),
	})
}

func removeTopSruct(fileds map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fileds {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

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
	page := ctx.DefaultQuery("page", "0")
	pageInt, _ := strconv.Atoi(page)
	pageSize := ctx.DefaultQuery("pageSize", "10")
	pageSizeInt, _ := strconv.Atoi(pageSize)
	rsp, err := userClient.GetUserPage(context.Background(), &proto.UserPageRequest{
		Page:     uint32(pageInt),
		PageSize: uint32(pageSizeInt),
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

func PasswordLogin(ctx *gin.Context) {
	//表单验证
	passwordLoginForm := forms.PasswordLoginForm{}
	if err := ctx.ShouldBind(&passwordLoginForm); err != nil {
		HandlerValidatorErr(err, ctx)
		return
	}
}
