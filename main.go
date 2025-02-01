package main

import (
	"github/kijunpos/app"
	"github/kijunpos/cmd/grpc"
)

func main() {
	setupData := app.Init()
	grpc.StartGrpcServer(setupData)
}
