package privatelan

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).Computing
	req := svc.NiftyDeletePrivateLanRequest(&computing.NiftyDeletePrivateLanInput{NetworkId: nifcloud.String(d.Id())})

	_, err := req.Send(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed deleting private_lan: %s", err))
	}

	d.SetId("")
	return nil
}
