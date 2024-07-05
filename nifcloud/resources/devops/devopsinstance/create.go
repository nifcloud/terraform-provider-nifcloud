package devopsinstance

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func createInstance(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).DevOps

	input := expandCreateInstanceInput(d)

	res, err := svc.CreateInstance(ctx, input)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create a DevOps instance: %s", err))
	}

	d.SetId(nifcloud.ToString(res.Instance.InstanceId))

	err = waitUntilInstanceRunning(ctx, d, svc)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to wait for the DevOps instance to become ready: %s", err))
	}

	return updateInstance(ctx, d, meta)
}
