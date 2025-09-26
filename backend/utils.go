package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/argon2"
)

const (
	h_time = 3
	h_memory = 32*1024
	h_threads = 4
	h_keyLen = 32
)

func HashPassword(password string) string {
	salt := make([]byte, 32)
	rand.Read(salt)
	key := argon2.Key([]byte(password), salt, h_time, h_memory, h_threads, h_keyLen)

	hash := base64.StdEncoding.EncodeToString(key)
	salt_b64 := base64.StdEncoding.EncodeToString(salt)

	return fmt.Sprintf("$argon2id$v=1$m=%d,t=%d,p=%d$%s$%s", 
	h_memory, h_time, h_threads, salt_b64, hash)
}

func ValidatePassword(password string, digest string) bool {
	parts := strings.Split(digest, "$")
	
	params := strings.Split(parts[3], ",")
	if len(params) < 3 {
		log.Printf("Error validating password: %s", digest)
		return false
	}

	memoryStr := strings.TrimPrefix(params[0], "m=")
	timeStr := strings.TrimPrefix(params[1], "t=")
	threadsStr := strings.TrimPrefix(params[2], "p=")

	memoryDigest, err := strconv.Atoi(memoryStr)
	if err != nil {
		log.Printf("Error converting memory string to integer: %v", err)
		return false
	}

	timeDigest, err := strconv.Atoi(timeStr)
	if err != nil {
		log.Printf("Error converting time string to integer: %v", err)
		return false
	}

	threadsDigest, err := strconv.Atoi(threadsStr)
	if err != nil {
		log.Printf("Error converting threads string to integer: %v", err)
		return false
	}

	salt, err := base64.StdEncoding.DecodeString(parts[4])
	if err != nil {
		log.Printf("Error decoding salt from base64: %v", err)
		return false
	}

	hash, err := base64.StdEncoding.DecodeString(parts[5])
	if err != nil {
		log.Printf("Error decoding hash from base64: %v", err)
		return false 
	}

	checkHash := argon2.Key([]byte(password), salt, uint32(timeDigest), 
	uint32(memoryDigest), uint8(threadsDigest), uint32(len(hash)))
	return bytes.Equal(checkHash, hash)
}

type GroupedRecord struct {
	Timestamp time.Time `json:"timestamp"`
	Window int `json:"window"`
	Event string `json:"event"`
	Count int64 `json:"count"`
}

func GroupRecords(event string, records []EventRecord, window int) []GroupedRecord {
	groupedRecords := make([]GroupedRecord, 0)
	currentTime := records[0].Timestamp
	count := 0
	for _, rec := range records {
		if rec.Timestamp.Equal(currentTime) || (rec.Timestamp.After(currentTime) && rec.Timestamp.Before(currentTime.Add(time.Duration(window) * time.Second))) {
			count++
		} else {
			groupedRecords = append(groupedRecords, GroupedRecord{
				Timestamp: currentTime,
				Window: window,
				Event: event,
				Count: int64(count),
			})

			currentTime = rec.Timestamp
			count = 1
		}
	}

	groupedRecords = append(groupedRecords, GroupedRecord{
		Timestamp: currentTime,
		Window: window,
		Event: event,
		Count: int64(count),
	})

	return groupedRecords
}