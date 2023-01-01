package main

import (
	"MDServer/parameter"
	"MDServer/server"
)

func main() {

	res := parameter.InitPar()
	if res { //由于不能双向引用则采用提升变量来启动服务器
		server.Listen()
	}
}
