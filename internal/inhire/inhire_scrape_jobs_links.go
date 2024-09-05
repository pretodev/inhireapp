package inhire

import (
	"errors"
	"fmt"
	"log"
	"regexp"

	"github.com/pretodev/inhireapp/pkg/google"
)

func scrapeJobsLinks() ([]string, error) {
	q := "site:inhire.app"
	page := 0
	mapLinks := make(map[string]bool)
	for {
		log.Printf("getting job links from google (page: %d)", page)
		htmlContent, err := google.FetchPage(q, page)
		if err != nil {
			return nil, fmt.Errorf("failed get google page: %v", err)
		}
		re := regexp.MustCompile(`https?://([a-zA-Z0-9-]+)\.inhire\.app`)
		matches := re.FindAllStringSubmatch(htmlContent, -1)
		if len(matches) == 0 {
			break
		}
		for _, match := range matches {
			if len(match) > 1 {
				link := match[0]
				mapLinks[link] = true
			}
		}
		page = page + 10
	}

	if len(mapLinks) == 0 {
		return nil, errors.New("not found job links")
	}

	links := make([]string, 0)
	for key := range mapLinks {
		links = append(links, key)
	}
	return links, nil
}
