package bucket

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).Storage

	if d.HasChange("versioning") {
		input := expandPutBucketVersioningInput(d)
		if _, err := svc.PutBucketVersioning(ctx, input); err != nil {
			return diag.FromErr(fmt.Errorf("failed updating bucket versioning: %s", err))
		}
	}

	if d.HasChange("policy") {
		if d.Get("policy") == nil || d.Get("policy").(string) == "" {
			input := expandDeleteBucketPolicyInput(d)
			if _, err := svc.DeleteBucketPolicy(ctx, input); err != nil {
				return diag.FromErr(fmt.Errorf("failed deleting bucket policy: %s", err))
			}
		} else {
			input := expandPutBucketPolicyInput(d)
			if _, err := svc.PutBucketPolicy(ctx, input); err != nil {
				return diag.FromErr(fmt.Errorf("failed updating bucket policy: %s", err))
			}
		}
	}

	return read(ctx, d, meta)
}
