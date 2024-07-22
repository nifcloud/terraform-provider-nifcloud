package multiipaddressgroup

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).Computing

	input := expandDescribeMultiIPAddressGroupsInput(d)
	res, err := svc.DescribeMultiIpAddressGroups(ctx, input)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed reading multi IP address group: %s", err))
	}

	if err := flatten(d, res); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
