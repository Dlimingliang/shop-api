package global

import (
	ut "github.com/go-playground/universal-translator"

	"github.com/Dlimingliang/shop-api/user-web/config"
)

const (
	JWTGinContextKey string = "claims"
)

var (
	ServerConfig   = &config.ServerConfig{}
	ValidatorTrans ut.Translator
)
