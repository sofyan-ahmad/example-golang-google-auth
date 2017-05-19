package services

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SetSession(email string, context *gin.Context) error {
	session := sessions.Default(context)
	session.Set("user-id", email)
	err := session.Save()
	return err
}
