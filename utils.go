package pgmigrate

import (
	"io/ioutil"
	"os"
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
	if _, err := os.Stat("/db/migrations/00000_init.sql"); os.IsNotExist(err) {
		os.MkdirAll("/db/migrations/00000_init.sql", os.ModePerm)
		initSQL, _ := ioutil.ReadFile("init.sql")
		ioutil.WriteFile("/db/migrations/00000_init.sql", initSQL, 0644)
	}

	if _, err := os.Stat("/db/config.yaml"); os.IsNotExist(err) {
		os.MkdirAll("/db/config.yaml", os.ModePerm)
		config, _ := ioutil.ReadFile("config.yaml")
		ioutil.WriteFile("/db/config.yaml", config, 0644)
	}
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
