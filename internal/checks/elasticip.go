package checks

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

func init() {
	Register(runElasticIP)
}

func runElasticIP(ctx context.Context) ([]Result, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS config: %w", err)
	}

	client := ec2.NewFromConfig(cfg)

	// 呼叫 DescribeAddresses 並計算未綁定的 Elastic IP
	out, err := client.DescribeAddresses(ctx, &ec2.DescribeAddressesInput{})
	if err != nil {
		return nil, err
	}

	total := len(out.Addresses)
	unattached := 0
	for _, addr := range out.Addresses {
		if addr.AssociationId == nil || *addr.AssociationId == "" {
			unattached++
		}
	}

	return []Result{
		{Name: "elastic_ips_total", Count: total},
		{Name: "elastic_ips_unattached", Count: unattached},
	}, nil
}
