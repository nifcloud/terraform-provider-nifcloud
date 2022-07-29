package bucket

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/smithy-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/storage/types"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).Storage
	res, err := svc.GetService(ctx, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed reading buckets: %s", err))
	}

	bucket, found := findBucket(res.Buckets, d.Id())
	if !found {
		d.SetId("")
		return nil
	}

	getBucketVersioningInput := expandGetBucketVersioningInput(d)
	versioningRes, err := svc.GetBucketVersioning(ctx, getBucketVersioningInput)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed reading bucket versioning: %s", err))
	}

	getBucketPolicyInput := expandGetBucketPolicyInput(d)
	policyRes, err := svc.GetBucketPolicy(ctx, getBucketPolicyInput)
	if err != nil {
		var awsErr smithy.APIError
		if errors.As(err, &awsErr) && awsErr.ErrorCode() == "NoSuchBucketPolicy" {
			policyRes = nil
		} else {
			return diag.FromErr(fmt.Errorf("failed reading bucket policy: %s", err))
		}
	}

	if err := flatten(d, bucket, versioningRes, policyRes); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func findBucket(buckets []types.Buckets, name string) (types.Buckets, bool) {
	for _, bucket := range buckets {
		if nifcloud.ToString(bucket.Name) == name {
			return bucket, true
		}
	}

	return types.Buckets{}, false
}
