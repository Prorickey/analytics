package main

import (
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	MustConnect()
	defer CloseDatabase()

	doomAndDespair()

	// router := gin.New()

	// router.Use(gin.Recovery())
	// router.Use(gin.Logger())

	// router.POST("/event/:event", PostEvent)
	// router.POST("/login", PostLogin)
	// router.GET("/analytics", AuthMiddlware, GetAnalytics)

	// router.Run("0.0.0.0:8080")
}

func doomAndDespair() {
	var wg sync.WaitGroup

	for range 10 {
		wg.Add(1)
		go func() {
			_, err := glob_db.Exec("INSERT INTO analytics(event, timestamp) SELECT 'testEvent', NOW() - (random() * interval '7 days') FROM generate_series(1, 100000);")
			if err != nil {
				log.Fatalf("Error inserting into analytics: %v", err)
			}

			wg.Done()
		}()
	}

	wg.Wait()

	log.Println("Finished inserting")
}

const (
	defaultWindow = 7 * 24 * time.Hour
)

func GetAnalytics(ctx *gin.Context) {
	event := ctx.Query("event")
	if event == "" {
		ctx.JSON(400, gin.H{"error": "malformed request"})
		return
	}

	windowSr := ctx.DefaultQuery("window", "5")
	window, err := strconv.Atoi(windowSr)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "malformed request"})
		return
	}

	startTimeStr := ctx.DefaultQuery("startTime", time.Now().Add(-defaultWindow).UTC().Format(time.RFC3339))
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

	records, err := GetEventRecords(event, startTime, endTime)
	if err != nil {
		log.Printf("Error getting event records: %v", err)
		ctx.JSON(500, gin.H{"error": "internal server error"})
		return 
	}

	log.Printf("Number of records: %d", len(records))

	groupedRecords := GroupRecords(event, records, window)

	ctx.JSON(200, groupedRecords)
}