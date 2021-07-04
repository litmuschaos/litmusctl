package types

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	Type        string `json:"type"`
}

type AuthInput struct {
	Endpoint string
	Username string
	Password string
}

type Credentials struct {
	Username string
	Token    string
	Endpoint string
}
