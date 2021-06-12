package wbDB

import (
	"database/sql"
	"fmt"

	// sql driver
	_ "github.com/go-sql-driver/mysql"
)

var (
	// Core database
	coreDB *sql.DB
)

// Database configuration
type dbConfig struct {
	// Host to connect to (ip or url)
	host string
	// port to connect to
	port int
	// user username
	user string
	// user password
	passwd string
}

func (c *dbConfig) ConnecterString(database string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", c.user, c.passwd, c.host, c.port, database)
}

var configDB = dbConfig{
	host:   "127.0.0.1",
	port:   3307,
	user:   "test_user",
	passwd: "test_passwd",
}

func Init() error {
	// Create the database handle
	var err error
	coreDB, err = sql.Open("mysql", configDB.ConnecterString("wb_core"))
	if err != nil {
		return err
	}
	// Some drivers may not create a connection initally, this verifies that the connection is alive
	err = coreDB.Ping()
	return err
}

func Get() (exists bool, coreDB *sql.DB) {
	return coreDB != nil, coreDB
}
