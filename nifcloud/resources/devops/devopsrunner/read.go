package devopsrunner

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/smithy-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func readRunner(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).DevOpsRunner

	runRes, err := svc.GetRunner(ctx, expandGetRunnerInput(d))
	if err != nil {
		var awsErr smithy.APIError
		if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameterNotFound.Runner" {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("failed to read a DevOps Runner: %s", err))
	}

	if err := flatten(d, runRes); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
