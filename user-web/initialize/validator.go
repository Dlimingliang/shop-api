package initialize

import (
	"fmt"
	"github.com/Dlimingliang/shop-api/user-web/global"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
	zhtranslations "github.com/go-playground/validator/v10/translations/zh"
)

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
