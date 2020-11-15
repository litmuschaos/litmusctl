package cmd

import (
	"fmt"
	"net/url"
	"os"

	"github.com/go-resty/resty/v2"
)

type Cred struct {
	Username string
	Token    string
}

type Credentials struct {
	Host     *url.URL
	Username string
	Password []byte
}

type Token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type AuthError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// getToken fetches JWT token for the entered user credentials
func getToken(c Credentials) (Token, AuthError) {

	var authErr AuthError
	client := resty.New()
	token := Token{}
	bodyData := map[string]interface{}{
		"username": fmt.Sprintf("%s", c.Username),
		"password": fmt.Sprintf("%s", c.Password),
	}
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(bodyData).
		//SetResult automatic unmarshalling for the request,
		// if response status code is between 200 and 299
		SetResult(&token).
		SetError(&authErr).
		Post(
			fmt.Sprintf(
				"%s/api/auth/v1/token",
				c.Host,
			),
		)

	if err != nil || !resp.IsSuccess() {
		return Token{}, authErr
	}

	return token, AuthError{}
}

func login(c Credentials) Token {
	t, err := getToken(c)
	if err.Error != "" || t.AccessToken == "" {
		fmt.Println("\nError: ", err.ErrorDescription, ":", err.Error)
		fmt.Println("❌ Login Failed!!")
		os.Exit(1)
	}
	fmt.Println("\n✅ Login Successful!")
	return t
}
