package nassecuritygroup

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandDeleteNASSecurityGroupInput(d)
	svc := meta.(*client.Client).NAS
	req := svc.DeleteNASSecurityGroupRequest(input)

	if _, err := req.Send(ctx); err != nil {
		return diag.FromErr(fmt.Errorf("failed deleting NAS security group: %s", err))
	}

	d.SetId("")

	return nil
}
