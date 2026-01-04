package checks

import (
	"context"
)

// Result 檢查結果
type Result struct {
	Name  string      `json:"name"`
	Count int         `json:"count"`
	Data  interface{} `json:"data,omitempty"`
}

type Checker interface {
	Run(ctx context.Context) ([]Result, error)
}

func RunAll(ctx context.Context) ([]Result, error) {
	res := []Result{}

	// 檢查 EC2
	if r, err := runEC2(ctx); err == nil {
		res = append(res, r...)
	}

	return res, nil
}
