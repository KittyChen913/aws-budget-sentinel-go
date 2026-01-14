package checks

import (
	"context"
	"log"
)

// Result 檢查結果
type Result struct {
	Name  string      `json:"name"`
	Count int         `json:"count"`
	Data  interface{} `json:"data,omitempty"`
}

// CheckFunc 定義檢查函數類型
type CheckFunc func(ctx context.Context) ([]Result, error)

// 註冊的檢查函數列表
var registeredChecks []CheckFunc

// Register 註冊一個檢查函數
func Register(fn CheckFunc) {
	registeredChecks = append(registeredChecks, fn)
}

func RunAll(ctx context.Context) ([]Result, error) {
	res := []Result{}

	// 自動執行所有已註冊的檢查
	for _, checkFn := range registeredChecks {
		r, err := checkFn(ctx)
		if err != nil {
			log.Printf("Check failed: %v", err)
			continue
		}
		res = append(res, r...)
	}

	return res, nil
}
