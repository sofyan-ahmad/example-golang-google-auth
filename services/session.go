package services

import (
	"github.com/sofyanhadia/example-golang-google-auth/structs"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SetSession(user structs.User, context *gin.Context) error {
	session := sessions.Default(context)

	session.Set("user-id", user.Email)
	session.Set("user-name", user.GivenName)

	err := session.Save()
	return err
}

func ClearSession(context *gin.Context) error {
	session := sessions.Default(context)
	session.Clear()

	err := session.Save()
	return err
}
