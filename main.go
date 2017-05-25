package main

import (
	"flag"

	"fmt"

	"github.com/sofyanhadia/example-golang-google-auth/api"
	"github.com/sofyanhadia/example-golang-google-auth/database"
	"github.com/sofyanhadia/example-golang-google-auth/middleware"
	"github.com/sofyanhadia/example-golang-google-auth/utils"
	"github.com/sofyanhadia/example-golang-google-auth/views"
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

	utils.SetConfig(*dbUrl, fmt.Sprintf("%s:%s/", *baseUrl, *port))

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
	router.GET("/register/detail", views.RegisterDetailView)
	router.GET("/user/resetpassword", views.ResetPasswordView)
	router.GET("/user/changepassword/do", views.ChangePasswordView)

	// API
	router.POST("/api/login", api.Login)
	router.GET("/api/logout", api.Logout)
	router.POST("/api/register", api.Register)
	router.GET("/api/googleLogin", api.GoogleAuthLogin)
	router.GET("/api/googleRegister", api.GoogleAuthRegister)
	router.GET("/api/userProfile", api.GetCurrentUser)
	router.PUT("/api/user/current", api.UpdateCurrentUser)
	router.GET("/api/user/resetPassword", api.GenerateResetToken)
	router.GET("/api/user/validateResetToken", api.ValidateResetToken)
	router.POST("/api/user/changePassword", api.ChangePassword)

	authorized := router.Group("/secure")
	authorized.Use(middleware.AuthorizeRequest())
	{
		authorized.GET("/", views.UserProfileView)
	}

	router.Run("0.0.0.0:9090")
}

func createMyRender() multitemplate.Render {
	templates := multitemplate.New()

	templates.AddFromFiles("login",
		"./views/templates/login.tmpl",
		"./views/templates/header.tmpl",
		"./views/templates/footer.tmpl")

	templates.AddFromFiles("register",
		"./views/templates/register.tmpl",
		"./views/templates/header.tmpl",
		"./views/templates/footer.tmpl")

	templates.AddFromFiles("registerDetail",
		"./views/templates/register-detail.tmpl",
		"./views/templates/header.tmpl",
		"./views/templates/footer.tmpl")

	templates.AddFromFiles("userProfile",
		"./views/templates/user-profile.tmpl",
		"./views/templates/header.tmpl",
		"./views/templates/footer.tmpl")

	templates.AddFromFiles("resetPassword",
		"./views/templates/reset-password.tmpl",
		"./views/templates/header.tmpl",
		"./views/templates/footer.tmpl")

	templates.AddFromFiles("changePassword",
		"./views/templates/reset-password-change.tmpl",
		"./views/templates/header.tmpl",
		"./views/templates/footer.tmpl")

	templates.AddFromFiles("error",
		"./views/templates/error.tmpl",
		"./views/templates/header.tmpl",
		"./views/templates/footer.tmpl")

	return templates
}
