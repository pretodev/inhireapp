package inhire

import (
	"context"
	"regexp"
)

type JobInfo struct {
	ID           string `json:"id"`
	PositionName string `json:"position_name"`
	PageURL      string `json:"page_url"`
}

func JobID(pageURL string) string {
	re := regexp.MustCompile(`https://[a-zA-Z0-9-]+\.inhire\.app/vagas/([a-f0-9-]+)/`)
	match := re.FindStringSubmatch(pageURL)
	if len(match) <= 1 {
		return ""
	}
	return match[1]
}

type LinkStore interface {
	SaveLinks(ctx context.Context, link ...string) error
	CachedLinks(ctx context.Context) ([]string, error)
}

type JobStore interface {
	SaveJobs(ctx context.Context, jobs ...JobInfo) error
	CachedJobs(ctx context.Context) ([]JobInfo, error)
}
