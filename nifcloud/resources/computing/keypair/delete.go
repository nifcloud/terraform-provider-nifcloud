package keypair

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

	_, err := svc.DeleteKeyPair(ctx, &computing.DeleteKeyPairInput{
		KeyName: nifcloud.String(d.Id()),
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed deleting: %s", err))
	}

	d.SetId("")
	return nil
}
