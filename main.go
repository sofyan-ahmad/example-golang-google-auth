package main

import (
	"flag"

	"bitbucket.org/Sofyan_A/sofyan_ahmad_oauth/api/handlers"
	"bitbucket.org/Sofyan_A/sofyan_ahmad_oauth/database"
	"bitbucket.org/Sofyan_A/sofyan_ahmad_oauth/middleware"
	"bitbucket.org/Sofyan_A/sofyan_ahmad_oauth/views"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func main() {
	dbUrl := flag.String(
		"dbUrl",
		"root:root@tcp(127.0.0.1:3306)/user_sso",
		"set database url",
	)

	flag.Parse()

	router := gin.Default()
	store := sessions.NewCookieStore([]byte("super-secret-key"))

	store.Options(sessions.Options{
		Path:   "/",
		MaxAge: 86400 * 7,
	})

	database.New(*dbUrl)

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(sessions.Sessions("my-session", store))

	router.Static("/css", "./static/css")
	router.Static("/img", "./static/img")
	router.LoadHTMLGlob("./views/templates/*")

	// Views
	router.GET("/", views.IndexView)
	router.GET("/login", views.LoginView)
	router.GET("/register", views.RegisterView)

	// API
	router.POST("/api/login", handlers.Login)
	router.POST("/api/register", handlers.Register)
	router.GET("/api/auth", handlers.GoogleAuth)

	authorized := router.Group("/secure")
	authorized.Use(middleware.AuthorizeRequest())
	{
		authorized.GET("/", views.UserProfileView)
	}

	router.Run("127.0.0.1:9090")
}
