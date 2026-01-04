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
	for _, r := range results {
		findings[r.Name] = r.Count
	}

	report := Report{Findings: findings}
	_ = json.NewEncoder(log.Writer()).Encode(report)

	// 傳送結果到 Discord（若有設定 webhook URL）
	webhookURL := os.Getenv("DISCORD_WEBHOOK_URL")
	if webhookURL != "" {
		log.Println("Sending results to Discord...")
		msg := discord.FormatCheckResults(findings)
		if err := discord.SendWebhook(ctx, webhookURL, msg); err != nil {
			log.Printf("Failed to send Discord webhook: %v", err)
			// 不中斷執行，繼續回傳結果
		} else {
			log.Println("Successfully sent to Discord")
		}
	} else {
		log.Println("DISCORD_WEBHOOK_URL not set, skipping Discord notification")
	}

	return report, nil
}

func main() {
	lambda.Start(handler)
}
