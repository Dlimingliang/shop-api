package initialize

import (
	"fmt"

	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/Dlimingliang/shop-api/goods-web/global"
	"github.com/Dlimingliang/shop-api/goods-web/proto"
)

func InitSrvConn() {

	conn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", global.ServerConfig.ConsulConfig.Host, global.ServerConfig.ConsulConfig.Port, "goods-srv"),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Panic("连接[商品服务失败]", "msg", err.Error())
		return
	}
	global.GoodsSrvConn = conn
	userClient := proto.NewGoodsClient(conn)
	global.GoodsSrvClient = userClient
}
