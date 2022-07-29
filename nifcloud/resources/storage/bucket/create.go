package bucket

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandPutBucketInput(d)

	svc := meta.(*client.Client).Storage
	if _, err := svc.PutBucket(ctx, input); err != nil {
		return diag.FromErr(fmt.Errorf("failed creating bucket: %s", err))
	}

	d.SetId(d.Get("bucket").(string))

	return update(ctx, d, meta)
}
