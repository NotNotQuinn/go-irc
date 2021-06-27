package data

import (
	"database/sql"

	"github.com/NotNotQuinn/go-irc/config"
	// sql driver
	_ "github.com/go-sql-driver/mysql"
)

var (
	// Core database
	CoreDB *sql.DB
)

func Init() error {
	// Create the database handle
	var err error
	conf := config.Private.Database
	CoreDB, err = sql.Open("mysql", conf.ConnecterString("wb_core"))
	if err != nil {
		return err
	}
	// Some drivers may not create a connection initally, this verifies that the connection is alive
	err = CoreDB.Ping()
	return err
}
