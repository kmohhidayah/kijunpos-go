package repository

import "github/kijunpos/config/db"

type Repository struct{}

func New(kijundbConn *db.Connection) Repository {
	return Repository{}
}
