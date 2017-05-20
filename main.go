package main

import (
	"flag"

	"bitbucket.org/Sofyan_A/sofyan_ahmad_oauth/api/handlers"
	"bitbucket.org/Sofyan_A/sofyan_ahmad_oauth/database"
	"bitbucket.org/Sofyan_A/sofyan_ahmad_oauth/middleware"
	"bitbucket.org/Sofyan_A/sofyan_ahmad_oauth/utils"
	"bitbucket.org/Sofyan_A/sofyan_ahmad_oauth/views"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/contrib/renders/multitemplate"
	"github.com/gin-gonic/gin"
)

func main() {
	dbUrl := flag.String(
		"dbUrl",
		"root:root@tcp(127.0.0.1:3306)/user_sso",
		"set database url",
	)

	baseUrl := flag.String(
		"baseUrl",
		"http://localhost",
		"set base url and port",
	)

	port := flag.String(
		"port",
		"9090",
		"application port",
	)

	flag.Parse()

	utils.SetConfig(*dbUrl, *baseUrl+":"+*port)

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

	router.HTMLRender = createMyRender()

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

	router.Run("0.0.0.0:" + *port)
}

func createMyRender() multitemplate.Render {
	templates := multitemplate.New()
	templates.AddFromFiles("index",
		"./views/templates/index.tmpl",
		"./views/templates/header.tmpl",
		"./views/templates/footer.tmpl")

	templates.AddFromFiles("login",
		"./views/templates/login.tmpl",
		"./views/templates/header.tmpl",
		"./views/templates/footer.tmpl")

	templates.AddFromFiles("register",
		"./views/templates/register.tmpl",
		"./views/templates/header.tmpl",
		"./views/templates/footer.tmpl")

	templates.AddFromFiles("userProfile",
		"./views/templates/user-profile.tmpl",
		"./views/templates/header.tmpl",
		"./views/templates/footer.tmpl")

	templates.AddFromFiles("error",
		"./views/templates/error.tmpl",
		"./views/templates/header.tmpl",
		"./views/templates/footer.tmpl")

	return templates
}
