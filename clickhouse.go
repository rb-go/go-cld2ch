package main

import (
	"database/sql"
	"log"

	"time"

	"github.com/kshvakov/clickhouse"
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
	createCollectdTable()
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

func createCollectdTable() {
	_, err := clickHouseDB.Exec(`
CREATE TABLE IF NOT EXISTS CollectD (
	EventDate DEFAULT toDate(now()),
	EventTime DateTime DEFAULT now(),
	Hostname String,
	Plugin String,
	ParamName String,
	ParamValue Float64
) ENGINE = MergeTree(EventDate, (Hostname, Plugin, ParamName, EventTime, EventDate), 8192)
`)
	if err != nil {
		log.Fatalln("Creating table error:", err)
	}
}

type cdElementData struct {
	EventDateTime time.Time
	Hostname      string
	Plugin        string
	ParamName     string
	ParamValue    float64
}

func insertCollectDToCH(insData []cdElementData) error {
	if len(insData) < 1 {
		return nil
	}
	var (
		tx, _   = clickHouseDB.Begin()
		stmt, _ = tx.Prepare("INSERT INTO CollectD VALUES (?, ?, ?, ?, ?, ?")
	)
	for _, element := range insData {
		if _, err := stmt.Exec(
			element.EventDateTime,
			element.EventDateTime,
			element.Hostname,
			element.Plugin,
			element.ParamName,
			element.ParamValue,
		); err != nil {
			return err
		}
	}
	err := tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
