package instance

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandRunInstancesInput(d)

	svc := meta.(*client.Client).Computing
	res, err := svc.RunInstances(ctx, input)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed creating Instance: %s", err))
	}

	if res == nil || len(res.InstancesSet) == 0 {
		return diag.FromErr(fmt.Errorf("launching source instance: no instances returned in response"))
	}

	instance := res.InstancesSet[0]
	d.SetId(nifcloud.ToString(instance.InstanceId))

	return update(ctx, d, meta)
}
