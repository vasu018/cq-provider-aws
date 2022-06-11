// Code generated by internal/generators/main.go; DO NOT EDIT.

package ec2

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudcontrol"
	"github.com/cloudquery/cq-provider-aws/client"
	"github.com/cloudquery/cq-provider-sdk/provider/diag"
	"github.com/cloudquery/cq-provider-sdk/provider/schema"
	"golang.org/x/sync/errgroup"
)

func Aws_ec2_routetable() *schema.Table {
  return &schema.Table{
    Name: "cloudcontrol_aws_ec2_routetable",
    Description: "aws_ec2_routetable",
    Resolver: fetchaws_ec2_routetable,
    Multiplex: client.ServiceAccountRegionMultiplexer(""),
    DeleteFilter: client.DeleteAccountRegionFilter,
    Columns: []schema.Column {
    	{
				Name: "account_id",
				Type: schema.TypeString,
			},
      {
				Name: "region",
				Type: schema.TypeString,
			},
      
      {
        Name: "RouteTableId",
        Description: `The route table ID.`,
        Type: schema.TypeString,
      },
      {
        Name: "Tags",
        Description: `Any tags assigned to the route table.`,
        Type: schema.TypeJSON,
      },
      {
        Name: "VpcId",
        Description: `The ID of the VPC.`,
        Type: schema.TypeString,
      },
    },
  }
}


func fetchaws_ec2_routetable(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
  config := cloudcontrol.ListResourcesInput{
		TypeName:   aws.String("AWS::EC2::RouteTable"),
		MaxResults: aws.Int32(100),
	}
	c := meta.(*client.Client)
  svc := cloudcontrol.NewFromConfig(c.AWSCfg.Copy(), func(o *cloudcontrol.Options) {
		o.Region = c.Region
	})
	for {
		listResources, err := svc.ListResources(ctx, &config)
		if err != nil {
			return diag.WrapError(err)
		}
		// batchResults := make([]map[string]interface{}, len(listResources.ResourceDescriptions))
		errGroup, ctx := errgroup.WithContext(ctx)
		for _, item := range listResources.ResourceDescriptions {
			// i := i
			it := item
			errGroup.Go(func() error {
				r, err := svc.GetResource(ctx, &cloudcontrol.GetResourceInput{
					Identifier: it.Identifier,
					TypeName:   aws.String("AWS::EC2::RouteTable"),
				})
				if err != nil {
					return diag.WrapError(err)
				}
				var resourceJson map[string]interface{}
				if err := json.Unmarshal([]byte(*r.ResourceDescription.Properties), &resourceJson); err != nil {
					return diag.WrapError(err)
				}
				resourceJson["Arn"] = *it.Identifier
				resourceJson["account_id"] = c.AccountID
        resourceJson["region"] = c.Region
				res <- resourceJson
				return nil
			})
		}
		if err := errGroup.Wait(); err != nil {
			return diag.WrapError(err)
		}
		// res <- batchResults
		if aws.ToString(listResources.NextToken) == "" {
			break
		}
		config.NextToken = listResources.NextToken
	}

	return nil
}
