package separateinstancerule

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

const waiterInitialDelayForCreate = 5

func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandNiftyCreateSeparateInstanceRuleInput(d)

	svc := meta.(*client.Client).Computing
	_, err := svc.NiftyCreateSeparateInstanceRule(ctx, input)

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed creating SeparateInstanceRule: %s", err))
	}

	SeparateInstanceRuleName := d.Get("name").(string)
	d.SetId(SeparateInstanceRuleName)

	// lintignore:R018
	time.Sleep(waiterInitialDelayForCreate * time.Second)

	return read(ctx, d, meta)
}
