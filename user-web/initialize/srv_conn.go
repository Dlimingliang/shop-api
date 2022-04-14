package initialize

import (
	"fmt"

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/Dlimingliang/shop-api/user-web/global"
	"github.com/Dlimingliang/shop-api/user-web/proto"
)

func InitSrvConn() {
	cfg := api.DefaultConfig()
	consulConfig := global.ServerConfig.ConsulConfig
	cfg.Address = fmt.Sprintf("%s:%d", consulConfig.Host, consulConfig.Port)
	userSrvHost := ""
	userSrvPort := 0

	client, err := api.NewClient(cfg)
	if err != nil {
		zap.S().Panic("连接consul失败", err.Error())
	}
	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf(`service == user-srv`))
	for _, value := range data {
		userSrvHost = value.Address
		userSrvPort = value.Port
		break
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", userSrvHost, userSrvPort), grpc.WithInsecure())
	if err != nil {
		zap.S().Panic("连接[用户服务失败]", "msg", err.Error())
		return
	}
	userClient := proto.NewUserClient(conn)
	global.UserSrvClient = userClient
}
