package db

import (
	"database/sql"
	"log"
	"strings"
)

type Options struct {
	DNS string
}

func NewConnection(opt Options) *sql.DB {
	driver := strings.SplitN(opt.DNS, ":", 2)[0]
	db, err := sql.Open(driver, opt.DNS)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Panicln("db not connected! ", err)
	}
	return db
}
