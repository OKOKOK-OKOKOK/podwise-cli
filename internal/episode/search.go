package episode

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/hardhacker/podwise-cli/internal/api"
)

// SearchHit is a single episode result returned by the search API.
type SearchHit struct {
	Seq         int    `json:"seq"`
	Title       string `json:"title"`
	PodcastName string `json:"podcastName"`
	Content     string `json:"content"`
	EpisodeID   string `json:"episodeId"`
	PodcastID   string `json:"podcastId"`
	PublishTime int64  `json:"publishTime"`
	Cover       string `json:"cover"`
}

// SearchResult holds the full search response from the API.
type SearchResult struct {
	Hits               []SearchHit `json:"result"`
	EstimatedTotalHits int         `json:"estimatedTotalHits"`
	Page               int         `json:"page"`
	HitsPerPage        int         `json:"hitsPerPage"`
}

// Search queries the Podwise episode search API and returns the first page of results.
// limit is passed as hitsPerPage; the API is always queried at page 0.
func Search(ctx context.Context, client *api.Client, query string, limit int) (*SearchResult, error) {
	if query == "" {
		return nil, fmt.Errorf("search query must not be empty")
	}

	q := url.Values{}
	q.Set("q", query)
	q.Set("page", "0")
	q.Set("hitsPerPage", strconv.Itoa(limit))

	var result SearchResult
	if err := client.Get(ctx, "/open/v1/episodes/search", q, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
