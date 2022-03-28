package dbparametergroup

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/smithy-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/rdb/types"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).RDB

	describeDBParameterGroupsRes, err := svc.DescribeDBParameterGroups(ctx, expandDescribeDBParameterGroupsInput(d))
	if err != nil {
		var awsErr smithy.APIError
		if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameterNotFound.DBParameterGroup" {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("failed reading DBParameterGroup: %s", err))
	}

	marker := ""
	parameters := []types.Parameters{}
	for {
		describeDBParametersRes, err := svc.DescribeDBParameters(ctx, expandDescribeDBParametersInput(d, marker))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed reading DBParameters: %s", err))
		}

		parameters = append(parameters, describeDBParametersRes.Parameters...)

		if describeDBParametersRes.Marker == nil {
			break
		}
		marker = nifcloud.ToString(describeDBParametersRes.Marker)
	}

	if err := flatten(d, describeDBParameterGroupsRes, parameters); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
