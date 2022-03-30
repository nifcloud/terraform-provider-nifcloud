package domainidentity

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandVerifyDomainIdentityInput(d)

	svc := meta.(*client.Client).ESS
	_, err := svc.VerifyDomainIdentity(ctx, input)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed creating domain: %s", err))
	}

	d.SetId(d.Get("domain").(string))

	return read(ctx, d, meta)
}
