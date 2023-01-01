package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"ServerLoginAuth/services"
)

type RegisterJSONBody struct {
	FullName     string `json:"fullName"`
	Username     string `json:"username"`
	EmailAddress string `json:"emailAddress"`
	Password     string `json:"password"`
	PhoneNumber  string `json:"phoneNumber"`
}

type LoginJSONBody struct {
	EmailAddress string `json:"emailAddress"`
	Password     string `json:"password"`
}

func Register(ctx *gin.Context) {

	reqBytes, reqErr := io.ReadAll(ctx.Request.Body)

	var requestBody RegisterJSONBody

	if reqErr != nil {
		panic(reqErr)
	}

	if err := json.Unmarshal(reqBytes, &requestBody); err != nil {
		panic(err)
	}

	fmt.Println(requestBody)

	services.RegisterUser(requestBody.FullName, requestBody.Username, requestBody.EmailAddress, requestBody.Password, requestBody.PhoneNumber)

	ctx.JSON(http.StatusCreated, gin.H{"message": "Successfully creted user!"})

}
func Login(ctx *gin.Context) {
	session := sessions.Default(ctx)

	reqBytes, reqErr := io.ReadAll(ctx.Request.Body)

	var requestBody LoginJSONBody

	if reqErr != nil {
		panic(reqErr)
	}

	if err := json.Unmarshal(reqBytes, &requestBody); err != nil {
		panic(err)
	}

	// Validate form input
	if strings.Trim(requestBody.EmailAddress, " ") == "" || strings.Trim(requestBody.Password, " ") == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Parameters can't be empty"})
		return
	}

	// Check for username and password match, usually from a database
	// if username != "1729963" || password != "aimanrahim" {
	// 	ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
	// 	return
	// }
	loginResult := services.TryLogin(requestBody.EmailAddress, requestBody.Password)

	if loginResult.Error {
		ctx.JSON(http.StatusUnauthorized, loginResult)
		return
	}

	// Save the username in the session
	session.Set("user", loginResult.Data.Uid) // In real world usage you'd set this to the users ID
	if err := session.Save(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	ctx.JSON(http.StatusOK, loginResult)
}

func Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Delete("user")
	session.Save()
}

func Balance(ctx *gin.Context) {
	session := sessions.Default(ctx)
	user := session.Get("user")
	ctx.JSON(http.StatusOK, gin.H{"user": user})
}
