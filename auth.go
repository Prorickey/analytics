package main

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddlware(ctx *gin.Context) {
	auth := strings.Split(ctx.GetHeader("Authorization"), ":")
	if len(auth) != 2 {
		ctx.JSON(400, gin.H{"error": "malformed authorization"})
		ctx.Abort()
		return
	}

	id := auth[0]
	token := auth[1]

	if ValidateAuthorization(id, token) {
		ctx.Next()
	} else {
		ctx.JSON(401, gin.H{"error": "invalid authentication"})
		ctx.Abort()
	}
}