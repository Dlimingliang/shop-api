package global

import (
	ut "github.com/go-playground/universal-translator"
	"google.golang.org/grpc"

	"github.com/Dlimingliang/shop-api/goods-web/config"
	"github.com/Dlimingliang/shop-api/goods-web/proto"
)

const (
	JWTGinContextKey string = "claims"
)

var (
	ServerConfig   = &config.ServerConfig{}
	ValidatorTrans ut.Translator
	GoodsSrvConn   = &grpc.ClientConn{}
	GoodsSrvClient proto.GoodsClient
)
