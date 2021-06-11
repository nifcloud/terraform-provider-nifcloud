package cluster

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).Hatoba

	if d.IsNewResource() {
		err := svc.WaitUntilClusterRunning(ctx, expandGetClusterInput(d))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed waiting for Hatoba cluster to become ready: %s", err))
		}

		return nil
	}

	if d.HasChanges("name", "description", "kubernetes_version", "addons_config") {
		input := expandUpdateClusterInput(d)
		req := svc.UpdateClusterRequest(input)
		if _, err := req.Send(ctx); err != nil {
			return diag.FromErr(fmt.Errorf("failed updating Hatoba cluster: %s", err))
		}

		d.SetId(d.Get("name").(string))
	}

	if d.HasChange("node_pools") {
		o, n := d.GetChange("node_pools")
		toDelete := o.(*schema.Set).Difference(n.(*schema.Set))
		toAdd := n.(*schema.Set).Difference(o.(*schema.Set))
		toChangeSize := n.(*schema.Set).Intersection(o.(*schema.Set))

		for _, l := range toChangeSize.List() {
			change := l.(map[string]interface{})
			input := expandSetNodePoolSizeInput(d, change)
			req := svc.SetNodePoolSizeRequest(input)
			if _, err := req.Send(ctx); err != nil {
				return diag.Errorf(err.Error())
			}

			if err := svc.WaitUntilClusterRunning(ctx, expandGetClusterInput(d)); err != nil {
				return diag.FromErr(fmt.Errorf("failed wait Hatoba cluster available: %s", err))
			}
		}

		toDeleteNames := []string{}
		for _, l := range toDelete.List() {
			del := l.(map[string]interface{})
			toDeleteNames = append(toDeleteNames, del["name"].(string))
		}

		if len(toDeleteNames) != 0 {
			deleteNodePoolsInput := expandDeleteNodePoolsInput(d, toDeleteNames)
			req := svc.DeleteNodePoolsRequest(deleteNodePoolsInput)
			if _, err := req.Send(ctx); err != nil {
				return diag.FromErr(fmt.Errorf("failed deleting Hatoba cluster node pools: %s", err))
			}

			if err := svc.WaitUntilClusterRunning(ctx, expandGetClusterInput(d)); err != nil {
				return diag.FromErr(fmt.Errorf("failed wait Hatoba cluster available: %s", err))
			}
		}

		for _, l := range toAdd.List() {
			add := l.(map[string]interface{})
			input := expandCreateNodePoolInput(d, add)
			req := svc.CreateNodePoolRequest(input)
			if _, err := req.Send(ctx); err != nil {
				return diag.FromErr(fmt.Errorf("failed creating Hatoba cluster node pool: %s", err))
			}

			if err := svc.WaitUntilClusterRunning(ctx, expandGetClusterInput(d)); err != nil {
				return diag.FromErr(fmt.Errorf("failed wait Hatoba cluster available: %s", err))
			}
		}
	}

	return read(ctx, d, meta)
}
