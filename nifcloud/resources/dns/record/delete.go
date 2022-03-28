package record

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandDeleteChangeResourceRecordSetsInput(d)

	svc := meta.(*client.Client).DNS

	if _, err := svc.ChangeResourceRecordSets(ctx, input); err != nil {
		return diag.FromErr(fmt.Errorf("failed deleting dns record error: %s", err))
	}

	d.SetId("")

	return nil
}
