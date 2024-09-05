package browser

import (
	"context"
	"os"

	"github.com/chromedp/chromedp"
)

const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3830.0 Safari/537.36"

func WithBrowserContext(ctx context.Context) (context.Context, context.CancelFunc) {
	headless := os.Getenv("BROWSER_VISIBILITY") != "Active"
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("mute-audio", true),
		chromedp.Flag("headless", headless),
		chromedp.Flag("user-agent", userAgent),
	)
	return chromedp.NewExecAllocator(ctx, opts...)
}
