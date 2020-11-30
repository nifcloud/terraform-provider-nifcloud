package instance

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
	describeInstancesInput := expandDescribeInstancesInput(d)
	describeInstanceAttributeInput := expandDescribeInstanceAttributeInputWithDisableAPITermination(d)

	svc := meta.(*client.Client).Computing

	if d.IsNewResource() {
		err := svc.WaitUntilInstanceRunning(ctx, expandDescribeInstancesInput(d))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed waiting for instance to become ready: %s", err))
		}
	}

	describeInstancesReq := svc.DescribeInstancesRequest(describeInstancesInput)
	describeInstanceAttribeteReq := svc.DescribeInstanceAttributeRequest(describeInstanceAttributeInput)

	describeInstancesRes, err := describeInstancesReq.Send(ctx)
	if err != nil {
		var awsErr awserr.Error
		if errors.As(err, &awsErr) && awsErr.Code() == "Client.InvalidParameterNotFound.Instance" {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("failed reading: %s", err))
	}

	describeInstanceAttribeteRes, err := describeInstanceAttribeteReq.Send(ctx)
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
