package instance

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/smithy-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	describeInstancesInput := expandDescribeInstancesInput(d)
	describeInstanceAttributeInput := expandDescribeInstanceAttributeInputWithDisableAPITermination(d)

	svc := meta.(*client.Client).Computing
	describeInstancesRes, err := svc.DescribeInstances(ctx, describeInstancesInput)
	if err != nil {
		var awsErr smithy.APIError
		if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameterNotFound.Instance" {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("failed reading: %s", err))
	}

	describeInstanceAttribeteRes, err := svc.DescribeInstanceAttribute(ctx, describeInstanceAttributeInput)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed reading: %s", err))
	}

	if err := flatten(d, describeInstancesRes); err != nil {
		return diag.FromErr(err)
	}

	if err := flattenDisableAPITermination(d, describeInstanceAttribeteRes); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
