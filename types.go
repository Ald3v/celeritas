package celeritas

import "database/sql"

type initPaths struct {
	rootPath    string
	folderNames []string
}

type cookieConfig struct {
	name     string
	lifetime string
	persist  string
	secure   string
	domain   string
}

type databaseConfig struct {
	dbType string 
	dsn      string
	database string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime string
}

type limiterConfig struct {
	Rps float64
	Burst int
	Enabled bool
}

type Database struct {
	DataType string
	Pool     *sql.DB
}