// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package aws

import (
	"context"

	stateless "github.com/elastic/inputrunner/input/v2/input-stateless"

	"github.com/elastic/elastic-agent-libs/logp"
	"github.com/elastic/elastic-agent-libs/mapstr"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func collectVPCAssets(ctx context.Context, cfg aws.Config, log *logp.Logger, publisher stateless.Publisher) {
	client := ec2.NewFromConfig(cfg)
	vpcs, err := describeVPCs(ctx, client)
	if err != nil {
		log.Errorf("could not describe VPCs for %s: %v", cfg.Region, err)
		return
	}

	for _, vpc := range vpcs {
		publishAWSAsset(publisher,
			cfg.Region,
			*vpc.OwnerId,
			"aws.vpc",
			*vpc.VpcId,
			nil,
			nil,
			flattenEC2Tags(vpc.Tags),
			mapstr.M{
				"isDefault": vpc.IsDefault,
			},
		)
	}
}

func collectSubnetAssets(ctx context.Context, cfg aws.Config, log *logp.Logger, publisher stateless.Publisher) {
	client := ec2.NewFromConfig(cfg)
	subnets, err := describeSubnets(ctx, client)
	if err != nil {
		log.Errorf("could not describe Subnets for %s: %v", cfg.Region, err)
		return
	}

	for _, subnet := range subnets {
		publishAWSAsset(
			publisher,
			cfg.Region,
			*subnet.OwnerId,
			"aws.subnet",
			*subnet.SubnetId,
			[]string{*subnet.VpcId},
			nil,
			flattenEC2Tags(subnet.Tags),
			mapstr.M{
				"state": string(subnet.State),
			},
		)
	}
}

func describeVPCs(ctx context.Context, client *ec2.Client) ([]types.Vpc, error) {
	vpcs := make([]types.Vpc, 0, 100)
	paginator := ec2.NewDescribeVpcsPaginator(client, &ec2.DescribeVpcsInput{})
	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		vpcs = append(vpcs, resp.Vpcs...)
	}

	return vpcs, nil
}

func describeSubnets(ctx context.Context, client *ec2.Client) ([]types.Subnet, error) {
	subnets := make([]types.Subnet, 0, 100)
	paginator := ec2.NewDescribeSubnetsPaginator(client, &ec2.DescribeSubnetsInput{})
	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		subnets = append(subnets, resp.Subnets...)
	}

	return subnets, nil
}