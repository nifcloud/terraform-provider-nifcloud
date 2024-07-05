package devopsinstance

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/smithy-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func readInstance(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).DevOps

	res, err := svc.GetInstance(ctx, expandGetInstanceInput(d))
	if err != nil {
		var awsErr smithy.APIError
		if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameterNotFound.Instance" {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("failed to read a DevOps instance: %s", err))
	}

	if err := flatten(d, res); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
