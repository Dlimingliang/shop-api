package global

import (
	"github.com/Dlimingliang/shop-api/user-web/proto"
	ut "github.com/go-playground/universal-translator"
	"google.golang.org/grpc"

	"github.com/Dlimingliang/shop-api/user-web/config"
)

const (
	JWTGinContextKey string = "claims"
)

var (
	ServerConfig   = &config.ServerConfig{}
	ValidatorTrans ut.Translator
	UserSrvConn    = &grpc.ClientConn{}
	UserSrvClient  proto.UserClient
)
