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

	var lastErr error

	const maxBodySize = 1 << 20 // 1MB

	for attempt := 1; attempt <= 3; attempt++ {

		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("http error: %w", err)
			time.Sleep(time.Duration(attempt) * time.Second)
			continue
		}

		// Read limited body
		bodyBytes, readErr := io.ReadAll(io.LimitReader(resp.Body, maxBodySize))
		errClose := resp.Body.Close()
		if errClose != nil {
			return fmt.Errorf("close response body: %w", errClose)
		}

		if readErr != nil {
			lastErr = fmt.Errorf("read body: %w", readErr)
			time.Sleep(time.Duration(attempt) * time.Second)
			continue
		}

		// 429
		if resp.StatusCode == http.StatusTooManyRequests {
			retryAfter := 1
			if hdr := resp.Header.Get("Retry-After"); hdr != "" {
				if v, err := strconv.Atoi(hdr); err == nil {
					retryAfter = v
				}
			}
			time.Sleep(time.Duration(retryAfter) * time.Second)
			continue
		}

		if resp.StatusCode >= 500 {
			lastErr = fmt.Errorf("server error %d: %s", resp.StatusCode, string(bodyBytes))
			time.Sleep(time.Duration(attempt) * time.Second)
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
