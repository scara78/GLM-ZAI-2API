// collector.go — Token collector integrated into the main binary.
// Uses chromedp (headless Chromium) to collect device tokens from chat.z.ai.
// Called from the dashboard via /admin/collect-tokens (SSE stream).
package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

const (
	collectorMaxTokens     = 1250
	collectorDefaultTokens = 750
	collectorMaxRetries    = 3
	collectorTimeout       = 120 * time.Second
	collectorZaiURL        = "https://chat.z.ai"
)

// CollectTokens runs the browser-based token collector.
// progress is called with log lines as they happen (for SSE streaming).
// count = number of tokens to collect; outPath = sqlite db path.
func CollectTokens(ctx context.Context, count int, outPath string, progress func(string)) error {
	if count <= 0 {
		count = collectorDefaultTokens
	}
	if count > collectorMaxTokens {
		count = collectorMaxTokens
	}

	progress(fmt.Sprintf("Starting collection of %d tokens → %s", count, outPath))

	// Find Chromium binary
	chromePath := os.Getenv("CHROME_PATH")
	if chromePath == "" {
		// Common paths in Alpine/Docker
		for _, p := range []string{
			"/usr/bin/chromium-browser",
			"/usr/bin/chromium",
			"/usr/bin/google-chrome",
			"/usr/bin/google-chrome-stable",
		} {
			if _, err := os.Stat(p); err == nil {
				chromePath = p
				break
			}
		}
	}

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-software-rasterizer", true),
		chromedp.Flag("disable-extensions", true),
	)
	if chromePath != "" {
		progress(fmt.Sprintf("Using Chromium: %s", chromePath))
		opts = append(opts, chromedp.ExecPath(chromePath))
	}

	allocCtx, allocCancel := chromedp.NewExecAllocator(ctx, opts...)
	defer allocCancel()

	browserCtx, browserCancel := chromedp.NewContext(allocCtx)
	defer browserCancel()

	// Auto-dismiss JS dialogs
	chromedp.ListenTarget(browserCtx, func(ev interface{}) {
		if _, ok := ev.(*page.EventJavascriptDialogOpening); ok {
			go chromedp.Run(browserCtx, page.HandleJavaScriptDialog(true))
		}
	})

	var lastErr error
	for attempt := 1; attempt <= collectorMaxRetries; attempt++ {
		progress(fmt.Sprintf("Attempt %d/%d", attempt, collectorMaxRetries))
		err := collectorTryCollect(browserCtx, count, outPath, progress)
		if err == nil {
			progress("✓ Collection complete!")
			return nil
		}
		lastErr = err
		progress(fmt.Sprintf("✗ Attempt %d failed: %v", attempt, err))
		if attempt < collectorMaxRetries {
			progress("Retrying...")
			time.Sleep(2 * time.Second)
		}
	}
	return fmt.Errorf("all %d attempts failed: %w", collectorMaxRetries, lastErr)
}

func collectorTryCollect(ctx context.Context, total int, outPath string, progress func(string)) error {
	progress("Navigating to chat.z.ai...")
	if err := chromedp.Run(ctx,
		chromedp.Navigate(collectorZaiURL),
		chromedp.WaitVisible(`#chat-input`, chromedp.ByQuery),
	); err != nil {
		return fmt.Errorf("page load failed: %w", err)
	}
	progress("Page loaded, triggering captcha SDK...")

	if err := chromedp.Run(ctx,
		chromedp.SendKeys(`#chat-input`, "__", chromedp.ByQuery),
		chromedp.Click(`#send-message-button`, chromedp.ByQuery),
	); err != nil {
		return fmt.Errorf("send click failed: %w", err)
	}

	progress("Waiting for z_um token endpoint...")
	waitCtx, waitCancel := context.WithTimeout(ctx, 30*time.Second)
	defer waitCancel()
	if err := collectorWaitForZUM(waitCtx); err != nil {
		return fmt.Errorf("token endpoint not ready: %w", err)
	}
	progress("Token endpoint ready, collecting...")

	collectCtx, collectCancel := context.WithTimeout(ctx, collectorTimeout)
	defer collectCancel()

	t0 := time.Now()

	jsExpr := fmt.Sprintf(`(async () => {
		const out = new Array(%d);
		for (let i = 0; i < %d; i++) {
			const tok = window.z_um.getToken();
			out[i] = (tok && typeof tok.then === 'function') ? await tok : tok;
			if (i %% 50 === 0) await new Promise(r => setTimeout(r, 0));
		}
		return out;
	})()`, total, total)

	var tokens []string
	err := chromedp.Run(collectCtx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			result, exception, err := runtime.Evaluate(jsExpr).
				WithAwaitPromise(true).
				WithReturnByValue(true).
				Do(ctx)
			if err != nil {
				return err
			}
			if exception != nil {
				return fmt.Errorf("JS exception: %s", exception.Text)
			}
			if result == nil || result.Value == nil {
				return fmt.Errorf("empty result from getToken")
			}
			return json.Unmarshal(result.Value, &tokens)
		}),
	)
	if err != nil {
		return fmt.Errorf("token collection failed: %w", err)
	}

	var valid []string
	for _, tok := range tokens {
		tok = strings.TrimSpace(tok)
		if tok != "" && tok != "null" && tok != "undefined" {
			valid = append(valid, tok)
		}
	}

	if len(valid) == 0 {
		return fmt.Errorf("no valid tokens collected")
	}

	elapsed := time.Since(t0).Seconds()
	progress(fmt.Sprintf("Collected %d valid tokens in %.1fs, saving to DB...", len(valid), elapsed))

	if err := collectorSaveTokens(outPath, valid); err != nil {
		return fmt.Errorf("save failed: %w", err)
	}

	info, _ := os.Stat(outPath)
	sizeKB := 0.0
	if info != nil {
		sizeKB = float64(info.Size()) / 1024
	}
	progress(fmt.Sprintf("Saved %d tokens to %s (%.1f KB)", len(valid), outPath, sizeKB))
	return nil
}

func collectorWaitForZUM(ctx context.Context) error {
	for {
		var ready bool
		if err := chromedp.Run(ctx,
			chromedp.Evaluate(`typeof window.z_um !== 'undefined' && typeof window.z_um.getToken === 'function'`, &ready),
		); err != nil {
			return err
		}
		if ready {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(500 * time.Millisecond):
		}
	}
}

func collectorSaveTokens(path string, tokens []string) error {
	// Remove old DB and recreate fresh
	os.Remove(path)

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return err
	}
	defer db.Close()

	if _, err := db.Exec(`CREATE TABLE tokens (id INTEGER PRIMARY KEY, token TEXT);`); err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO tokens (id, token) VALUES (?, ?);")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for i, tok := range tokens {
		if _, err := stmt.Exec(i, tok); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	// Reopen the token store so the server picks up the new tokens immediately
	return nil
}