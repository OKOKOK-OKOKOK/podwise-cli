package episode

import (
	"context"
	"fmt"

	"github.com/hardhacker/podwise-cli/internal/api"
	"github.com/hardhacker/podwise-cli/internal/cache"
)

// Segment is one transcript chunk returned by the API.
type Segment struct {
	Time     string  `json:"time"`
	Start    float64 `json:"start,omitempty"`
	End      float64 `json:"end,omitempty"`
	Content  string  `json:"content"`
	Speaker  string  `json:"speaker,omitempty"`
	Language string  `json:"language,omitempty"`
}

type transcriptResponse struct {
	Success bool      `json:"success"`
	Result  []Segment `json:"result"`
}

// FetchTranscripts returns the transcript segments for the given episode seq.
// Results are transparently cached in ~/.cache/podwise/<seq>_transcript.json;
// subsequent calls return the cached copy without hitting the network.
func FetchTranscripts(ctx context.Context, client *api.Client, seq int) ([]Segment, error) {
	const cacheType = "transcript"

	var cached []Segment
	if hit, err := cache.Read(seq, cacheType, &cached); err != nil {
		return nil, fmt.Errorf("cache: %w", err)
	} else if hit {
		return cached, nil
	}

	var resp transcriptResponse
	path := fmt.Sprintf("/open/v1/episodes/%d/transcripts", seq)
	if err := client.Get(ctx, path, nil, &resp); err != nil {
		return nil, err
	}

	if err := cache.Write(seq, cacheType, resp.Result); err != nil {
		// Non-fatal: log but don't fail the command.
		fmt.Printf("warning: could not write cache: %v\n", err)
	}

	return resp.Result, nil
}
