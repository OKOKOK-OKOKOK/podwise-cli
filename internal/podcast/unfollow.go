package podcast

import (
	"context"
	"fmt"

	"github.com/hardhacker/podwise-cli/internal/api"
)

// Unfollow unfollows the podcast identified by seq. The operation is idempotent —
// unfollowing a podcast you do not follow succeeds silently.
func Unfollow(ctx context.Context, client *api.Client, seq int) error {
	path := fmt.Sprintf("/open/v1/podcasts/%d/unfollow", seq)
	var resp struct {
		Success bool `json:"success"`
	}
	return client.Post(ctx, path, nil, &resp)
}
