package nasinstance

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandCreateNASInstanceInput(d)

	svc := meta.(*client.Client).NAS
	res, err := svc.CreateNASInstance(ctx, input)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed creating NAS instance: %s", err))
	}

	d.SetId(nifcloud.ToString(res.NASInstance.NASInstanceIdentifier))

	return update(ctx, d, meta)
}
