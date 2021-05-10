package nassecuritygroup

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandDescribeNASSecurityGroupsInput(d)
	svc := meta.(*client.Client).NAS
	req := svc.DescribeNASSecurityGroupsRequest(input)

	res, err := req.Send(ctx)
	if err != nil {
		var awsErr awserr.Error
		if errors.As(err, &awsErr) && awsErr.Code() == "Client.InvalidParameter.NotFound.NASSecurityGroupName" {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("failed reading NAS security group: %s", err))
	}

	if err := flatten(d, res); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
