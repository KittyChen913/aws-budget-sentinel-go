package checks

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func init() {
	Register(runNATGateway)
}

func runNATGateway(ctx context.Context) ([]Result, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS config: %w", err)
	}

	client := ec2.NewFromConfig(cfg)

	// 呼叫 DescribeNatGateways 並計算可用的 NAT Gateway
	out, err := client.DescribeNatGateways(ctx, &ec2.DescribeNatGatewaysInput{})
	if err != nil {
		return nil, err
	}

	available := 0
	for _, ng := range out.NatGateways {
		if ng.State == ec2types.NatGatewayStateAvailable {
			available++
		}
	}

	return []Result{{Name: "nat_gateways_available", Count: available}}, nil
}
