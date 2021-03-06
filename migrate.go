package pgmigrate

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	// pq is used for db driver
	_ "github.com/lib/pq"
)

// Init creates the db/migrations folder
// and creates the initial migration to create meta table
func Init() {
	migrationPath := filepath.Join(".", "db/migrations")

	if _, err := os.Stat(migrationPath); os.IsNotExist(err) {
		os.MkdirAll(migrationPath, os.ModePerm)
		ioutil.WriteFile(migrationPath+"/00000_init.sql", []byte(InitSQL), 0644)
		configPath := filepath.Join(".", "db")
		ioutil.WriteFile(configPath+"/config.yaml", []byte(ConfigYAML), 0644)
	}
}

// Migrate reads all migrations and migrates them
func Migrate() {

	loadConfig()

	pwd, _ := os.Getwd()
	migrationsPath := pwd + "/db/migrations/"

	var allMigrations []string

	filepath.Walk(migrationsPath, func(path string, file os.FileInfo, err error) error {
		migrationName := file.Name()
		if !file.IsDir() && migrationName != ".gitkeep" {
			allMigrations = append(allMigrations, migrationName)
		}
		return nil
	})

	var migrated []string
	migratedRows, err := DB.Query("SELECT name FROM _migrations")
	if err != nil {
		log.Fatal(err)
	}
	defer migratedRows.Close()
	for migratedRows.Next() {
		var name string
		if err := migratedRows.Scan(&name); err != nil {
			log.Fatal(err)
		}
		migrated = append(migrated, name)
	}

	for _, migrationName := range allMigrations {
		if !contains(migrated, migrationName) {
			migration, _ := ioutil.ReadFile(migrationsPath + migrationName)
			if _, err = DB.Exec(string(migration)); err != nil {
				log.Fatal(err)
			}
			if _, err = DB.Exec(`INSERT INTO _migrations (name) VALUES ($1)`, migrationName); err != nil {
				log.Fatal(err)
			}
			fmt.Println("Migrated: ", migrationName)
		}
	}
}
