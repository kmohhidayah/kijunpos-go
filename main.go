package main

import (
	"github/kijunpos/app"
	"github/kijunpos/cmd/grpc"
)

func main() {
	// init app
	app.Init()
	grpc.StartServer()
}
