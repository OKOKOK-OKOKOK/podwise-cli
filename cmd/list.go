package cmd

import (
	"context"
	"fmt"

	"github.com/hardhacker/podwise-cli/internal/api"
	"github.com/hardhacker/podwise-cli/internal/config"
	"github.com/hardhacker/podwise-cli/internal/episode"
	"github.com/hardhacker/podwise-cli/internal/podcast"
	"github.com/spf13/cobra"
)

// podwise list <subcommand>
var listCmd = &cobra.Command{
	Use:   "list <subcommand>",
	Short: "List episodes or podcasts from your account",
	Long:  "List episodes or podcasts related to your Podwise account, such as episodes or podcasts you follow.",
	Example: `  podwise list followed-episodes --date today
  podwise list followed-episodes --date yesterday
  podwise list followed-episodes --latest 3 --json
  podwise list followed-podcasts --date today
  podwise list followed-podcasts --latest 14 --json
  podwise list podcast https://podwise.ai/dashboard/podcasts/397
  podwise list podcast https://podwise.ai/dashboard/podcasts/397 --latest 30 --json`,
}

const defaultFollowedLatest = 7

var followedEpisodesDate string
var followedEpisodesLatest int
var followedEpisodesJSONOutput bool

// podwise list followed-episodes
var listFollowedEpisodesCmd = &cobra.Command{
	Use:   "followed-episodes",
	Short: "List recent episodes from podcasts you follow",
	Long: `List episodes published by podcasts the authenticated user follows.

Episodes are sorted by publish time (newest first).

With no flags, shows today's episodes by default.
Use --date to show episodes for a specific day (today, yesterday, or YYYY-MM-DD).
Use --latest N to show the last N days ending today (max 30).
When --date is provided it takes priority and --latest is ignored.`,
	Example: `  podwise list followed-episodes --date today
  podwise list followed-episodes --date yesterday
  podwise list followed-episodes --date 2025-03-01
  podwise list followed-episodes --latest 7 --json`,
	Args: cobra.NoArgs,
	RunE: runListFollowedEpisodes,
}

var followedPodcastsDate string
var followedPodcastsLatest int
var followedPodcastsJSONOutput bool

// podwise list followed-podcasts
var listFollowedPodcastsCmd = &cobra.Command{
	Use:   "followed-podcasts",
	Short: "List followed podcasts with recent new episodes",
	Long: `List podcasts the authenticated user follows that have new episodes within a date range.

Podcasts are sorted by last publish time (newest first).

With no flags, shows podcasts updated today by default.
Use --date to show podcasts updated on a specific day (today, yesterday, or YYYY-MM-DD).
Use --latest N to show the last N days ending today (max 30).
When --date is provided it takes priority and --latest is ignored.`,
	Example: `  podwise list followed-podcasts --date today
  podwise list followed-podcasts --date yesterday
  podwise list followed-podcasts --date 2025-03-01
  podwise list followed-podcasts --latest 14 --json`,
	Args: cobra.NoArgs,
	RunE: runListFollowedPodcasts,
}

const defaultPodcastLatest = 30

var listPodcastLatest int
var listPodcastJSONOutput bool

// podwise list podcast <podcast-url>
var listPodcastCmd = &cobra.Command{
	Use:   "podcast <podcast-url>",
	Short: "List episodes for a specific podcast",
	Long: `List episodes for a specific podcast within a date range, sorted by publish time (newest first).

The podcast-url must be a Podwise podcast URL, e.g. https://podwise.ai/dashboard/podcasts/397.

With no flags, shows episodes from the last 7 days ending today by default.
Use --latest N to look back N days ending today (max 365).`,
	Example: `  podwise list podcast https://podwise.ai/dashboard/podcasts/397
  podwise list podcast https://podwise.ai/dashboard/podcasts/397 --latest 30
  podwise list podcast https://podwise.ai/dashboard/podcasts/397 --latest 30 --json`,
	Args: cobra.ExactArgs(1),
	RunE: runListPodcast,
}

func init() {
	listFollowedEpisodesCmd.Flags().StringVar(&followedEpisodesDate, "date", "", "show episodes for a specific day: today, yesterday, or YYYY-MM-DD (takes priority over --latest)")
	listFollowedEpisodesCmd.Flags().IntVar(&followedEpisodesLatest, "latest", defaultFollowedLatest, "show the last N days ending today (max 30)")
	listFollowedEpisodesCmd.Flags().BoolVar(&followedEpisodesJSONOutput, "json", false, "output results as formatted JSON instead of markdown")
	listCmd.AddCommand(listFollowedEpisodesCmd)

	listFollowedPodcastsCmd.Flags().StringVar(&followedPodcastsDate, "date", "", "show podcasts updated on a specific day: today, yesterday, or YYYY-MM-DD (takes priority over --latest)")
	listFollowedPodcastsCmd.Flags().IntVar(&followedPodcastsLatest, "latest", defaultFollowedLatest, "show podcasts with new episodes in the last N days ending today (max 30)")
	listFollowedPodcastsCmd.Flags().BoolVar(&followedPodcastsJSONOutput, "json", false, "output results as formatted JSON instead of markdown")
	listCmd.AddCommand(listFollowedPodcastsCmd)

	listPodcastCmd.Flags().IntVar(&listPodcastLatest, "latest", defaultPodcastLatest, "show episodes from the last N days ending today (max 365)")
	listPodcastCmd.Flags().BoolVar(&listPodcastJSONOutput, "json", false, "output results as formatted JSON instead of markdown")
	listCmd.AddCommand(listPodcastCmd)
}

func runListFollowedEpisodes(cmd *cobra.Command, args []string) error {
	var date string
	var days int

	if followedEpisodesDate != "" {
		// --date takes priority: show exactly that one day
		parsed, err := episode.ParseDate(followedEpisodesDate)
		if err != nil {
			return err
		}
		date = parsed
		days = 1
	} else if cmd.Flags().Changed("latest") {
		// --latest N explicitly provided: look back N days from today
		if followedEpisodesLatest < 1 || followedEpisodesLatest > 30 {
			return fmt.Errorf("--latest must be between 1 and 30")
		}
		date = episode.Today()
		days = followedEpisodesLatest
	} else {
		// no flags: default to today only
		date = episode.Today()
		days = 1
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	if err := config.Validate(cfg); err != nil {
		return err
	}

	client := api.New(cfg.APIBaseURL, cfg.APIKey)
	result, err := episode.FetchFollowedEpisodes(context.Background(), client, date, days)
	if err != nil {
		return err
	}

	if followedEpisodesJSONOutput {
		data, err := result.FormatJSON()
		if err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), string(data))
		return nil
	}

	fmt.Print(result.FormatText(date, days))
	return nil
}

func runListFollowedPodcasts(cmd *cobra.Command, args []string) error {
	var date string
	var days int

	if followedPodcastsDate != "" {
		parsed, err := episode.ParseDate(followedPodcastsDate)
		if err != nil {
			return err
		}
		date = parsed
		days = 1
	} else if cmd.Flags().Changed("latest") {
		if followedPodcastsLatest < 1 || followedPodcastsLatest > 30 {
			return fmt.Errorf("--latest must be between 1 and 30")
		}
		date = episode.Today()
		days = followedPodcastsLatest
	} else {
		date = episode.Today()
		days = 1
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	if err := config.Validate(cfg); err != nil {
		return err
	}

	client := api.New(cfg.APIBaseURL, cfg.APIKey)
	result, err := podcast.FetchFollowedPodcasts(context.Background(), client, date, days)
	if err != nil {
		return err
	}

	if followedPodcastsJSONOutput {
		data, err := result.FormatJSON()
		if err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), string(data))
		return nil
	}

	fmt.Print(result.FormatText(date, days))
	return nil
}

func runListPodcast(cmd *cobra.Command, args []string) error {
	podcastSeq, err := podcast.ParseSeq(args[0])
	if err != nil {
		return fmt.Errorf("invalid podcast: %w", err)
	}

	if listPodcastLatest < 1 || listPodcastLatest > 365 {
		return fmt.Errorf("--latest must be between 1 and 365")
	}
	date := episode.Today()
	days := listPodcastLatest

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	if err := config.Validate(cfg); err != nil {
		return err
	}

	client := api.New(cfg.APIBaseURL, cfg.APIKey)
	result, err := podcast.FetchPodcastEpisodes(context.Background(), client, podcastSeq, date, days)
	if err != nil {
		return err
	}

	if listPodcastJSONOutput {
		data, err := result.FormatJSON()
		if err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), string(data))
		return nil
	}

	fmt.Print(result.FormatText(date, days))
	return nil
}
