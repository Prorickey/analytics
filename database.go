package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var glob_db *sql.DB

func MustConnect() {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
	os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"),
	os.Getenv("DB_PASSWORD"), os.Getenv("DB_DATABASE"))

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Error opening connection: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	var dat []byte
	if os.Getenv("GIN_MODE") == "release" {
		dat, err = os.ReadFile("/root/schema.sql")
		if err != nil {
			log.Fatalf("Error reading schema file: %v", err)
		}
	} else {
		dat, err = os.ReadFile("./schema.sql")
		if err != nil {
			log.Fatalf("Error reading schema file: %v", err)
		}
	}

	schema := string(dat)
	_, err = db.Exec(schema)
	if err != nil {
		log.Fatalf("Error executing schema: %v", err)
	}

	glob_db = db

	fmt.Println("Successfully connected to database!")
}

func CloseDatabase() {
	glob_db.Close()
}

func CreateEvent(event string, metadata string) error {
	_, err := glob_db.Exec("INSERT INTO analytics(event, metadata) VALUES($1, $2)", event, metadata)
	if err != nil {
		log.Printf("Error creating event: %v", err)
		return err
	}

	return nil
}

func ValidateAuthorization(id string, token string) bool {
	var temp string
	err := glob_db.QueryRow("SELECT service FROM service_auth WHERE id=$1 AND token=$2", id, token).Scan(&temp)
	return err == nil
}

type EventRecord struct {
	Timestamp time.Time `json:"timestamp"`
}

func GetGroupedRecords(event string, startTime time.Time, endTime time.Time, table string, window int) ([]GroupedRecord, error) {
	stmt := fmt.Sprintf("SELECT timestamp, count FROM %s WHERE event=$1 AND timestamp>=$2 AND timestamp<=$3 ORDER BY timestamp", table)
	rows, err := glob_db.Query(stmt, event, startTime, endTime)
	if err != nil {
		log.Printf("Error selecting from analytics: %v", err)
		return []GroupedRecord{}, err 
	}
	defer rows.Close()

	records := make([]GroupedRecord, 0)

	for rows.Next() {
		var record GroupedRecord
		if err := rows.Scan(&record.Timestamp, &record.Count); err != nil {
			log.Printf("Error scanning row: %v", err)
			return records, err
		}

		record.Event = event 
		record.Window = window

		records = append(records, record)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error concerning rows: %v", err)
		return records, err 
	}

	return records, nil
}

func GetEventRecords(event string, startTime time.Time, endTime time.Time) ([]EventRecord, error) {
	rows, err := glob_db.Query("SELECT timestamp FROM analytics WHERE event=$1 AND timestamp>=$2 AND timestamp<=$3 ORDER BY timestamp", 
	event, startTime, endTime)
	if err != nil {
		log.Printf("Error selecting from analytics: %v", err)
		return []EventRecord{}, err 
	}
	defer rows.Close()

	records := make([]EventRecord, 0)

	for rows.Next() {
		var record EventRecord
		if err := rows.Scan(&record.Timestamp); err != nil {
			log.Printf("Error scanning row: %v", err)
			return records, err
		}

		records = append(records, record)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error concerning rows: %v", err)
		return records, err 
	}

	return records, nil
}