package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/hardhacker/podwise-cli/internal/api"
	"github.com/hardhacker/podwise-cli/internal/config"
	"github.com/hardhacker/podwise-cli/internal/episode"
	"github.com/spf13/cobra"
)

var processNoWait bool
var processPollInterval time.Duration
var processTimeout time.Duration

// podwise process <episode-url>
var processCmd = &cobra.Command{
	Use:   "process <episode-url>",
	Short: "Submit an episode or YouTube video for AI processing",
	Long: `Submit a podcast episode or YouTube video for AI processing (transcription and analysis).

Processing consumes credits from your account. The API is asynchronous —
the request returns immediately and the command polls for status until complete.

Status values:
  waiting     episode is queued and will be picked up shortly
  processing  transcription and AI analysis is in progress
  done        processing is complete; use "podwise get" to fetch results

Use --no-wait to submit without waiting for completion.
Use --timeout to override the maximum wait time (default 30m).`,
	Example: `  podwise process https://podwise.ai/dashboard/episodes/7360326`,
	Args:    cobra.ExactArgs(1),
	RunE:    runProcess,
}

func init() {
	processCmd.Flags().BoolVar(&processNoWait, "no-wait", false, "submit and return immediately without polling for completion")
	processCmd.Flags().DurationVar(&processPollInterval, "interval", 10*time.Second, "how often to poll for status updates (min 10s)")
	processCmd.Flags().DurationVar(&processTimeout, "timeout", 30*time.Minute, "maximum time to wait for processing to complete")
}

func runProcess(cmd *cobra.Command, args []string) error {
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

	if processPollInterval < 10*time.Second {
		processPollInterval = 10 * time.Second
	}

	client := api.New(cfg.APIBaseURL, cfg.APIKey)
	ctx := context.Background()

	fmt.Printf("Submitting episode %d for processing...\n", seq)

	result, err := episode.SubmitProcess(ctx, client, seq)
	if err != nil {
		return err
	}

	var initialProgress float64
	if result.Progress != nil {
		initialProgress = *result.Progress
	}
	printProcessStatus(result, initialProgress)

	if processNoWait || result.Status == "done" {
		if result.Status == "done" {
			printProcessDoneHint(seq)
		}
		return nil
	}

	deadline := time.Now().Add(processTimeout)
	ticker := time.NewTicker(processPollInterval)
	defer ticker.Stop()

	var maxProgress float64
	if result.Progress != nil {
		maxProgress = *result.Progress
	}

	for range ticker.C {
		if time.Now().After(deadline) {
			return fmt.Errorf("timed out after %s waiting for episode %d to finish processing", processTimeout, seq)
		}
		status, err := episode.FetchStatus(ctx, client, seq)
		if err != nil {
			return err
		}
		if status.Progress != nil && *status.Progress > maxProgress {
			maxProgress = *status.Progress
		}
		printProcessStatus(status, maxProgress)
		switch status.Status {
		case "done":
			printProcessDoneHint(seq)
			return nil
		case "failed":
			return fmt.Errorf("processing failed for episode %d", seq)
		}
	}
	return nil
}

// printProcessStatus prints a single status line. maxProgress is the
// highest progress value observed so far across all polls, used to
// suppress any regressive values returned by the API.
func printProcessStatus(r *episode.ProcessResult, maxProgress float64) {
	ts := time.Now().Format("15:04:05")
	switch r.Status {
	case "waiting":
		fmt.Printf("  [%s] → waiting       episode is queued for processing\n", ts)
	case "processing":
		if maxProgress >= 0.0 {
			fmt.Printf("  [%s] → processing    %.0f%% complete\n", ts, maxProgress)
		}
	case "done":
		fmt.Printf("  [%s] ✓ done          processing complete (100%%)\n", ts)
	case "not_requested":
		fmt.Printf("  [%s] → not_requested  transcription has not been requested yet\n", ts)
	case "failed":
		fmt.Printf("  [%s] ✗ failed         transcription failed\n", ts)
	default:
		fmt.Printf("  [%s] ? %s\n", ts, r.Status)
	}
}

func printProcessDoneHint(seq int) {
	fmt.Printf("\nRun \"podwise get transcript https://podwise.ai/dashboard/episodes/%d\" to fetch the transcript.", seq)
	fmt.Printf("\nRun \"podwise get summary https://podwise.ai/dashboard/episodes/%d\" to fetch the summary.", seq)
	fmt.Printf("\nRun \"podwise get --help\" for more results.\n")
}
