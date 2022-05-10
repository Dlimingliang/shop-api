package main

import (
	"flag"
	"fmt"

	"go.uber.org/zap"

	"github.com/Dlimingliang/shop-api/goods-web/global"
	"github.com/Dlimingliang/shop-api/goods-web/initialize"
)

var (
	ip   = flag.String("ip", "0.0.0.0", "IP")
	port = flag.Int("port", 8070, "端口号")
)

func main() {
	//初始化logger
	initialize.InitLogger()
	//初始化config
	initialize.InitConfig()
	//初始化路由
	ginRouter := initialize.Routers()
	//初始换validator
	initialize.InitValidator("zh")
	//初始化grpc服务连接
	initialize.InitSrvConn()
	defer global.GoodsSrvConn.Close()

	flag.Parse()
	zap.S().Info(fmt.Sprintf("shop-api项目启动, 访问地址: http://%s:%d", *ip, *port))
	if err := ginRouter.Run(fmt.Sprintf("%s:%d", *ip, *port)); err != nil {
		zap.S().Panic(err)
	}
}
