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
	HttpClient   *http.Client
}

func Authenticate(ctx context.Context, arlCookie string) (*Session, error) {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Timeout: 20 * time.Second,
		Jar:     jar,
	}

	url := "https://www.deezer.com/ajax/gw-light.php?method=deezer.getUserData&input=3&api_version=1.0&api_token="
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.AddCookie(&http.Cookie{
		Name:  "arl",
		Value: arlCookie,
	})

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var res UserDataResponse
	json.Unmarshal(body, &res)

	if res.Results.User.Id == 0 {
		return nil, fmt.Errorf("invalid arl cookie")
	}

	return &Session{
		APIToken:     res.Results.APIToken,
		LicenseToken: res.Results.User.Options.LicenseToken,
		HttpClient:   client,
	}, nil
}
