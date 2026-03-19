package cmd

import (
	"context"
	"fmt"

	"github.com/hardhacker/podwise-cli/internal/api"
	"github.com/hardhacker/podwise-cli/internal/config"
	"github.com/hardhacker/podwise-cli/internal/podcast"
	"github.com/spf13/cobra"
)

// podwise follow <podcast-url>
var followCmd = &cobra.Command{
	Use:   "follow <podcast-url>",
	Short: "Follow a podcast by its Podwise URL",
	Long: `Follow a podcast by its Podwise URL.

The podcast-url must be a Podwise podcast URL, e.g. https://podwise.ai/dashboard/podcasts/386.

Following a podcast you already follow succeeds silently (idempotent).`,
	Example: `  podwise follow https://podwise.ai/dashboard/podcasts/386`,
	Args:    cobra.ExactArgs(1),
	RunE:    runFollow,
}

func runFollow(cmd *cobra.Command, args []string) error {
	seq, err := podcast.ParseSeq(args[0])
	if err != nil {
		return fmt.Errorf("invalid podcast: %w", err)
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	if err := config.Validate(cfg); err != nil {
		return err
	}

	client := api.New(cfg.APIBaseURL, cfg.APIKey)
	if err := podcast.Follow(context.Background(), client, seq); err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Following podcast %s\n", podcast.BuildPodcastURL(seq))
	return nil
}
