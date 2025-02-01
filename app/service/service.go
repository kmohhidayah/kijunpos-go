package service

import (
	"github/kijunpos/app/service/user"
	"github/kijunpos/config"
)

type Services struct {
	User user.Service
}

func New(cfg *config.Config) Services {
	return Services{
		User: user.New(cfg),
	}
}
