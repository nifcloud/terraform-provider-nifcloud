package dbparametergroup

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/service/rdb"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).RDB

	describeDBParameterGroupsReq := svc.DescribeDBParameterGroupsRequest(expandDescribeDBParameterGroupsInput(d))
	describeDBParameterGroupsRes, err := describeDBParameterGroupsReq.Send(ctx)
	if err != nil {
		var awsErr awserr.Error
		if errors.As(err, &awsErr) && awsErr.Code() == "Client.InvalidParameterNotFound.DBParameterGroup" {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("failed reading DBParameterGroup: %s", err))
	}

	marker := ""
	parameters := []rdb.Parameter{}
	for {
		describeDBParametersReq := svc.DescribeDBParametersRequest(expandDescribeDBParametersInput(d, marker))
		describeDBParametersRes, err := describeDBParametersReq.Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed reading DBParameters: %s", err))
		}
		if describeDBParametersRes.Marker == nil {
			break
		}
		marker = *describeDBParametersRes.Marker
		parameters = append(parameters, describeDBParametersRes.Parameters...)
	}

	if err := flatten(d, describeDBParameterGroupsRes, parameters); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
