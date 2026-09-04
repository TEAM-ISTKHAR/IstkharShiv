package utils

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// YouTubeResult is the small, provider-neutral result shape used by the
// platform adapters.
type YouTubeResult struct {
	ID        string
	Title     string
	Link      string
	Duration  string
	Thumbnail string
}

type ytDLResult struct {
	ID        string  `json:"id"`
	Title     string  `json:"title"`
	URL       string  `json:"webpage_url"`
	Duration  float64 `json:"duration"`
	Thumbnail string  `json:"thumbnail"`
	Entries   []struct {
		ID        string  `json:"id"`
		Title     string  `json:"title"`
		URL       string  `json:"webpage_url"`
		Duration  float64 `json:"duration"`
		Thumbnail string  `json:"thumbnail"`
	} `json:"entries"`
}

// SearchYouTube uses yt-dlp's built-in search provider. Keeping search behind
// one helper prevents the platform adapters from depending on a fragile HTML
// scraper and makes the required runtime dependency explicit.
func SearchYouTube(query string, limit int) []YouTubeResult {
	query = strings.TrimSpace(query)
	if query == "" || limit < 1 {
		return nil
	}

	raw, err := exec.Command(
		"yt-dlp",
		"--dump-single-json",
		"--flat-playlist",
		"--skip-download",
		"--no-warnings",
		fmt.Sprintf("ytsearch%d:%s", limit, query),
	).Output()
	if err != nil {
		return nil
	}

	var payload ytDLResult
	if json.Unmarshal(raw, &payload) != nil {
		return nil
	}

	entries := payload.Entries
	if len(entries) == 0 && payload.ID != "" {
		entries = append(entries, struct {
			ID        string  `json:"id"`
			Title     string  `json:"title"`
			URL       string  `json:"webpage_url"`
			Duration  float64 `json:"duration"`
			Thumbnail string  `json:"thumbnail"`
		}{payload.ID, payload.Title, payload.URL, payload.Duration, payload.Thumbnail})
	}

	results := make([]YouTubeResult, 0, len(entries))
	for _, entry := range entries {
		link := entry.URL
		if link == "" {
			link = "https://www.youtube.com/watch?v=" + entry.ID
		}
		results = append(results, YouTubeResult{
			ID:        entry.ID,
			Title:     entry.Title,
			Link:      link,
			Duration:  formatDuration(entry.Duration),
			Thumbnail: entry.Thumbnail,
		})
	}
	return results
}

func formatDuration(seconds float64) string {
	if seconds <= 0 {
		return "0:00"
	}
	total := int64(seconds)
	minutes := total / 60
	remaining := total % 60
	if minutes >= 60 {
		return strconv.FormatInt(minutes/60, 10) + ":" +
			fmt.Sprintf("%02d:%02d", minutes%60, remaining)
	}
	return fmt.Sprintf("%d:%02d", minutes, remaining)
}
