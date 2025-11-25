package operator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	APIKey     string
}

func NewClient(baseUrl string) *Client {
	return &Client{
		BaseURL: baseUrl,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func jitter(d time.Duration) time.Duration {
	return d + time.Duration(50-time.Now().UnixNano()%100)*time.Millisecond
}

func urlContains(url, part string) bool {
	return len(url) >= len(part) && (url[len(url)-len(part):] == part || url[len(url)-len(part)-1] == '/')
}

func (c *Client) doRequest(method, url string, body any, out any) error {
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return fmt.Errorf("encode body: %w", err)
		}
	}

	req, err := http.NewRequest(method, url, &buf)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	if m, ok := body.(map[string]any); ok {
		if ref, ok := m["refId"].(string); ok {
			if method == "POST" && urlContains(url, "withdraw") {
				req.Header.Set("X-Idempotency-Key", "withdraw-"+ref)
			}
			if method == "POST" && urlContains(url, "deposit") {
				req.Header.Set("X-Idempotency-Key", "deposit-"+ref)
			}
		}
	}

	var lastErr error
	const maxBodySize = 1 << 20 // 1MB
	const maxRetries = 5

	for attempt := 1; attempt <= maxRetries; attempt++ {
		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("http error: %w", err)
			time.Sleep(jitter(time.Duration(attempt) * time.Second))
			continue
		}

		bodyBytes, readErr := io.ReadAll(io.LimitReader(resp.Body, maxBodySize))
		_ = resp.Body.Close()

		if readErr != nil {
			lastErr = fmt.Errorf("read body: %w", readErr)
			time.Sleep(jitter(time.Duration(attempt) * time.Second))
			continue
		}

		// 429
		if resp.StatusCode == http.StatusTooManyRequests {
			retryAfter := 1
			if hdr := resp.Header.Get("Retry-After"); hdr != "" {
				if v, err := strconv.Atoi(hdr); err == nil && v > 0 {
					retryAfter = v
				}
			}
			time.Sleep(time.Duration(retryAfter) * time.Second)
			continue
		}

		if resp.StatusCode >= 500 {
			lastErr = fmt.Errorf("server error %d: %s", resp.StatusCode, string(bodyBytes))
			time.Sleep(jitter(time.Duration(attempt) * time.Second))
			continue
		}

		if resp.StatusCode >= 400 {
			return fmt.Errorf("operator error %d: %s", resp.StatusCode, string(bodyBytes))
		}

		if out != nil {
			if err := json.Unmarshal(bodyBytes, out); err != nil {
				return fmt.Errorf("decode operator response: %w", err)
			}
		}

		return nil
	}

	if lastErr != nil {
		return fmt.Errorf("failed after retries: %w", lastErr)
	}

	return fmt.Errorf("failed after retries without specific error")
}

func (c *Client) Withdraw(playerID string, req WithdrawRequest) (*WithdrawResponse, error) {
	url := fmt.Sprintf("%s/v2/players/%s/withdraw", c.BaseURL, playerID)

	var resp WithdrawResponse
	err := c.doRequest("POST", url, req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) Deposit(playerID string, req DepositRequest) (*DepositResponse, error) {
	url := fmt.Sprintf("%s/v2/players/%s/deposit", c.BaseURL, playerID)

	var resp DepositResponse
	err := c.doRequest("POST", url, req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
