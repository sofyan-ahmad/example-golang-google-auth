package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/sofyanhadia/example-golang-google-auth/structs"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	DBUrl        string
	BaseUrl      string
	Cred         structs.OAuthCredentials
	ConfLogin    *oauth2.Config
	ConfRegister *oauth2.Config
)

func SetConfig(dbUrl string, baseUrl string) {
	file, err := ioutil.ReadFile("./creds.json")

	if err != nil {
		log.Printf("File error: %v\n", err)
		os.Exit(1)
	}

	json.Unmarshal(file, &Cred)

	DBUrl = dbUrl
	BaseUrl = baseUrl

	ConfLogin = generateLoginConf()
	ConfRegister = generateRegisterConf()
}

func generateLoginConf() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     Cred.Cid,
		ClientSecret: Cred.Csecret,
		RedirectURL:  BaseUrl + "api/googleLogin",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
		},
		Endpoint: google.Endpoint,
	}
}

func generateRegisterConf() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     Cred.Cid,
		ClientSecret: Cred.Csecret,
		RedirectURL:  BaseUrl + "api/googleRegister",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
		},
		Endpoint: google.Endpoint,
	}
}
