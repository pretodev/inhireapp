package inhire

import (
	"context"

	"github.com/chromedp/chromedp"
)

const (
	searchJobEl = "div[data-sentry-component='SearchJob']"
)

const extractJobInfosScript = `(() => {
	const btns = document.querySelectorAll("a[data-sentry-element='NavLink']");
	const jobs = [];
	for (const btn of btns) {
		const pageUrl = btn.href;
		const jobBtnEl = btn.querySelector("div[data-sentry-element='JobPositionName']");
		const positionName = jobBtnEl.innerText;
		jobs.push({pageUrl, positionName});
	}
	return jobs;
})()`

type jsobject map[string]interface{}

func scrapeJobInfos(ctx context.Context, url string) ([]jsobject, error) {
	var jobInfosList []jsobject
	err := chromedp.Run(
		ctx,
		chromedp.Navigate(url+"/vagas"),
		chromedp.WaitReady(searchJobEl),
		chromedp.EvaluateAsDevTools(extractJobInfosScript, &jobInfosList),
	)
	return jobInfosList, err
}
