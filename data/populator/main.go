package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/NotNotQuinn/go-irc/data"
)

func main() {
	stats, err := runSqlFiles("./data/sql")
	fmt.Printf("Ran %d files. (%d queries, %d rows affected)\n", stats.NumFiles, stats.NumQueries, stats.NumRows)
	if err != nil {
		panic(err)
	}
}

// Runs SQL files in a directory in alphabetical order.
// Works best when files are labeled on their order, e.g. `00-database1.sql`, `01-table1.sql`
//
// Warning: SQL is split crudely by semi-colon and each part run individually.
//
// dir: Path to directory containing SQL files. Ignores directories and other files.
func runSqlFiles(dir string) (stats struct{ NumFiles, NumQueries, NumRows int }, err error) {
	stat, err := os.Stat(dir)
	if err != nil {
		return stats, err
	}
	if !stat.IsDir() {
		return stats, errors.New("not a directory")
	}

	if items, err := os.ReadDir(dir); err == nil {
		for _, item := range items {
			if item.IsDir() || !strings.HasSuffix(item.Name(), ".sql") {
				// Ignore
				continue
			}
			bytes, err := os.ReadFile(filepath.Join(dir, item.Name()))
			if err != nil {
				return stats, err
			}
			rawQueries := strings.Split(string(bytes), ";")
			var queries []string
			// Must be a seperate loop to count the ammount of actual queries are in a file beforehand
			for i, query := range rawQueries {
				// Append missing semicolon, ignore if trimming it results in empty string
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
				// Execute, print stats, and update stats
				res, err := data.CoreDB.Exec(query)
				if err != nil {
					return stats, err
				}
				rows, err := res.RowsAffected()
				if err != nil {
					return stats, err
				}

				fmt.Printf("%s (%d/%d): %d rows affected\n", item.Name(), i+1, len(queries), rows)
				stats.NumRows += int(rows)
				stats.NumQueries++
			}
			stats.NumFiles++
		}
	} else {
		return stats, err
	}
	return stats, nil
}
