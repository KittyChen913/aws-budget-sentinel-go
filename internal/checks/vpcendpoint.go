package checks

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func init() {
	Register(runVPCEndpoint)
}

func runVPCEndpoint(ctx context.Context) ([]Result, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS config: %w", err)
	}

	client := ec2.NewFromConfig(cfg)

	out, err := client.DescribeVpcEndpoints(ctx, &ec2.DescribeVpcEndpointsInput{})
	if err != nil {
		return nil, err
	}

	total := len(out.VpcEndpoints)
	interfaceCount := 0
	interfaceAvailable := 0

	for _, ve := range out.VpcEndpoints {
		if ve.VpcEndpointType == ec2types.VpcEndpointTypeInterface {
			interfaceCount++
			if string(ve.State) == "available" {
				interfaceAvailable++
			}
		}
	}

	return []Result{
		{Name: "vpc_endpoints_total", Count: total},
		{Name: "vpc_interface_endpoints_total", Count: interfaceCount},
		{Name: "vpc_interface_endpoints_available", Count: interfaceAvailable},
	}, nil
}
