package structs

// Credentials which stores google ids.
type OAuthCredentials struct {
	Cid     string `json:"cid"`
	Csecret string `json:"csecret"`
}
