package views

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"bitbucket.org/Sofyan_A/sofyan_ahmad_oauth/services"
)

// RandToken generates a random @l length token.
func RandToken(l int) string {
	b := make([]byte, l)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func getLoginURL(state string) string {
	return services.Conf.AuthCodeURL(state)
}

// IndexHandler handels /.
func IndexView(c *gin.Context) {
	session := sessions.Default(c)

	if session.Get("user-id") != nil {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{})
	} else {
		req := c.Request
		path := req.URL.Host
		http.Redirect(c.Writer, req, path+"/login", 301)
	}
}

// LoginView handles the login procedure.
func LoginView(c *gin.Context) {
	state := RandToken(32)
	session := sessions.Default(c)
	session.Set("state", state)
	session.Save()
	link := getLoginURL(state)
	c.HTML(http.StatusOK, "login.tmpl", gin.H{"link": link})
}

func RegisterView(c *gin.Context) {
	state := RandToken(32)
	session := sessions.Default(c)
	session.Set("state", state)
	session.Save()
	link := getLoginURL(state)
	c.HTML(http.StatusOK, "register.tmpl", gin.H{"link": link})
}

// FieldView is a rudementary View for logged in users.
func FieldView(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user-id")
	c.HTML(http.StatusOK, "field.tmpl", gin.H{"user": userID})
}
