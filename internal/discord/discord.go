package discord

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// 發送到 Discord Webhook 的訊息結構
type Message struct {
	Content string  `json:"content,omitempty"`
	Embeds  []Embed `json:"embeds,omitempty"`
}

// Discord 訊息中的 embed 部分
type Embed struct {
	Title       string  `json:"title,omitempty"`
	Description string  `json:"description,omitempty"`
	Color       int     `json:"color,omitempty"`
	Fields      []Field `json:"fields,omitempty"`
	Timestamp   string  `json:"timestamp,omitempty"`
}

// Discord Embed 中的欄位
type Field struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

// 傳送訊息到 Discord Webhook
func SendWebhook(ctx context.Context, webhookURL string, msg Message) error {
	if webhookURL == "" {
		return fmt.Errorf("Discord webhook URL is empty")
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal Discord message: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", webhookURL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send Discord webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("Discord webhook returned status %d", resp.StatusCode)
	}

	return nil
}

// 將檢查結果格式轉為 Discord Embed
func FormatCheckResults(findings map[string]interface{}) Message {
	fields := []Field{}

	for name, count := range findings {
		fields = append(fields, Field{
			Name:   name,
			Value:  fmt.Sprintf("%v", count),
			Inline: true,
		})
	}

	embed := Embed{
		Title:       "AWS Budget Sentinel - 檢查結果",
		Description: "目前 AWS 資源使用狀況",
		Color:       3447003,
		Fields:      fields,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}

	return Message{
		Embeds: []Embed{embed},
	}
}
