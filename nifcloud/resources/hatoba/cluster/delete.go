package cluster

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).Hatoba

	getClusterInput := expandGetClusterInput(d)

	if _, err := svc.GetClusterRequest(getClusterInput).Send(ctx); err != nil {
		var awsErr awserr.Error
		if errors.As(err, &awsErr) && awsErr.Code() == "Client.InvalidParameterNotFound.Cluster" {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("failed reading cluster: %s", err))
	}

	if err := svc.WaitUntilClusterDeleted(ctx, getClusterInput); err != nil {
		return diag.FromErr(fmt.Errorf("failed to wait for deletion of Hatoba cluster: %s", err))
	}

	input := expandDeleteClusterInput(d)

	req := svc.DeleteClusterRequest(input)

	if _, err := req.Send(ctx); err != nil {
		return diag.FromErr(fmt.Errorf("failed deleting Hatoba cluster error: %s", err))
	}

	if err := svc.WaitUntilClusterDeleted(ctx, getClusterInput); err != nil {
		return diag.FromErr(fmt.Errorf("failed to wait for deletion of Hatoba cluster: %s", err))
	}

	d.SetId("")

	return nil
}
