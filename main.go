package main

import (
	"github/kijunpos/app"
	"github/kijunpos/cmd/grpc"
)

func main() {
	app.Init()
	grpc.StartServer()
}
