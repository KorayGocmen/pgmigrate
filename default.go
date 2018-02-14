package pgmigrate

// InitSQL is the initialization for the meta table
var InitSQL = `CREATE TABLE IF NOT EXISTS _migrations (
  name varchar(255) NOT NULL
);`

// ConfigYAML is the default db config file
var ConfigYAML = `development:
  username: ""
  password: ""
  database: ""
  host: ""
  sslmode: "disable"

test:
  username: ""
  password: ""
  database: ""
  host: ""
  sslmode: "disable"

production:
  username: ""
  password: ""
  database: ""
  host: ""
  sslmode: "enable"`
