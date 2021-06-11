package cluster

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandCreateClusterInput(d)

	svc := meta.(*client.Client).Hatoba
	req := svc.CreateClusterRequest(input)

	res, err := req.Send(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed creating Hatoba cluster: %s", err))
	}

	d.SetId(nifcloud.StringValue(res.Cluster.Name))

	return update(ctx, d, meta)
}
