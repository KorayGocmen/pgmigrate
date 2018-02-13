package pgmigrate

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// GetEnvVar gets the env var or fallbacks to default
func GetEnvVar(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// InitFiles creates the db/migrations folder
// and creates the initial migration to create meta table
func InitFiles() {
	migrationPath := filepath.Join(".", "db/migrations")
	if _, err := os.Stat(migrationPath); os.IsNotExist(err) {
		os.MkdirAll(migrationPath, os.ModePerm)
		initSQL, _ := ioutil.ReadFile("init.sql")
		ioutil.WriteFile(migrationPath+"/00000_init.sql", initSQL, 0644)
	}

	configPath := filepath.Join(".", "db")
	config, _ := ioutil.ReadFile("config.yaml")
	ioutil.WriteFile(configPath+"/config.yaml", config, 0644)
}

// Contains checks if a given string exists in an array of strings
func Contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
