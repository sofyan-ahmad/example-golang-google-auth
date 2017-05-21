package views

import (
	"log"
	"net/http"

	"bitbucket.org/Sofyan_A/sofyan_ahmad_oauth/database"
	"bitbucket.org/Sofyan_A/sofyan_ahmad_oauth/services"
	"bitbucket.org/Sofyan_A/sofyan_ahmad_oauth/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func getLoginURL(state string) string {
	return utils.ConfLogin.AuthCodeURL(state)
}

func getRegisterURL(state string) string {
	return utils.ConfRegister.AuthCodeURL(state)
}

// IndexHandler handels /.
func IndexView(c *gin.Context) {
	req := c.Request
	session := sessions.Default(c)
	userId := session.Get("user-id")

	if userId == nil || userId == "" {
		http.Redirect(c.Writer, req, utils.BaseUrl+"login", 307)
	} else {
		http.Redirect(c.Writer, req, utils.BaseUrl+"secure", 307)
	}
}

// LoginView handles the login procedure.
func LoginView(c *gin.Context) {
	state := utils.RandToken(32)
	session := sessions.Default(c)
	session.Set("state", state)
	session.Save()

	link := getLoginURL(state)
	c.HTML(http.StatusOK, "login", gin.H{"baseUrl": utils.BaseUrl, "link": link})
}

func RegisterView(c *gin.Context) {
	state := utils.RandToken(32)
	session := sessions.Default(c)
	session.Set("state", state)
	session.Save()

	link := getRegisterURL(state)
	c.HTML(http.StatusOK, "register", gin.H{"baseUrl": utils.BaseUrl, "link": link})
}

func RegisterDetailView(c *gin.Context) {
	state := utils.RandToken(32)
	link := getLoginURL(state)
	session := sessions.Default(c)
	userId := session.Get("user-id")

	if userId == nil || userId == "" {
		c.HTML(http.StatusBadRequest, "error", gin.H{"message": "Invalid register session. Please try again."})
		return
	}

	c.HTML(http.StatusOK, "registerDetail", gin.H{"baseUrl": utils.BaseUrl, "link": link})
}

// UserProfileView is a View for logged in users.
func UserProfileView(c *gin.Context) {
	session := sessions.Default(c)
	userId := session.Get("user-id")
	userName := session.Get("user-name")

	c.HTML(http.StatusOK, "userProfile", gin.H{"baseUrl": utils.BaseUrl, "userId": userId, "userName": userName})
}

func ResetPasswordView(c *gin.Context) {
	services.ClearSession(c)
	c.HTML(http.StatusOK, "resetPassword", gin.H{"baseUrl": utils.BaseUrl})
}

func ChangePasswordView(c *gin.Context) {
	email := c.Request.URL.Query().Get("email")
	token := c.Request.URL.Query().Get("t")

	user, dbError := database.Read(email)
	if dbError != nil {
		log.Println(dbError)
		c.HTML(http.StatusBadRequest, "error", gin.H{"message": "Error while fetching current user data. Please try again."})
		return
	}

	dbError = database.CheckResetToken(user, token)
	if dbError != nil {
		log.Println(dbError)
		c.HTML(http.StatusBadRequest, "error", gin.H{"message": "Invalid token."})
		return
	}

	services.ClearSession(c)

	session := sessions.Default(c)
	session.Set("reset-token", token)
	session.Set("reset-email", email)
	session.Save()

	c.HTML(http.StatusOK, "changePassword", gin.H{"baseUrl": utils.BaseUrl, "token": token, "email": email})
}
