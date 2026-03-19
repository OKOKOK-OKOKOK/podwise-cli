package podcast

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// ParseSeq extracts the integer podcast seq from a podwise podcast URL.
// Expected format: https://podwise.ai/dashboard/podcasts/<seq>
func ParseSeq(input string) (int, error) {
	const hint = "(expected https://podwise.ai/dashboard/podcasts/<id>)"

	u, err := url.Parse(input)
	if err != nil || u.Scheme != "https" || (u.Host != "podwise.ai" && u.Host != "beta.podwise.ai") {
		return 0, fmt.Errorf("%q is not a valid podwise podcast URL %s", input, hint)
	}

	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) != 3 || parts[0] != "dashboard" || parts[1] != "podcasts" || parts[2] == "" {
		return 0, fmt.Errorf("%q is not a valid podwise podcast URL %s", input, hint)
	}

	seq, err := strconv.Atoi(parts[2])
	if err != nil || seq <= 0 {
		return 0, fmt.Errorf("podcast ID %q is not a positive integer %s", parts[2], hint)
	}
	return seq, nil
}

// BuildPodcastURL builds a podwise podcast URL from a sequence number.
func BuildPodcastURL(seq int) string {
	return fmt.Sprintf("https://podwise.ai/dashboard/podcasts/%d", seq)
}
