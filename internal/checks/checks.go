package checks

import (
	"context"
	"fmt"
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

// RunAllChecksWithErrors 執行所有檢查並返回錯誤信息，用於需要診斷錯誤的場景
func RunAllChecksWithErrors(ctx context.Context) ([]Result, map[string]string) {
	res := []Result{}
	errors := make(map[string]string)

	// 自動執行所有已註冊的檢查
	for i, checkFn := range registeredChecks {
		r, err := checkFn(ctx)
		if err != nil {
			checkID := fmt.Sprintf("check_%d", i)
			errors[checkID] = err.Error()
			log.Printf("[%s] failed: %v", checkID, err)
			continue
		}
		res = append(res, r...)
	}

	return res, errors
}
