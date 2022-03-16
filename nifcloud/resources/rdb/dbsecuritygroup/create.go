package dbsecuritygroup

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandCreateDBSecurityGroupInput(d)

	svc := meta.(*client.Client).RDB

	res, err := svc.CreateDBSecurityGroup(ctx, input)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed creating db security group: %s", err))
	}

	d.SetId(nifcloud.ToString(res.DBSecurityGroup.DBSecurityGroupName))

	return update(ctx, d, meta)
}
