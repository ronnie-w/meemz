package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

const (
	DRIVER_NAME = "mysql"
	DATA_SRC    = "meemz:R0nni3W3k35@@/meemz"
)

var db *sql.DB
var err error

func init() {
	db, err = sql.Open(DRIVER_NAME, DATA_SRC)
	if err != nil {
		log.Panicln(err)
	}
}

func Conn() *sql.DB {
	return db
}
