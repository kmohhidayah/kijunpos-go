package user

import "github/kijunpos/config/db"

type User struct {
	Query   Query
	Command Command
}

type Query struct {
	dbConn *db.Connection
}

type Command struct {
	dbConn *db.Connection
}

func newQuery(dbConn *db.Connection) Query {
	return Query{dbConn: dbConn}
}
func newCommand(dbConn *db.Connection) Command {
	return Command{dbConn: dbConn}
}

func New(dbConn *db.Connection) User {
	return User{
		Query:   newQuery(dbConn),
		Command: newCommand(dbConn),
	}
}
