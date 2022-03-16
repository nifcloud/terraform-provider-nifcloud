package dbsecuritygroup

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandDeleteDBSecurityGroupInput(d)
	svc := meta.(*client.Client).RDB

	_, err := svc.DeleteDBSecurityGroup(ctx, input)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed deleting: %s", err))
	}

	d.SetId("")
	return nil
}
