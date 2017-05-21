package structs

// User is a retrieved and authentiacted user.
type User struct {
	Id            string `json:"id"`
	Sub           string `json:"sub"`
	GivenName     string `json:"givenName"`
	FamilyName    string `json:"familyName"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"emailVerified"`
	Password      string `json:"password"`
	Gender        string `json:"gender"`
	Phone         string `json:"phone"`
	Address       string `json:"address"`
}
