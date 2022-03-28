package dbparametergroup

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).RDB

	input := expandCreateDBParameterGroupInput(d)

	res, err := svc.CreateDBParameterGroup(ctx, input)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed creating DBParameterGroup: %s", err))
	}

	d.SetId(nifcloud.ToString(res.DBParameterGroup.DBParameterGroupName))

	return update(ctx, d, meta)
}
