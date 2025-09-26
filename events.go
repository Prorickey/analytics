package main

import (
	"encoding/json"
	"io"
	"log"

	"github.com/gin-gonic/gin"
)

func PostEvent(ctx *gin.Context) {
	event := ctx.Param("event")
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		ctx.JSON(400, gin.H{"error": "malformed body"})
		return
	}

	if !json.Valid(body) {
		ctx.JSON(400, gin.H{"error": "malformed body"})
		return
	}

	CreateEvent(event, string(body))
	ctx.JSON(200, gin.H{})
}