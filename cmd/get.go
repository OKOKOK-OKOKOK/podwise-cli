package cmd

import (
	"context"
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/hardhacker/podwise-cli/internal/api"
	"github.com/hardhacker/podwise-cli/internal/config"
	"github.com/hardhacker/podwise-cli/internal/episode"
	"github.com/spf13/cobra"
)

// podwise get <subcommand>
var getCmd = &cobra.Command{
	Use:     "get <subcommand>",
	Short:   "Get AI-processed content for a podcast episode",
	Long:    "Get AI-processed content for a podcast episode from podwise.ai.",
	Example: `podwise get transcript https://podwise.ai/dashboard/episodes/7360326`,
}

// podwise get transcript <episode-url>
var transcriptSeconds bool

var getTranscriptCmd = &cobra.Command{
	Use:     "transcript <episode-url>",
	Short:   "Get the full transcript of a podcast episode",
	Long:    "Get the full transcript of a podcast episode and print it to stdout.",
	Example: `podwise get transcript https://podwise.ai/dashboard/episodes/7360326`,
	Args:    cobra.ExactArgs(1),
	RunE:    runGetTranscript,
}

func init() {
	getTranscriptCmd.Flags().BoolVar(&transcriptSeconds, "seconds", false, "show time as start offset in seconds instead of hh:mm:ss")
	getCmd.AddCommand(getTranscriptCmd)
}

func runGetTranscript(cmd *cobra.Command, args []string) error {
	seq, err := parseSeq(args[0])
	if err != nil {
		return fmt.Errorf("invalid episode: %w", err)
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	if err := config.Validate(cfg); err != nil {
		return err
	}

	client := api.New(cfg.APIBaseURL, cfg.APIKey)
	segments, err := episode.FetchTranscripts(context.Background(), client, seq)
	if err != nil {
		return err
	}

	for _, seg := range segments {
		var timeLabel string
		if transcriptSeconds {
			timeLabel = strconv.FormatFloat(seg.Start/1000, 'f', -1, 64)
		} else {
			timeLabel = seg.Time
		}

		if seg.Speaker != "" {
			fmt.Printf("[%s] - %s: %s\n", timeLabel, seg.Speaker, seg.Content)
		} else {
			fmt.Printf("[%s] - %s\n", timeLabel, seg.Content)
		}
	}
	return nil
}

// parseSeq extracts the integer episode seq from a podwise episode URL.
// Expected format: https://podwise.ai/dashboard/episodes/<seq>
func parseSeq(input string) (int, error) {
	if !strings.HasPrefix(input, "http://") && !strings.HasPrefix(input, "https://") {
		return 0, fmt.Errorf("%q is not a valid episode URL", input)
	}
	raw := path.Base(strings.TrimRight(input, "/"))
	seq, err := strconv.Atoi(raw)
	if err != nil || seq <= 0 {
		return 0, fmt.Errorf("%q does not contain a valid episode ID", input)
	}
	return seq, nil
}
