package songlink

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

type TrackInfo struct {
	Artist string `json:"artist"`
	Title  string `json:"title"`
	Cover  string `json:"cover"`
	URL    string `json:"url"`
}

type apiResponse struct {
	LinksByPlatform    map[string]platformLink           `json:"linksByPlatform"`
	EntitiesByUniqueId map[string]entitiesByUniqueIdItem `json:"entitiesByUniqueId"`
}

type platformLink struct {
	EntityUniqueId string `json:"entityUniqueId"`
	URL            string `json:"url"`
}

type entitiesByUniqueIdItem struct {
	Title        string `json:"title,omitempty"`
	ArtistName   string `json:"artistName,omitempty"`
	ThumbnailUrl string `json:"thumbnailUrl,omitempty"`
}

func GetLink(ctx context.Context, raw string) (*TrackInfo, error) {
	u, err := url.Parse("https://api.song.link/v1-alpha.1/links")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("url", raw)
	q.Set("userCountry", "US")
	q.Set("songIfSingle", "true")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("song.link API returned status %d", resp.StatusCode)
	}

	var body apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	priority := []string{
		"deezer", "yandex", "tidal", "youtube", "youtubeMusic",
	}

	var chosen platformLink
	found := false
	for _, p := range priority {
		if pl, ok := body.LinksByPlatform[p]; ok && pl.URL != "" {
			chosen = pl
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("no link found for priority platforms")
	}

	ti := &TrackInfo{
		URL: chosen.URL,
	}

	if ent, ok := body.EntitiesByUniqueId[chosen.EntityUniqueId]; ok {
		ti.Title = ent.Title
		ti.Artist = ent.ArtistName
		ti.Cover = ent.ThumbnailUrl
	}

	return ti, nil
}
