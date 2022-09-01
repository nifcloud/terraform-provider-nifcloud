package bucket

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/storage"
	"github.com/nifcloud/nifcloud-sdk-go/service/storage/types"
)

func flatten(d *schema.ResourceData, bucket types.Buckets,
	versioningRes *storage.GetBucketVersioningOutput, policyRes *storage.GetBucketPolicyOutput) error {
	if err := d.Set("bucket", bucket.Name); err != nil {
		return err
	}

	if err := d.Set("versioning", flattenVersioning(versioningRes)); err != nil {
		return err
	}

	if policyRes != nil && nifcloud.ToString(policyRes.Policy) != "" {
		if err := d.Set("policy", policyRes.Policy); err != nil {
			return err
		}
	}

	return nil
}

func flattenVersioning(out *storage.GetBucketVersioningOutput) []map[string]interface{} {
	res := map[string]interface{}{}

	if out != nil && nifcloud.ToString(out.Status) == "Enabled" {
		res["enabled"] = true
	} else {
		res["enabled"] = false
	}

	return []map[string]interface{}{res}
}
