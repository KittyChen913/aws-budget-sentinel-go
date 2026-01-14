package checks

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	ec2 "github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func init() {
	Register(runEC2)
}

func runEC2(ctx context.Context) ([]Result, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS config: %w", err)
	}

	client := ec2.NewFromConfig(cfg)

	// 呼叫 DescribeInstances 並計算正在執行的 EC2 instances
	out, err := client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{})
	if err != nil {
		return nil, err
	}

	running := 0
	for _, r := range out.Reservations {
		for _, i := range r.Instances {
			if i.State != nil && i.State.Name == ec2types.InstanceStateNameRunning {
				running++
			}
		}
	}

	return []Result{{Name: "ec2_instances_running", Count: running}}, nil
}
