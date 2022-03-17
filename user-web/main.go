package main

import (
	"flag"
	"fmt"
	"github.com/Dlimingliang/shop-api/user-web/initialize"
)

var (
	ip   = flag.String("ip", "0.0.0.0", "IP")
	port = flag.Int("port", 8090, "端口号")
)

func main() {
	ginRouter := initialize.Routers()

	flag.Parse()
	if err := ginRouter.Run(fmt.Sprintf("%s:%d", ip, port)); err != nil {
		panic(err)
	}
}
