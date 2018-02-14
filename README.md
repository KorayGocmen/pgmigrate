## `pgmigrate`

pgmigrate is a simple package for postgresql migrations

### Usage

You can create the connection configuration file and initialization migration along with the necessary folders by using:

```go
pgmigrate.Init()
```

To run migrations use:

```go
pgmigrate.Migrate()
```

This package does not utilize up/down migration system. A down migration is just another migration.