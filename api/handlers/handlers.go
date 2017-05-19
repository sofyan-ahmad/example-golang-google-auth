package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"bitbucket.org/Sofyan_A/sofyan_ahmad_oauth/database"
	"bitbucket.org/Sofyan_A/sofyan_ahmad_oauth/services"
	"bitbucket.org/Sofyan_A/sofyan_ahmad_oauth/structs"
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

	err := services.SetSession(user.Email, c)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error while saving session. Please try again."})
		return
	}

	c.JSON(http.StatusOK, user)
}

func GoogleAuth(c *gin.Context) {
	// Handle the exchange code to initiate a transport.
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

	err = services.SetSession(u.Email, c)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error while saving session. Please try again."})
		return
	}

	if _, dbError := database.Read(u.Email); dbError != nil {
		_, err = database.Create(&u)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"message": "Error while saving user. Please try again."})
			return
		}
	}

	req := c.Request
	path := req.URL.Host
	http.Redirect(c.Writer, req, path+"/secure", 301)
}

func Register(c *gin.Context) {
	var json structs.User
	c.Bind(&json)

	user := structs.User{
		Sub:           json.Sub,
		Name:          json.Name,
		GivenName:     json.GivenName,
		FamilyName:    json.FamilyName,
		Profile:       json.Profile,
		Picture:       json.Picture,
		Email:         json.Email,
		EmailVerified: json.EmailVerified,
		Gender:        json.Gender,
	}

	if _, dbError := database.Read(user.Email); dbError == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Email already used"})
	} else {
		_, err := database.Create(&user)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"message": "Error while saving user. Please try again."})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "new user created"})
}
