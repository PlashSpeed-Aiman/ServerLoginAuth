package main

import (
	"ServerLoginAuth/model/DB"
	"ServerLoginAuth/routes"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"net/http"
)

func main() {
	DB.DBConn = DB.Setup()
	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("SESSION_ID", store))
	r.Use(gin.Logger())
	r.POST("/register", routes.Register)
	r.POST("/login", routes.Login)

	authenticated := r.Group("/auth")
	authenticated.Use(AuthRequired)
	{
		authenticated.GET("/logout", routes.Logout)
		authenticated.GET("/balance", routes.Balance)

	}
	r.Run(":8000")
}

func AuthRequired(context *gin.Context) {
	session := sessions.Default(context)
	user := session.Get("user")
	if user == nil {
		// Abort the request with the appropriate error code
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	// Continue down the chain to handler etc
	context.Next()
}
