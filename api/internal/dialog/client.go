package dialog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type DialogClient struct {
	baseURL    string
	httpClient *http.Client
}

type SendMessageRequest struct {
	Text string `json:"text"`
}

type Message struct {
	ID        string    `json:"id"`
	From      int       `json:"from"`
	To        int       `json:"to"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

type DialogResponse struct {
	Messages []Message `json:"messages"`
}

func NewDialogClient(baseURL string) *DialogClient {
	return &DialogClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *DialogClient) SendMessage(ctx context.Context, fromUserID, toUserID int, text string) (*Message, error) {
	url := fmt.Sprintf("%s/dialog/%d/send", c.baseURL, toUserID)

	reqBody := SendMessageRequest{Text: text}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-User-Id", fmt.Sprintf("%d", fromUserID))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("dialog service returned status %d", resp.StatusCode)
	}

	var message Message
	if err := json.NewDecoder(resp.Body).Decode(&message); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &message, nil
}

func (c *DialogClient) GetDialog(ctx context.Context, userID1, userID2 int) (*DialogResponse, error) {
	url := fmt.Sprintf("%s/dialog/%d/list", c.baseURL, userID2)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-User-Id", fmt.Sprintf("%d", userID1))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("dialog service returned status %d", resp.StatusCode)
	}

	var dialog DialogResponse
	if err := json.NewDecoder(resp.Body).Decode(&dialog); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &dialog, nil
}
