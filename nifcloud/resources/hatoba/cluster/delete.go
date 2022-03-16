package cluster

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/smithy-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/service/hatoba"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).Hatoba
	deadline, _ := ctx.Deadline()

	getClusterInput := expandGetClusterInput(d)

	if _, err := svc.GetCluster(ctx, getClusterInput); err != nil {
		var awsErr smithy.APIError
		if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameterNotFound.Cluster" {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("failed reading cluster: %s", err))
	}

	if err := hatoba.NewClusterRunningWaiter(svc).Wait(ctx, getClusterInput, time.Until(deadline)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to wait for running of Hatoba cluster: %s", err))
	}

	input := expandDeleteClusterInput(d)

	_, err := svc.DeleteCluster(ctx, input)

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed deleting Hatoba cluster error: %s", err))
	}

	if err := hatoba.NewClusterDeletedWaiter(svc).Wait(ctx, getClusterInput, time.Until(deadline)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to wait for deletion of Hatoba cluster: %s", err))
	}

	d.SetId("")

	return nil
}
