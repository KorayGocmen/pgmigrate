package pgmigrate

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	// pq is used for db driver
	_ "github.com/lib/pq"
	yaml "gopkg.in/yaml.v2"
)

// DatabaseConfig struct for the DBConfig
type DatabaseConfig struct {
	Development struct {
		Username string
		Password string
		Database string
		Host     string
		SslMode  string
		Driver   string
	}
	Test struct {
		Username string
		Password string
		Database string
		Host     string
		SslMode  string
		Driver   string
	}
	Production struct {
		Username string
		Password string
		Database string
		Host     string
		SslMode  string
		Driver   string
	}
}

// DBConfig the config file values
var (
	DBConfig DatabaseConfig
	DBConn   string
	DB       *sql.DB
)

func (dc *DatabaseConfig) readConfig() {
	pwd, _ := os.Getwd()
	yamlFile, err := ioutil.ReadFile(pwd + "/db/config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, dc)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
}

// Init creates the db connection
func Init() {

	InitFiles()

	DBConfig.readConfig()
	env := GetEnvVar("ENV", "development")

	if env == "production" {
		DBConn = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s",
			DBConfig.Production.Username,
			DBConfig.Production.Password,
			DBConfig.Production.Database,
			DBConfig.Production.SslMode,
		)
	} else if env == "test" {
		DBConn = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s",
			DBConfig.Test.Username,
			DBConfig.Test.Password,
			DBConfig.Test.Database,
			DBConfig.Test.SslMode,
		)
	} else {
		DBConn = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s",
			DBConfig.Development.Username,
			DBConfig.Development.Password,
			DBConfig.Development.Database,
			DBConfig.Development.SslMode,
		)
	}

	var err error
	DB, err = sql.Open("postgres", DBConn)
	if err != nil {
		log.Fatal(err)
	}

	pwd, _ := os.Getwd()
	migration, err := ioutil.ReadFile(pwd + "/db/migrations/00000_init.sql")
	if err != nil {
		log.Fatal(err)
	}
	if _, err = DB.Exec(string(migration)); err != nil {
		log.Fatal(err)
	}
}

// Migrate reads all migrations and migrates them
func Migrate() {
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
		if !Contains(migrated, migrationName) {
			migration, _ := ioutil.ReadFile(migrationsPath + migrationName)
			if _, err = DB.Exec(string(migration)); err != nil {
				fmt.Println("Error migrating", migrationName, err)
			}
			if _, err = DB.Exec(`INSERT INTO _migrations (name) VALUES ($1)`, migrationName); err != nil {
				fmt.Println("Error inserting migration name", migrationName, err)
			}
			fmt.Println("Migrated: ", migrationName)
		}
	}
}
