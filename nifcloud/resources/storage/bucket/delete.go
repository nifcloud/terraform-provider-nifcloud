package bucket

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/smithy-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandDeleteBucketInput(d)
	svc := meta.(*client.Client).Storage
	if _, err := svc.DeleteBucket(ctx, input); err != nil {
		var awsErr smithy.APIError
		if errors.As(err, &awsErr) && awsErr.ErrorCode() == "NoSuchBucket" {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("failed deleting bucket: %s", err))
	}

	d.SetId("")

	return nil
}
