package views

import (
	"net/http"

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

// FieldView is a View for logged in users.
func UserProfileView(c *gin.Context) {
	session := sessions.Default(c)
	userId := session.Get("user-id")
	userName := session.Get("user-name")

	c.HTML(http.StatusOK, "userProfile", gin.H{"baseUrl": utils.BaseUrl, "userId": userId, "userName": userName})
}
