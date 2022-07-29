package bucket

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/storage"
	"github.com/nifcloud/nifcloud-sdk-go/service/storage/types"
	"github.com/stretchr/testify/assert"
)

func TestExpandGetBucketVersioningInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"bucket": "test_bucket",
	})
	rd.SetId("test_bucket")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *storage.GetBucketVersioningInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &storage.GetBucketVersioningInput{
				Bucket: nifcloud.String("test_bucket"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandGetBucketVersioningInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandGetBucketPolicyInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"bucket": "test_bucket",
	})
	rd.SetId("test_bucket")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *storage.GetBucketPolicyInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &storage.GetBucketPolicyInput{
				Bucket: nifcloud.String("test_bucket"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandGetBucketPolicyInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandPutBucketInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"bucket": "test_bucket",
	})
	rd.SetId("test_bucket")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *storage.PutBucketInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &storage.PutBucketInput{
				Bucket: nifcloud.String("test_bucket"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandPutBucketInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandPutBucketVersioningInput(t *testing.T) {
	versioningEnabledRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"bucket": "test_bucket",
		"versioning": []interface{}{map[string]interface{}{
			"enabled": true,
		}},
	})
	versioningEnabledRd.SetId("test_bucket")

	versioningDisabledRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"bucket": "test_bucket",
		"versioning": []interface{}{map[string]interface{}{
			"enabled": false,
		}},
	})
	versioningDisabledRd.SetId("test_bucket")

	versioningUndefinedRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"bucket": "test_bucket",
	})
	versioningUndefinedRd.SetId("test_bucket")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *storage.PutBucketVersioningInput
	}{
		{
			name: "expands the resource data (versioning enabled)",
			args: versioningEnabledRd,
			want: &storage.PutBucketVersioningInput{
				Bucket: nifcloud.String("test_bucket"),
				VersioningConfiguration: &types.RequestVersioningConfiguration{
					Status: types.StatusOfVersioningConfigurationForPutBucketVersioningEnabled,
				},
			},
		},
		{
			name: "expands the resource data (versioning disabled)",
			args: versioningDisabledRd,
			want: &storage.PutBucketVersioningInput{
				Bucket: nifcloud.String("test_bucket"),
				VersioningConfiguration: &types.RequestVersioningConfiguration{
					Status: types.StatusOfVersioningConfigurationForPutBucketVersioningSuspended,
				},
			},
		},
		{
			name: "expands the resource data (versioning undefined)",
			args: versioningUndefinedRd,
			want: &storage.PutBucketVersioningInput{
				Bucket: nifcloud.String("test_bucket"),
				VersioningConfiguration: &types.RequestVersioningConfiguration{
					Status: types.StatusOfVersioningConfigurationForPutBucketVersioningSuspended,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandPutBucketVersioningInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandPutBucketPoliyInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"bucket": "test_bucket",
		"policy": `{"Statement": [{"Effect": "Allow"}]}`,
	})
	rd.SetId("test_bucket")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *storage.PutBucketPolicyInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &storage.PutBucketPolicyInput{
				Bucket: nifcloud.String("test_bucket"),
				Policy: nifcloud.String(`{"Statement": [{"Effect": "Allow"}]}`),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandPutBucketPolicyInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDeleteBucketPolicyInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"bucket": "test_bucket",
	})
	rd.SetId("test_bucket")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *storage.DeleteBucketPolicyInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &storage.DeleteBucketPolicyInput{
				Bucket: nifcloud.String("test_bucket"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDeleteBucketPolicyInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDeleteBucketInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"bucket": "test_bucket",
	})
	rd.SetId("test_bucket")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *storage.DeleteBucketInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &storage.DeleteBucketInput{
				Bucket: nifcloud.String("test_bucket"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDeleteBucketInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}
