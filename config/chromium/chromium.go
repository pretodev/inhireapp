package chromium

import (
	"context"

	"github.com/chromedp/chromedp"
	"github.com/pretodev/inhireapp/config/env"
)

const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3830.0 Safari/537.36"

func NewExecAllocator(ctx context.Context, cfg env.Config) (context.Context, context.CancelFunc) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("mute-audio", true),
		chromedp.Flag("headless", cfg.IsChromiumHeadlessEnabled()),
		chromedp.Flag("user-agent", userAgent),
		chromedp.Flag("no-sandbox", true),
	)
	return chromedp.NewExecAllocator(ctx, opts...)
}
