package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

const (
	host		= "localhost"
	port		= 5432
	user		= "postgres"
	password 	= "password"
	dbname	 	= "postgres"
)

var glob_db *sql.DB

func MustConnect() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
	host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Error opening connection: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	dat, err := os.ReadFile("./schema.sql")
	if err != nil {
		log.Fatalf("Error reading schema file: %v", err)
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

type User struct {
	Id string `json:"id"`
	Digest string `json:"-"`
	Token string `json:"token"`
}

func LoginAccount(username string, password string) (User, error) {
	var user User
	err := glob_db.QueryRow("SELECT id, digest, token FROM users WHERE username=$1", username).Scan(&user.Id, &user.Digest, &user.Token)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("Error selecting user: %v", err)
		}

		return User{}, err
	}

	if ValidatePassword(password, user.Digest) {
		return user, nil
	} else {
		return User{}, fmt.Errorf("invalid password")
	}
}

func ValidateAuthorization(id string, token string) bool {
	var temp string
	err := glob_db.QueryRow("SELECT username FROM users WHERE id=$1 AND token=$2", id, token).Scan(&temp)
	return err == nil
}

type EventRecord struct {
	Timestamp time.Time `json:"timestamp"`
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