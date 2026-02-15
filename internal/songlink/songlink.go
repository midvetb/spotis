package songlink

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

func Link(ctx context.Context, raw string) (string, error) {
	u, err := url.Parse("https://api.song.link/v1-alpha.1/links")
	if err != nil {
		return "", err
	}

	q := u.Query()
	q.Set("url", raw)
	q.Set("userCountry", "US")
	q.Set("songIfSingle", "true")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return "", err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("song.link API returned status %d", resp.StatusCode)
	}

	var body struct {
		LinksByPlatform map[string]struct {
			URL string `json:"url"`
		} `json:"linksByPlatform"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	if p, ok := body.LinksByPlatform["deezer"]; ok && p.URL != "" {
		return p.URL, nil
	}

	return "", errors.New("track not found :(")
}
