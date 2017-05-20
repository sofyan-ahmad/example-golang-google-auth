package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"bitbucket.org/Sofyan_A/sofyan_ahmad_oauth/structs"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	DBUrl   string
	BaseUrl string
	Cred    structs.OAuthCredentials
	Conf    *oauth2.Config
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

	Conf = &oauth2.Config{
		ClientID:     Cred.Cid,
		ClientSecret: Cred.Csecret,
		RedirectURL:  BaseUrl + "api/auth",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
		},
		Endpoint: google.Endpoint,
	}
}
