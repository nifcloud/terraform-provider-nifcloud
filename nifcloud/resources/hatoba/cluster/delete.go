package cluster

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandDeleteClusterInput(d)

	svc := meta.(*client.Client).Hatoba
	req := svc.DeleteClusterRequest(input)

	if _, err := req.Send(ctx); err != nil {
		return diag.FromErr(fmt.Errorf("failed deleting Hatoba cluster error: %s", err))
	}

	if err := svc.WaitUntilClusterDeleted(ctx, expandGetClusterInput(d)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to wait for deletion of Hatoba cluster: %s", err))
	}

	d.SetId("")

	return nil
}
