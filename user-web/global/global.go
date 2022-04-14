package global

import (
	"github.com/Dlimingliang/shop-api/user-web/proto"
	ut "github.com/go-playground/universal-translator"

	"github.com/Dlimingliang/shop-api/user-web/config"
)

const (
	JWTGinContextKey string = "claims"
)

var (
	ServerConfig   = &config.ServerConfig{}
	ValidatorTrans ut.Translator
	UserSrvClient  proto.UserClient
)
