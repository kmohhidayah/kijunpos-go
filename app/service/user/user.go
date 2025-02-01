package user

import "github/kijunpos/config"

type Service struct{}

func New(cfg *config.Config) Service {
	return Service{}
}
