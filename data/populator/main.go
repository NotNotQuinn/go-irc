package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/NotNotQuinn/go-irc/config"
	"github.com/NotNotQuinn/go-irc/data"
)

func main() {
	err := config.InitForTests("./config")
	if err != nil {
		panic(err)
	}
	err = data.Init()
	if err != nil {
		panic(err)
	}
	count, err := runSqlFiles("./data/sql")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Ran %d files.", count)
}

// Runs SQL files in a directory (or just one file) in alphabetical order.
// Works best when files are labeled on their order, e.g. `00-database1.sql`, `01-table1.sql`
//
// Will only run files ending with `.sql`
//
// `path`: path to file to read and execute as SQL. Will run recursively on directories.
func runSqlFiles(dir string) (fileCount int, err error) {
	stat, err := os.Stat(dir)
	if err != nil {
		return fileCount, err
	}
	if !stat.IsDir() {
		return fileCount, errors.New("not a directory")
	}

	// List of SQL queries to execute.
	// var ExecQueue = []string{}
	if items, err := os.ReadDir(dir); err == nil {
		for _, item := range items {
			if item.IsDir() {
				continue
			}
			bytes, err := os.ReadFile(filepath.Join(dir, item.Name()))
			if err != nil {
				return fileCount, err
			}
			query := string(bytes)
			fmt.Printf("Exec: \n%s\n", query)
			res, err := data.CoreDB.Exec(query)
			if err != nil {
				return fileCount, err
			}
			rows, err := res.RowsAffected()
			if err != nil {
				return fileCount, err
			}
			fmt.Printf("%s: %d rows affected\n", item.Name(), rows)
		}
	} else {
		return fileCount, err
	}
	return fileCount, nil
}
