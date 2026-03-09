package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hardhacker/podwise-cli/internal/api"
	"github.com/hardhacker/podwise-cli/internal/config"
	"github.com/hardhacker/podwise-cli/internal/episode"
	"github.com/spf13/cobra"
)

const defaultSearchLimit = 10

var searchLimit int

// podwise search <query>
var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search for podcast episodes",
	Long:  "Search for podcast episodes across the Podwise database and print results to stdout.",
	Example: `  podwise search "artificial intelligence"
  podwise search "machine learning" --limit 20`,
	Args: cobra.MinimumNArgs(1),
	RunE: runSearch,
}

func init() {
	searchCmd.Flags().IntVar(&searchLimit, "limit", defaultSearchLimit, "maximum number of results to return (max 50)")
}

func runSearch(cmd *cobra.Command, args []string) error {
	query := strings.Join(args, " ")

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	if err := config.Validate(cfg); err != nil {
		return err
	}

	client := api.New(cfg.APIBaseURL, cfg.APIKey)
	result, err := episode.Search(context.Background(), client, query, searchLimit)
	if err != nil {
		return err
	}

	if len(result.Hits) == 0 {
		fmt.Println("(no results found)")
		return nil
	}

	fmt.Printf("# Search: \"%s\"\n\n", query)
	fmt.Printf("**Found:** %d\n\n", len(result.Hits))
	fmt.Println("---")
	for i, hit := range result.Hits {
		publishDate := time.Unix(hit.PublishTime, 0).Format("2006-01-02")
		fmt.Printf("\n## %d. %s\n\n", i+1, hit.Title)
		fmt.Printf("- **Podcast:** %s\n", hit.PodcastName)
		fmt.Printf("- **Published:** %s\n", publishDate)
		fmt.Printf("- **Episode URL:** https://podwise.ai/dashboard/episodes/%d\n", hit.Seq)
		if hit.Content != "" {
			fmt.Printf("\n> %s\n", hit.Content)
		}
		fmt.Println("\n---")
	}
	return nil
}
