package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
)

func ApplyMigrations(db *sql.DB, migrationsPath string) error {
	files, err := ioutil.ReadDir(migrationsPath)
	if err != nil {
		return fmt.Errorf("error %v", err)
	}

	for _, file := range files {
		if strings.Contains(file.Name(),"up.sql"){
			filePath := filepath.Join(migrationsPath, file.Name())

			err := executeSQLFile(db, filePath)
			if err != nil {
				fmt.Println("faild at",file.Name(),err)
				continue
			}

			fmt.Println(file.Name())
		}
	}

	return nil
}

func executeSQLFile(db *sql.DB, filePath string) error {
	query, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	_, err = db.Exec(string(query))
	return err
}
