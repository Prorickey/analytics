package main

import (
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	MustConnect()
	defer CloseDatabase()

	log.Printf("testtest: %s", HashPassword("testtest"))

	// This will generate the database in the background
	// It also sort of simulates normal conditions of an
	// always growing database
	doomAndDespair()

	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	router.POST("/event/:event", PostEvent)
	router.POST("/login", PostLogin)
	router.GET("/analytics", AuthMiddlware, GetAnalytics)

	router.Run("0.0.0.0:8080")
}

func doomAndDespair() {
	const numThreads = 2

	for range numThreads {
		go func() {
			for {
				_, err := glob_db.Exec("INSERT INTO analytics(event, timestamp) VALUES('testEvent', NOW() - (random() * interval '7 days'));")
				if err != nil {
					log.Fatalf("Error inserting into analytics: %v", err)
				}
			}
		}()
	}
}

const (
	defaultLength = 7 * 24 * time.Hour
)

func GetAnalytics(ctx *gin.Context) {
	event := ctx.Query("event")
	if event == "" {
		ctx.JSON(400, gin.H{"error": "malformed request"})
		return
	}

	windowSr := ctx.DefaultQuery("window", "3600")
	window, err := strconv.Atoi(windowSr)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "malformed request"})
		return
	}

	startTimeStr := ctx.DefaultQuery("startTime", time.Now().Add(-defaultLength).UTC().Format(time.RFC3339))
	endTimeStr := ctx.DefaultQuery("endTime", time.Now().UTC().Format(time.RFC3339))

	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "malformed request"})
		return 
	}

	endTime, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "malformed request"})
		return 
	}

	var records []GroupedRecord

	if window * int(time.Second) == int(24 * time.Hour) {
		records, err = GetGroupedRecords(event, startTime, endTime, "analytics_day", window)
	} else if window * int(time.Second) == int(time.Hour) {
		records, err = GetGroupedRecords(event, startTime, endTime, "analytics_hour", window)
	} else if window * int(time.Second) == int(time.Minute) {
		records, err = GetGroupedRecords(event, startTime, endTime, "analytics_minute", window)
	} else {
		recs, err := GetEventRecords(event, startTime, endTime)
		if err != nil {
			log.Printf("Error getting event records: %v", err)
			ctx.JSON(500, gin.H{"error": "internal server error"})
			return 
		}

		records = GroupRecords(event, recs, window)
	}

	if err != nil {
		log.Printf("Error getting event records: %v", err)
		ctx.JSON(500, gin.H{"error": "internal server error"})
		return 
	}

	ctx.JSON(200, records)
}