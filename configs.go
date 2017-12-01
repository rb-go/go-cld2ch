package main

import (
	"database/sql"
)

var (
	version, build, buildDate string
)

type conf struct {
	System struct {
		ListenOn   string `yaml:"listenOn"`
		MaxThreads int    `yaml:"maxThreads"`
	}
	Clickhouse struct {
		ConnString string `yaml:"connString"`
		DBName     string `yaml:"dbname"`
	}
}

//Configs ...
var config conf

var configPath string

var clickHouseDB *sql.DB
