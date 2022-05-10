package api

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Dlimingliang/shop-api/goods-web/global"
	"github.com/Dlimingliang/shop-api/goods-web/models/response"
	"github.com/Dlimingliang/shop-api/goods-web/proto"
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
			case codes.InvalidArgument:
				ctx.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.NotFound:
				ctx.JSON(http.StatusNotFound, gin.H{
					"msg": s.Message(),
				})
			case codes.Internal:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"msg": "系统内部错误",
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

func GetGoodsPage(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "0")
	pageInt, _ := strconv.Atoi(page)
	pageSize := ctx.DefaultQuery("pageSize", "10")
	pageSizeInt, _ := strconv.Atoi(pageSize)
	rsp, err := global.GoodsSrvClient.GetGoodsPage(context.Background(), &proto.GoodsPageReq{
		Pages:    int32(pageInt),
		PageSize: int32(pageSizeInt),
	})
	if err != nil {
		zap.S().Errorw("[GetUserPage] 查询 [商品列表]失败", "msg", err.Error())
		HandlerGrpcErrToHttpErr(err, ctx)
		return
	}

	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		userResp := response.GoodsResponse{
			Id:   value.Id,
			Name: value.Name,
		}
		result = append(result, userResp)
	}
	ctx.JSON(http.StatusOK, result)
}
