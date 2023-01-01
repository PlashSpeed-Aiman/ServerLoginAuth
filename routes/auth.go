package routes

import (
	"ServerLoginAuth/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

func Register(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	userint, err := strconv.Atoi(username)
	if err != nil {

	}
	if model.CheckUserExist(userint) == false {
		r := make(chan error)
		go func() {
			_, err := model.RegisterUser(username, password)
			r <- err
		}()
		err := <-r
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Error on Backend!"})
			return
		} else {
			ctx.JSON(http.StatusOK, gin.H{"message": "Successfully registered user"})
			return
		}
	}
	ctx.JSON(http.StatusBadRequest, gin.H{"message": "User Already Exists"})

}
func Login(ctx *gin.Context) {
	session := sessions.Default(ctx)
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")

	// Validate form input
	if strings.Trim(username, " ") == "" || strings.Trim(password, " ") == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Parameters can't be empty"})
		return
	}

	// Check for username and password match, usually from a database
	userint, err := strconv.Atoi(username)
	if err != nil {

	}
	if model.CheckUser(userint, password) == false {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	// Save the username in the session
	session.Set("user", username) // In real world usage you'd set this to the users ID
	if err := session.Save(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully authenticated user"})
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
