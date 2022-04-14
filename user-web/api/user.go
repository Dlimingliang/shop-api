package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Dlimingliang/shop-api/user-web/forms"
	"github.com/Dlimingliang/shop-api/user-web/global"
	"github.com/Dlimingliang/shop-api/user-web/middlewares"
	"github.com/Dlimingliang/shop-api/user-web/models"
	"github.com/Dlimingliang/shop-api/user-web/models/response"
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

func GetUserList(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "0")
	pageInt, _ := strconv.Atoi(page)
	pageSize := ctx.DefaultQuery("pageSize", "10")
	pageSizeInt, _ := strconv.Atoi(pageSize)
	rsp, err := global.UserSrvClient.GetUserPage(context.Background(), &proto.UserPageRequest{
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

	//验证验证码是否正确
	if !store.Verify(passwordLoginForm.CaptchaID, passwordLoginForm.Captcha, true) {
		ctx.JSON(http.StatusOK, gin.H{"msg": "验证码错误"})
		return
	}

	//查询用户是否存在
	userRsp, err := global.UserSrvClient.GetUserByMobile(context.Background(), &proto.MobileRequest{Mobile: passwordLoginForm.Mobile})
	if err != nil {
		HandlerGrpcErrToHttpErr(err, ctx)
		return
	}

	//验证密码正确性
	ok, err := global.UserSrvClient.CheckPassword(context.Background(), &proto.PasswordCheckRequest{
		Password:          passwordLoginForm.Password,
		EncryptedPassword: userRsp.Password,
	})
	if err != nil {
		HandlerGrpcErrToHttpErr(err, ctx)
		return
	}
	if ok.Success {
		j := middlewares.NewJWT()
		claims := models.CustomClaims{
			UserId:   userRsp.Id,
			UserName: userRsp.UserName,
			RoleId:   int(userRsp.Role),
			StandardClaims: jwt.StandardClaims{
				NotBefore: time.Now().Unix() - 1000,     // 签名生效时间
				ExpiresAt: time.Now().Unix() + 60*60*24, // 签名过期时间
				Issuer:    "lml",                        // 签名颁发者
			},
		}
		token, err := j.CreateToken(claims)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg": err.Error(),
			})
		}
		ctx.JSON(http.StatusOK, gin.H{"msg": "登录成功", "token": token})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"msg": "用户名或密码不正确"})
	}
}

func Register(ctx *gin.Context) {
	//表单验证
	registerForm := forms.RegisterForm{}
	if err := ctx.ShouldBind(&registerForm); err != nil {
		HandlerValidatorErr(err, ctx)
		return
	}

	//验证验证码
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.ServerConfig.RedisConfig.Host, global.ServerConfig.RedisConfig.Port),
		Password: global.ServerConfig.RedisConfig.Password,
	})
	smsCode, err := rdb.Get(context.Background(), registerForm.Mobile).Result()
	if err == redis.Nil || registerForm.SmsCode != smsCode {
		//验证码不存在
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "验证码错误",
		})
		return
	}

	//注册用户,并且直接登录
	user, err := global.UserSrvClient.CreateUser(context.Background(), &proto.CreateUserRequest{
		UserName: registerForm.Mobile,
		Mobile:   registerForm.Mobile,
		Password: registerForm.Password,
	})
	if err != nil {
		HandlerGrpcErrToHttpErr(err, ctx)
		return
	}
	j := middlewares.NewJWT()
	claims := models.CustomClaims{
		UserId:   user.Id,
		UserName: user.UserName,
		RoleId:   int(user.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 1000,     // 签名生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24, // 签名过期时间
			Issuer:    "lml",                        // 签名颁发者
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "登录成功", "token": token})
}
