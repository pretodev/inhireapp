package inhire

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
)

type Service struct {
	store *Store
}

func NewService(
	store *Store,
) *Service {
	return &Service{
		store: store,
	}
}

func (srv *Service) GetJobInfos(ctx context.Context, url string) ([]JobInfo, error) {
	jobInfosList, err := scrapeJobInfos(ctx, url)
	if err != nil {
		return nil, err
	}
	result := make([]JobInfo, 0)
	for _, infos := range jobInfosList {
		pgurl := infos["pageUrl"].(string)
		job := JobInfo{
			ID:           JobID(pgurl),
			PageURL:      pgurl,
			PositionName: infos["positionName"].(string),
		}
		result = append(result, job)
	}
	return result, nil
}

func (srv *Service) UpdateJobInfos(ctx context.Context) error {
	links, err := srv.store.CachedLinks(ctx)
	if err != nil {
		return fmt.Errorf("failed acess cached links: %v", err)
	}
	if len(links) == 0 {
		return errors.New("not found cached links")
	}
	ch := make(chan int, 10)
	var wg sync.WaitGroup
	for _, link := range links {
		wg.Add(1)
		ch <- 1
		go func() {
			defer func() { wg.Done(); <-ch }()
			pageCtx, pageCancel := chromedp.NewContext(ctx)
			defer pageCancel()
			timeoutCtx, timeoutCancel := context.WithTimeout(pageCtx, time.Second*30)
			defer timeoutCancel()
			linkJobs, err := srv.GetJobInfos(timeoutCtx, link)
			if err != nil {
				log.Printf("failed get jobs from link %s: %v", link, err)
				return
			}
			if len(linkJobs) == 0 {
				log.Printf("not found jobs from link %s", link)
				return
			}
			log.Printf("encontrados %d vagas em %s", len(linkJobs), link)
			if err := srv.store.SaveJobs(ctx, linkJobs...); err != nil {
				log.Printf("failed to save jobs from link %s: %v", link, err)
				return
			}
		}()
	}
	wg.Wait()
	return nil
}

func (srv *Service) UpdateJobLinks(ctx context.Context) ([]string, error) {
	links, err := scrapeJobsLinks()
	if err != nil {
		return nil, err
	}
	if err := srv.store.SaveLinks(ctx, links...); err != nil {
		return nil, err
	}
	return links, nil
}
