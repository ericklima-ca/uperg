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
	driverDNS := strings.SplitN(opt.DNS, ":", 2)
	driver, dns := driverDNS[0], driverDNS[1]
	db, err := sql.Open(driver, dns)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Panicln("db not connected!")
	}
	return db
}
