package podcast

import (
	"context"
	"fmt"

	"github.com/hardhacker/podwise-cli/internal/api"
)

// Follow follows the podcast identified by seq. The operation is idempotent —
// following an already-followed podcast succeeds silently.
func Follow(ctx context.Context, client *api.Client, seq int) error {
	path := fmt.Sprintf("/open/v1/podcasts/%d/follow", seq)
	var resp struct {
		Success bool `json:"success"`
	}
	return client.Post(ctx, path, nil, &resp)
}
