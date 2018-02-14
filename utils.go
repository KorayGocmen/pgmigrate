package pgmigrate

import (
	"database/sql"
	"io/ioutil"
	"log"
	"os"

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
	}
	Test struct {
		Username string
		Password string
		Database string
		Host     string
		SslMode  string
	}
	Production struct {
		Username string
		Password string
		Database string
		Host     string
		SslMode  string
	}
}

// DBConfig the config file values
var (
	DBConfig DatabaseConfig
	DBConn   string
	DB       *sql.DB
)

// getEnvVar gets the env var or fallbacks to default
func getEnvVar(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// contains checks if a given string exists in an array of strings
func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

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
func loadConfig() {

	DBConfig.readConfig()
	env := getEnvVar("ENV", "development")

	if env == "production" {
		if DBConfig.Production.Username != "" {
			DBConn += "user=" + DBConfig.Production.Username
		}
		if DBConfig.Production.Password != "" {
			DBConn += " password=" + DBConfig.Production.Password
		}
		if DBConfig.Production.Database != "" {
			DBConn += " dbname=" + DBConfig.Production.Database
		}
		if DBConfig.Production.SslMode != "" {
			DBConn += " sslmode=" + DBConfig.Production.SslMode
		}
	} else if env == "test" {
		if DBConfig.Test.Username != "" {
			DBConn += "user=" + DBConfig.Test.Username
		}
		if DBConfig.Test.Password != "" {
			DBConn += " password=" + DBConfig.Test.Password
		}
		if DBConfig.Test.Database != "" {
			DBConn += " dbname=" + DBConfig.Test.Database
		}
		if DBConfig.Test.SslMode != "" {
			DBConn += " sslmode=" + DBConfig.Test.SslMode
		}
	} else {
		if DBConfig.Development.Username != "" {
			DBConn += "user=" + DBConfig.Development.Username
		}
		if DBConfig.Development.Password != "" {
			DBConn += " password=" + DBConfig.Development.Password
		}
		if DBConfig.Development.Database != "" {
			DBConn += " dbname=" + DBConfig.Development.Database
		}
		if DBConfig.Development.SslMode != "" {
			DBConn += " sslmode=" + DBConfig.Development.SslMode
		}
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
