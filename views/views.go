package views

import (
	"net/http"

	"bitbucket.org/Sofyan_A/sofyan_ahmad_oauth/services"
	"bitbucket.org/Sofyan_A/sofyan_ahmad_oauth/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func getLoginURL(state string) string {
	return services.Conf.AuthCodeURL(state)
}

// IndexHandler handels /.
func IndexView(c *gin.Context) {
	req := c.Request
	path := req.URL.Host
	http.Redirect(c.Writer, req, path+"/secure", 301)
}

// LoginView handles the login procedure.
func LoginView(c *gin.Context) {
	state := utils.RandToken(32)
	session := sessions.Default(c)
	session.Set("state", state)
	session.Save()

	link := getLoginURL(state)
	registerPath := c.Request.URL.Host + "/register"
	c.HTML(http.StatusOK, "login.tmpl", gin.H{"link": link, "registerPath": registerPath})
}

func RegisterView(c *gin.Context) {
	state := utils.RandToken(32)
	link := getLoginURL(state)
	c.HTML(http.StatusOK, "register.tmpl", gin.H{"link": link})
}

// FieldView is a rudementary View for logged in users.
func UserProfileView(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user-id")
	c.HTML(http.StatusOK, "user-profile.tmpl", gin.H{"user": userID})
}
