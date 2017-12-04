package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/kshvakov/clickhouse"
	"github.com/patrickmn/go-cache"
)

func connectClickDB() {
	var err error
	log.Println("Connecting to Clickhouse by connString", config.Clickhouse.ConnString)
	clickHouseDB, err = sql.Open("clickhouse", config.Clickhouse.ConnString)
	if err != nil {
		log.Fatalln("Clickhouse connection error: ", err, config.Clickhouse.ConnString)
	}

	log.Println("Connected to Clickhouse, pinging")
	if err := clickHouseDB.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			log.Fatalf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			log.Fatalln("Clickhouse ping error: ", err, config.Clickhouse.ConnString)
		}
	}
	log.Println("Ping to Clickhouse - OK")

	createClickDB()
	tables := getTablesList()
	for _, tbl := range tables {
		cached.Set(tbl, 1, cache.NoExpiration)
	}
}

func createClickDB() {
	_, err := clickHouseDB.Exec(`CREATE DATABASE IF NOT EXISTS ` + config.Clickhouse.DBName)
	if err != nil {
		log.Fatalln("Creating database error:", err)
	}
	_, err = clickHouseDB.Exec("USE " + config.Clickhouse.DBName)
	if err != nil {
		log.Fatalln("Selecting db for session error:", err)
	}
}

func getTablesList() []string {
	rows, err := clickHouseDB.Query("SHOW TABLES")
	if err != nil {
		log.Fatalln("Error on SHOW TABLES:", err)
	}
	var tables []string

	for rows.Next() {
		var tableRow string
		if err := rows.Scan(&tableRow); err != nil {
			log.Fatalln("Error on scanning SHOW TABLES:", err)
		}
		tables = append(tables, tableRow)
	}
	return tables
}

func createCollectdTable(tableName string) {
	_, err := clickHouseDB.Exec(`
CREATE TABLE IF NOT EXISTS ` + tableName + ` (
	EventDate DEFAULT toDate(EventTime),
	EventTime DateTime DEFAULT now(),
	Hostname String,
	ParamName String,
	ParamValue Float64
) ENGINE = MergeTree(EventDate, (Hostname, EventTime), 8192)
`)
	if err != nil {
		log.Fatalln("Creating table error:", err)
	}
	cached.Set(tableName, 1, cache.NoExpiration)
}

func checkAndCreateTableIfNeed(table string) {
	_, foundTable := cached.Get(table)
	if foundTable != true {
		createCollectdTable(table)
	}
}

type cdElementData struct {
	EventDateTime time.Time
	Hostname      string
	Plugin        string
	ParamName     string
	ParamValue    float64
}

func insertCollectDToCH(insData []cdElementData) {
	if len(insData) < 1 {
		return
	}
	for _, element := range insData {
		checkAndCreateTableIfNeed(element.Plugin)
		var tx, _ = clickHouseDB.Begin()
		var stmt, _ = tx.Prepare(`INSERT INTO ` + element.Plugin + ` (EventTime, Hostname, ParamName, ParamValue) VALUES (?, ?, ?, ?)`)
		_, err := stmt.Exec(
			element.EventDateTime,
			element.Hostname,
			element.ParamName,
			element.ParamValue,
		)
		if err != nil {
			log.Fatalln("Error on EXEC Statement:", err)
		}
		if err := tx.Commit(); err != nil {
			log.Fatalln("Error on Commit Statement:", err)
		}
	}
}
