package igdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// TwitchToken represents the OAuth2 response from Twitch.
type TwitchToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// GameData represents the game metadata from IGDB.
type GameData struct {
	ID               int        `json:"id"`
	Name             string     `json:"name"`
	Summary          string     `json:"summary"`
	FirstReleaseDate int64      `json:"first_release_date"`
	Cover            Cover      `json:"cover"`
	Platforms        []Platform `json:"platforms"`
}

// ReleaseYear returns the 4-digit year from the Unix timestamp, or "N/A".
func (g GameData) ReleaseYear() string {
	if g.FirstReleaseDate == 0 {
		return "N/A"
	}
	return time.Unix(g.FirstReleaseDate, 0).UTC().Format("2006")
}

// PlatformNames returns a comma-separated list of platform abbreviations.
func (g GameData) PlatformNames() string {
	if len(g.Platforms) == 0 {
		return "N/A"
	}
	names := make([]string, len(g.Platforms))
	for i, p := range g.Platforms {
		names[i] = p.Abbreviation
		if names[i] == "" {
			names[i] = p.Name
		}
	}
	return strings.Join(names, ", ")
}

// Cover represents the cover metadata from IGDB.
type Cover struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

// BigCoverURL returns the cover URL resized to t_cover_big (264x374).
func (c Cover) BigCoverURL() string {
	if c.URL == "" {
		return ""
	}
	return strings.Replace(c.URL, "t_thumb", "t_cover_big", 1)
}

// Platform represents a game platform from IGDB.
type Platform struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
}

// GetAccessToken retrieves a new access token from Twitch using client credentials.
func GetAccessToken(clientID, clientSecret string) (*TwitchToken, error) {
	url := fmt.Sprintf("https://id.twitch.tv/oauth2/token?client_id=%s&client_secret=%s&grant_type=client_credentials", clientID, clientSecret)

	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to request token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var token TwitchToken
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, fmt.Errorf("failed to decode token: %w", err)
	}

	return &token, nil
}

// SearchGame searches for games by name using the IGDB API.
func SearchGame(clientID, accessToken, query string) ([]GameData, error) {
	url := "https://api.igdb.com/v4/games"

	// IGDB Query Language (Apex)
	// We want ID, Name, Summary, Cover (URL), First Release Date
	q := fmt.Sprintf(`search "%s"; fields name, summary, first_release_date, cover.url, platforms.name, platforms.abbreviation; limit 10;`, query)

	req, err := http.NewRequest("POST", url, bytes.NewBufferString(q))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Client-ID", clientID)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute search: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("search failed with status %d: %s", resp.StatusCode, string(body))
	}

	var games []GameData
	if err := json.NewDecoder(resp.Body).Decode(&games); err != nil {
		return nil, fmt.Errorf("failed to decode games: %w", err)
	}

	return games, nil
}
