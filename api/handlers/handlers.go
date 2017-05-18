package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"bitbucket.org/Sofyan_A/sofyan_ahmad_oauth/database"
	"bitbucket.org/Sofyan_A/sofyan_ahmad_oauth/services"
	"bitbucket.org/Sofyan_A/sofyan_ahmad_oauth/structs"
	"golang.org/x/oauth2"
)

// RandToken generates a random @l length token.
func RandToken(l int) string {
	b := make([]byte, l)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

// AuthHandler handles authentication of a user and initiates a session.
func AuthHandler(c *gin.Context) {
	// Handle the exchange code to initiate a transport.
	session := sessions.Default(c)
	// retrievedState := session.Get("state")
	// queryState := c.Request.URL.Query().Get("state")

	// if retrievedState != queryState {
	// 	log.Printf("Invalid session state: retrieved: %s; Param: %s", retrievedState, queryState)
	// 	c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid session state."})
	// 	return
	// }
	code := c.Request.URL.Query().Get("code")
	tok, err := services.Conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Login failed. Please try again."})
		return
	}

	client := services.Conf.Client(oauth2.NoContext, tok)
	userinfo, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	defer userinfo.Body.Close()
	data, _ := ioutil.ReadAll(userinfo.Body)
	u := structs.User{}
	if err = json.Unmarshal(data, &u); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error marshalling response. Please try again."})
		return
	}
	session.Set("user-id", u.Email)
	err = session.Save()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error while saving session. Please try again."})
		return
	}
	seen := false

	if _, mongoErr := database.Read(u.Email); mongoErr == nil {
		seen = true
	} else {
		_, err = database.Create(&u)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"message": "Error while saving user. Please try again."})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"email": u.Email, "seen": seen})
}

// FieldHandler is a rudementary handler for logged in users.
func FieldHandler(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user-id")
	c.JSON(http.StatusOK, gin.H{"user": userID})
}
