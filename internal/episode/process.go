package episode

import (
	"context"
	"fmt"

	"github.com/hardhacker/podwise-cli/internal/api"
)

// ProcessResult holds the processing status for an episode.
// Status values: "waiting", "processing", "done", "not_requested", "failed".
type ProcessResult struct {
	Status   string   `json:"status"`
	Progress *float64 `json:"progress"`
}

type processResponse struct {
	Success bool          `json:"success"`
	Result  ProcessResult `json:"result"`
}

// SubmitProcess calls POST /open/v1/episodes/{seq}/process.
// It submits the episode for AI processing (transcription and analysis) and
// returns the initial processing status. Calling this on an already-processed
// episode returns status "done" with progress 100 without consuming credits.
func SubmitProcess(ctx context.Context, client *api.Client, seq int) (*ProcessResult, error) {
	var resp processResponse
	apiPath := fmt.Sprintf("/open/v1/episodes/%d/process", seq)
	if err := client.Post(ctx, apiPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Result, nil
}

// FetchStatus calls GET /open/v1/episodes/{seq}/status and returns the current
// transcription status and progress. Use this to poll for completion after
// SubmitProcess; it does not consume credits.
func FetchStatus(ctx context.Context, client *api.Client, seq int) (*ProcessResult, error) {
	var resp processResponse
	apiPath := fmt.Sprintf("/open/v1/episodes/%d/status", seq)
	if err := client.Get(ctx, apiPath, nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Result, nil
}
