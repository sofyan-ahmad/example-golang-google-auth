package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"bitbucket.org/Sofyan_A/sofyan_ahmad_oauth/database"
	"bitbucket.org/Sofyan_A/sofyan_ahmad_oauth/services"
	"bitbucket.org/Sofyan_A/sofyan_ahmad_oauth/structs"
	"bitbucket.org/Sofyan_A/sofyan_ahmad_oauth/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func Login(c *gin.Context) {
	var json structs.LoginCredential
	c.Bind(&json)

	loginData := structs.LoginCredential{
		Password: json.Password,
		Email:    json.Email,
	}

	user, dbError := database.Login(loginData)

	if dbError != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email or password"})
		return
	}

	err := services.SetSession(*user, c)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error while saving session. Please try again."})
		return
	}

	c.JSON(http.StatusOK, user)
}

func Logout(c *gin.Context) {
	err := services.ClearSession(c)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error while clearing session. Please try again."})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func GoogleAuthLogin(c *gin.Context) {
	session := sessions.Default(c)
	retrievedState := session.Get("state")
	queryState := c.Request.URL.Query().Get("state")

	if retrievedState != queryState {
		log.Printf("Invalid session state: retrieved: %s; Param: %s", retrievedState, queryState)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid session state."})
		return
	}

	// Handle the exchange code to initiate a transport.
	code := c.Request.URL.Query().Get("code")
	tok, err := utils.ConfLogin.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Login failed. Please try again."})
		return
	}

	client := utils.ConfLogin.Client(oauth2.NoContext, tok)
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

	if _, dbError := database.Read(u.Email); dbError != nil {
		log.Println(dbError)
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found. Please register first."})
		return
	}

	err = services.SetSession(u, c)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error while saving session. Please try again."})
		return
	}

	http.Redirect(c.Writer, c.Request, utils.BaseUrl+"secure", 301)
}

func GoogleAuthRegister(c *gin.Context) {
	session := sessions.Default(c)
	retrievedState := session.Get("state")
	queryState := c.Request.URL.Query().Get("state")

	if retrievedState != queryState {
		log.Printf("Invalid session state: retrieved: %s; Param: %s", retrievedState, queryState)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid session state."})
		return
	}

	// Handle the exchange code to initiate a transport.
	code := c.Request.URL.Query().Get("code")
	tok, err := utils.ConfRegister.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Register failed. Please try again."})
		return
	}

	client := utils.ConfRegister.Client(oauth2.NoContext, tok)
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

	user, dbError := database.Read(u.Email)

	if dbError == nil || user != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User already exist. Please login."})
		return
	}

	_, err = database.Create(&u)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error while saving user. Please try again."})
		return
	}

	err = services.SetSession(u, c)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error while saving session. Please try again."})
		return
	}

	http.Redirect(c.Writer, c.Request, utils.BaseUrl+"register/detail", 301)
}

func Register(c *gin.Context) {
	var json structs.User
	c.Bind(&json)

	user := structs.User{
		Sub:           json.Sub,
		GivenName:     json.GivenName,
		FamilyName:    json.FamilyName,
		Profile:       json.Profile,
		Picture:       json.Picture,
		Email:         json.Email,
		EmailVerified: json.EmailVerified,
		Password:      json.Password,
		Gender:        json.Gender,
	}

	if _, dbError := database.Read(user.Email); dbError == nil {
		log.Println(dbError)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Email already used"})
		return
	} else {
		_, err := database.Create(&user)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"message": "Error while saving user. Please try again."})
			return
		}
	}

	err := services.SetSession(user, c)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error while saving session. Please try again."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "new user created"})
}

func GetCurrentUser(c *gin.Context) {
	session := sessions.Default(c)
	userId := session.Get("user-id")

	user, dbError := database.Read(fmt.Sprintf("%s", userId))

	if dbError != nil {
		log.Println(dbError)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error while fetching current user data. Please try again."})
		return
	}

	c.JSON(http.StatusOK, user)
}

func UpdateCurrentUser(c *gin.Context) {
	session := sessions.Default(c)
	userId := fmt.Sprintf("%s", session.Get("user-id"))

	user, dbError := database.Read(userId)

	if dbError != nil {
		log.Println(dbError)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error while fetching current user data. Please try again."})
		return
	}

	var json structs.User
	c.Bind(&json)

	if len(json.Email) > 0 {
		user.Email = json.Email
	}

	user.GivenName = json.GivenName
	user.FamilyName = json.FamilyName
	user.Address = json.Address
	user.Phone = json.Phone
	user.Gender = json.Gender

	_, dbError = database.Update(user)

	if dbError != nil {
		log.Println(dbError)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error while updating user data. Please try again."})
		return
	}

	c.JSON(http.StatusOK, user)
}
