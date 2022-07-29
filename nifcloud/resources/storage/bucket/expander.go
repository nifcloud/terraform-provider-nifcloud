package bucket

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/storage"
	"github.com/nifcloud/nifcloud-sdk-go/service/storage/types"
)

func expandGetBucketVersioningInput(d *schema.ResourceData) *storage.GetBucketVersioningInput {
	input := &storage.GetBucketVersioningInput{
		Bucket: nifcloud.String(d.Id()),
	}
	return input
}

func expandGetBucketPolicyInput(d *schema.ResourceData) *storage.GetBucketPolicyInput {
	input := &storage.GetBucketPolicyInput{
		Bucket: nifcloud.String(d.Id()),
	}
	return input
}

func expandPutBucketInput(d *schema.ResourceData) *storage.PutBucketInput {
	input := &storage.PutBucketInput{
		Bucket: nifcloud.String(d.Get("bucket").(string)),
	}
	return input
}

func expandPutBucketVersioningInput(d *schema.ResourceData) *storage.PutBucketVersioningInput {
	input := &storage.PutBucketVersioningInput{
		Bucket:                  nifcloud.String(d.Id()),
		VersioningConfiguration: expandVersioning(d.Get("versioning").([]interface{})),
	}
	return input
}

func expandVersioning(raw []interface{}) *types.RequestVersioningConfiguration {
	if len(raw) == 0 || raw[0] == nil {
		return &types.RequestVersioningConfiguration{
			Status: types.StatusOfVersioningConfigurationForPutBucketVersioningSuspended,
		}
	}

	config := raw[0].(map[string]interface{})
	if value, ok := config["enabled"]; ok && value.(bool) {
		return &types.RequestVersioningConfiguration{
			Status: types.StatusOfVersioningConfigurationForPutBucketVersioningEnabled,
		}
	}

	return &types.RequestVersioningConfiguration{
		Status: types.StatusOfVersioningConfigurationForPutBucketVersioningSuspended,
	}
}

func expandPutBucketPolicyInput(d *schema.ResourceData) *storage.PutBucketPolicyInput {
	input := &storage.PutBucketPolicyInput{
		Bucket: nifcloud.String(d.Id()),
		Policy: nifcloud.String(d.Get("policy").(string)),
	}
	return input
}

func expandDeleteBucketPolicyInput(d *schema.ResourceData) *storage.DeleteBucketPolicyInput {
	input := &storage.DeleteBucketPolicyInput{
		Bucket: nifcloud.String(d.Id()),
	}
	return input
}

func expandDeleteBucketInput(d *schema.ResourceData) *storage.DeleteBucketInput {
	input := &storage.DeleteBucketInput{
		Bucket: nifcloud.String(d.Id()),
	}
	return input
}
