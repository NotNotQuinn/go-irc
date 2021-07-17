package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
			rawQueries := strings.Split(string(bytes), ";")
			var queries []string
			for i, query := range rawQueries {
				// Append missing semicolon, Trim query and ignore if its empty
				if i+1 != len(query) {
					query += ";"
				}
				trimmed := strings.Trim(query, " \r\n\t")
				if trimmed == "" || trimmed == ";" {
					continue
				}
				queries = append(queries, query)
			}

			for i, query := range queries {
				// Still use the original string
				res, err := data.CoreDB.Exec(query)
				if err != nil {
					return fileCount, err
				}
				rows, err := res.RowsAffected()
				if err != nil {
					return fileCount, err
				}
				fmt.Printf("%s (%d/%d): %d rows affected\n", item.Name(), i+1, len(queries), rows)
			}
			fileCount++
		}
	} else {
		return fileCount, err
	}
	return fileCount, nil
}
