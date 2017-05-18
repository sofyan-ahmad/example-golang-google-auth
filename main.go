package main

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"bitbucket.org/Sofyan_A/sofyan_ahmad_oauth/api/handlers"
	"bitbucket.org/Sofyan_A/sofyan_ahmad_oauth/database"
	"bitbucket.org/Sofyan_A/sofyan_ahmad_oauth/middleware"
	views "bitbucket.org/Sofyan_A/sofyan_ahmad_oauth/views"
)

func main() {
	router := gin.Default()
	store := sessions.NewCookieStore([]byte(handlers.RandToken(64)))
	store.Options(sessions.Options{
		Path:   "/",
		MaxAge: 86400 * 7,
	})

	database.New("root:root@tcp(127.0.0.1:3306)/cyza")

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(sessions.Sessions("goquestsession", store))
	router.Static("/css", "./static/css")
	router.Static("/img", "./static/img")
	router.LoadHTMLGlob("./views/templates/*")

	// Views
	router.GET("/", views.IndexView)
	router.GET("/login", views.LoginView)
	router.GET("/register", views.RegisterView)
	// router.GET("/profile", handlers.FieldHandler)

	// API
	router.GET("/api/auth", handlers.AuthHandler)

	authorized := router.Group("/battle")
	authorized.Use(middleware.AuthorizeRequest())
	{
		authorized.GET("/field", handlers.FieldHandler)
	}

	router.Run("127.0.0.1:9090")
}
