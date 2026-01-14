package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/KittyChen913/aws-budget-sentinel-go/internal/checks"
	"github.com/KittyChen913/aws-budget-sentinel-go/internal/discord"
	"github.com/aws/aws-lambda-go/lambda"
)

// 儲存彙總的檢查結果
type Report struct {
	Findings map[string]interface{} `json:"findings"`
}

func handler(ctx context.Context) (Report, error) {
	log.Println("Starting aws-budget-sentinel checks")
	results, err := checks.RunAll(ctx)
	if err != nil {
		log.Println("checks error:", err)
	}

	findings := map[string]interface{}{}
	hasRunningServices := false
	for _, r := range results {
		findings[r.Name] = r.Count
		if r.Count > 0 {
			hasRunningServices = true
		}
	}

	report := Report{Findings: findings}
	_ = json.NewEncoder(log.Writer()).Encode(report)

	// 傳送結果到 Discord（若有設定 webhook URL && 有執行中的服務）
	webhookURL := os.Getenv("DISCORD_WEBHOOK_URL")
	if webhookURL != "" && hasRunningServices {
		log.Println("Sending results to Discord...")
		msg := discord.FormatCheckResults(findings)
		if err := discord.SendWebhook(ctx, webhookURL, msg); err != nil {
			log.Printf("Failed to send Discord webhook: %v", err)
		} else {
			log.Println("Successfully sent to Discord")
		}
	} else if webhookURL == "" {
		log.Println("DISCORD_WEBHOOK_URL not set, skipping Discord notification")
	} else {
		log.Println("No running services found, skipping Discord notification")
	}

	return report, nil
}

func main() {
	lambda.Start(handler)
}
