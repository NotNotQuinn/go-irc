package data

import (
	"database/sql"

	"github.com/NotNotQuinn/go-irc/config"
	_ "github.com/go-sql-driver/mysql"
)

var (
	// Core database - should NEVER be nil.
	CoreDB *sql.DB
)

func init() {
	// Create the database handle
	var err error
	conf := config.Private.Database
	CoreDB, err = sql.Open("mysql", conf.ConnecterString("wb_core"))
	if err != nil {
		panic(err)
	}
	// Some drivers may not create a connection initally, this verifies that the connection is alive
	err = CoreDB.Ping()
	if err != nil {
		panic(err)
	}
}
