package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/hardhacker/podwise-cli/internal/api"
	"github.com/hardhacker/podwise-cli/internal/config"
	"github.com/hardhacker/podwise-cli/internal/utils"
	"github.com/spf13/cobra"
)

const (
	defaultAuthPollInterval = 5 * time.Second
	defaultAuthTimeout      = 2 * time.Minute
	maxAuthInitAttempts     = 3
)

var browserOpener = openBrowser

type cliAuthStatus struct {
	Status      string `json:"status"`
	AccessToken string `json:"accessToken"`
}

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authorize the CLI in your browser and save the API key automatically",
	Long: `Authorize the podwise CLI in your browser and save the returned API key
to your local configuration file automatically.

The command opens your default browser, shows a short confirmation code, and
waits until you approve the CLI or the 2-minute confirmation window expires.`,
	Args: cobra.NoArgs,
	RunE: runAuth,
}

func init() {
}

func runAuth(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	client := api.New(cfg.APIBaseURL, "")
	parentCtx := cmd.Context()
	if parentCtx == nil {
		parentCtx = context.Background()
	}
	ctx, cancel := context.WithTimeout(parentCtx, defaultAuthTimeout)
	defer cancel()

	fmt.Println("Initializing browser authorization...")
	confirmCode, err := initCLIAuth(ctx, client)
	if err != nil {
		return err
	}

	authURL, err := cliBrowserAuthURL(cfg.APIBaseURL, confirmCode)
	if err != nil {
		return err
	}

	fmt.Printf("Confirm code: %s\n", confirmCode)
	fmt.Printf("Opening browser: %s\n", authURL)
	if err := browserOpener(authURL); err != nil {
		fmt.Printf("Could not open your browser automatically: %v\n", err)
		fmt.Println("Open the URL above manually and approve the request.")
	}

	fmt.Printf("Waiting for authorization (timeout: %s)...\n", utils.FormatDuration(defaultAuthTimeout))
	accessToken, err := pollCLIAuth(ctx, client, confirmCode, defaultAuthPollInterval)
	if err != nil {
		return err
	}

	cfg.APIKey = accessToken
	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	fmt.Printf("Authorization complete. Run `podwise config show` to verify your configuration.\n")
	return nil
}

func initCLIAuth(ctx context.Context, client *api.Client) (string, error) {
	var confirmCode string
	backoff := time.Second

	for attempt := 1; attempt <= maxAuthInitAttempts; attempt++ {
		err := client.Post(ctx, "/no-auth/cli/auth/init", nil, &confirmCode)
		if err == nil {
			if confirmCode == "" {
				return "", errors.New("authorization init returned an empty confirmation code")
			}
			return confirmCode, nil
		}

		var apiErr *api.APIError
		if !errors.As(err, &apiErr) || apiErr.StatusCode < 500 || apiErr.StatusCode >= 600 || attempt == maxAuthInitAttempts {
			return "", fmt.Errorf("initialize authorization: %w", err)
		}

		fmt.Printf("Authorization init failed (attempt %d/%d): %v; retrying in %s...\n", attempt, maxAuthInitAttempts, err, utils.FormatDuration(backoff))

		timer := time.NewTimer(backoff)
		select {
		case <-ctx.Done():
			timer.Stop()
			return "", authContextError(ctx.Err(), "authorization initialization")
		case <-timer.C:
		}
		backoff *= 2
	}

	return "", errors.New("authorization initialization failed")
}

func pollCLIAuth(ctx context.Context, client *api.Client, confirmCode string, interval time.Duration) (string, error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		status, err := fetchCLIAuthStatus(ctx, client, confirmCode)
		if err == nil {
			switch status.Status {
			case "pending":
			case "authorized":
				if status.AccessToken == "" {
					return "", errors.New("authorization succeeded but no access token was returned")
				}
				return status.AccessToken, nil
			default:
				return "", fmt.Errorf("unexpected authorization status %q", status.Status)
			}
		} else {
			var apiErr *api.APIError
			if errors.As(err, &apiErr) && apiErr.StatusCode == 404 {
				return "", errors.New("authorization code expired after 5 minutes; run `podwise auth` again")
			}
			if ctx.Err() != nil {
				return "", authContextError(ctx.Err(), "authorization")
			}
			fmt.Printf("Authorization poll failed: %v; continuing...\n", err)
		}

		select {
		case <-ctx.Done():
			return "", authContextError(ctx.Err(), "authorization")
		case <-ticker.C:
		}
	}
}

func fetchCLIAuthStatus(ctx context.Context, client *api.Client, confirmCode string) (*cliAuthStatus, error) {
	query := url.Values{}
	query.Set("confirmCode", confirmCode)

	var status cliAuthStatus
	if err := client.Get(ctx, "/no-auth/cli/auth", query, &status); err != nil {
		return nil, err
	}
	return &status, nil
}

func cliBrowserAuthURL(apiBaseURL, confirmCode string) (string, error) {
	u, err := url.Parse(apiBaseURL)
	if err != nil {
		return "", fmt.Errorf("invalid API base URL: %w", err)
	}

	basePath := strings.TrimSuffix(strings.TrimSuffix(u.Path, "/"), "/api")
	pageURL := &url.URL{
		Scheme: u.Scheme,
		Host:   u.Host,
		Path:   strings.TrimRight(basePath, "/") + "/auth/cli",
	}
	if pageURL.Path == "/auth/cli" || pageURL.Path == "auth/cli" {
		pageURL.Path = "/auth/cli"
	}

	query := url.Values{}
	query.Set("confirm_code", confirmCode)
	pageURL.RawQuery = query.Encode()
	return pageURL.String(), nil
}

func openBrowser(target string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", target)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", target)
	default:
		cmd = exec.Command("xdg-open", target)
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("launch browser: %w", err)
	}
	return nil
}

func authContextError(err error, action string) error {
	if errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("%s timed out after %s; run `podwise auth` again", action, utils.FormatDuration(defaultAuthTimeout))
	}
	if errors.Is(err, context.Canceled) {
		return fmt.Errorf("%s canceled", action)
	}
	return err
}
