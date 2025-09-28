package main

import (
	"strings"

	"github.com/gin-gonic/gin"
)

type LoginBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func PostLogin(ctx *gin.Context) {
	var body LoginBody
	err := ctx.ShouldBindBodyWithJSON(&body)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "malformed body"})
		return
	}

	user, err := LoginAccount(body.Username, body.Password)
	if err != nil {
		ctx.JSON(401, gin.H{"error": "invalid username or password"})
		return
	}

	ctx.JSON(200, user)
}

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