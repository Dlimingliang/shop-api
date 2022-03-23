package initialize

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
	zhtranslations "github.com/go-playground/validator/v10/translations/zh"
	"go.uber.org/zap"

	"github.com/Dlimingliang/shop-api/user-web/custom_validator"
	"github.com/Dlimingliang/shop-api/user-web/global"
)

func InitValidator(local string) {
	//初始化翻译
	InitValidatorTrans(local)
	//注册自定义验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("mobile", custom_validator.ValidateMobile)
		if err != nil {
			zap.S().Panic("注册手机号验证器失败", err.Error())
		}
		err = v.RegisterTranslation("mobile", global.ValidatorTrans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0} 手机号码不合法!", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
		if err != nil {
			zap.S().Panic("注册手机号翻译器失败", err.Error())
		}
	}
}

func InitValidatorTrans(locale string) (err error) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {

		//注册一个获取jsontag的方法
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
		zhT := zh.New()
		enT := en.New()
		uni := ut.New(enT, zhT)
		global.ValidatorTrans, ok = uni.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s)", locale)
		}

		switch locale {
		case "en":
			entranslations.RegisterDefaultTranslations(v, global.ValidatorTrans)
		case "zh":
			zhtranslations.RegisterDefaultTranslations(v, global.ValidatorTrans)
		default:
			zhtranslations.RegisterDefaultTranslations(v, global.ValidatorTrans)
		}
		return
	}
	return
}
