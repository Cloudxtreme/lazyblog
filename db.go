package main

import (
	rdb "github.com/dancannon/gorethink"
)

// DB is the struct the holds our session.
type DB struct {
	Session *rdb.Session
}
