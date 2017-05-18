package services

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	Cred Credentials
	Conf *oauth2.Config
)

func init() {
	file, err := ioutil.ReadFile("./creds.json")
	if err != nil {
		log.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	json.Unmarshal(file, &Cred)

	Conf = &oauth2.Config{
		ClientID:     Cred.Cid,
		ClientSecret: Cred.Csecret,
		RedirectURL:  "http://127.0.0.1:9090/api/auth",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
		},
		Endpoint: google.Endpoint,
	}
}

// Credentials which stores google ids.
type Credentials struct {
	Cid     string `json:"cid"`
	Csecret string `json:"csecret"`
}
