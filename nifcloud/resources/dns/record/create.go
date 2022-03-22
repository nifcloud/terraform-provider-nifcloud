package record

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandCreateChangeResourceRecordSetsInput(d)

	svc := meta.(*client.Client).DNS
	req := svc.ChangeResourceRecordSetsRequest(input)

	_, err := req.Send(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed creating dns record: %s", err))
	}

	return read(ctx, d, meta)
}
