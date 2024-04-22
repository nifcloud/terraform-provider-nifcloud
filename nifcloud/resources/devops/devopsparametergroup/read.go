package devopsparametergroup

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/smithy-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func readParameterGroup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).DevOps

	getParameterGroupRes, err := svc.GetParameterGroup(ctx, expandGetParameterGroupInput(d))
	if err != nil {
		var awsErr smithy.APIError
		if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameterNotFound.ParameterGroup" {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("failed reading ParameterGroup: %s", err))
	}

	if err := flatten(d, getParameterGroupRes); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
