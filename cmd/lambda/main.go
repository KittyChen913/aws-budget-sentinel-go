package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/KittyChen913/aws-budget-sentinel-go/internal/checks"
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

	return report, nil
}

func main() {
	lambda.Start(handler)
}
