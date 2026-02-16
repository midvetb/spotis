package deezer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"time"
)

type UserDataResponse struct {
	Results struct {
		APIToken string `json:"checkForm"`
		User     struct {
			Id      int `json:"USER_ID"`
			Options struct {
				LicenseToken string `json:"license_token"`
			} `json:"OPTIONS"`
		} `json:"USER"`
	} `json:"results"`
}

type Session struct {
	APIToken     string
	LicenseToken string
	Client       *http.Client
}

const url = "https://www.deezer.com/ajax/gw-light.php?method=deezer.getUserData&input=3&api_version=1.0&api_token="

func Authenticate(ctx context.Context, arl string) (*Session, error) {
	jar, _ := cookiejar.New(nil)
	c := &http.Client{Timeout: 20 * time.Second, Jar: jar}

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.AddCookie(&http.Cookie{Name: "arl", Value: arl})

	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	var response UserDataResponse
	json.Unmarshal(body, &response)

	if response.Results.User.Id == 0 {
		return nil, fmt.Errorf("invalid arl cookie")
	}

	return &Session{
		APIToken:     response.Results.APIToken,
		LicenseToken: response.Results.User.Options.LicenseToken,
		Client:       c,
	}, nil
}
