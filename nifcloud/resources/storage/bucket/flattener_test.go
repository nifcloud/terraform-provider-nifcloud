package bucket

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/storage"
	"github.com/nifcloud/nifcloud-sdk-go/service/storage/types"
	"github.com/stretchr/testify/assert"
)

func TestFlatten(t *testing.T) {
	versioningEnabledRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"bucket": "test_bucket",
		"versioning": []interface{}{map[string]interface{}{
			"enabled": true,
		}},
		"policy": `{"Statement": [{"Effect": "Allow"}]}`,
	})
	versioningEnabledRd.SetId("test_bucket")

	versioningDisabledRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"bucket": "test_bucket",
		"versioning": []interface{}{map[string]interface{}{
			"enabled": false,
		}},
	})
	versioningDisabledRd.SetId("test_bucket")

	type args struct {
		bucket        types.Buckets
		versioningRes *storage.GetBucketVersioningOutput
		policyRes     *storage.GetBucketPolicyOutput
		d             *schema.ResourceData
	}
	tests := []struct {
		name string
		args args
		want *schema.ResourceData
	}{
		{
			name: "flattens the versioning enabled response",
			args: args{
				d: versioningEnabledRd,
				bucket: types.Buckets{
					Name: nifcloud.String("test_bucket"),
				},
				versioningRes: &storage.GetBucketVersioningOutput{
					Status: nifcloud.String("Enabled"),
				},
				policyRes: &storage.GetBucketPolicyOutput{
					Policy: nifcloud.String(`{"Statement": [{"Effect": "Allow"}]}`),
				},
			},
			want: versioningEnabledRd,
		},
		{
			name: "flattens the versioning disabled response",
			args: args{
				d: versioningDisabledRd,
				bucket: types.Buckets{
					Name: nifcloud.String("test_bucket"),
				},
				versioningRes: &storage.GetBucketVersioningOutput{
					Status: nifcloud.String("Suspended"),
				},
				policyRes: &storage.GetBucketPolicyOutput{},
			},
			want: versioningDisabledRd,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := flatten(tt.args.d, tt.args.bucket, tt.args.versioningRes, tt.args.policyRes)
			assert.NoError(t, err)

			wantState := tt.want.State()
			if wantState == nil {
				tt.want.SetId("some")
				wantState = tt.want.State()
			}

			gotState := tt.args.d.State()
			if gotState == nil {
				tt.args.d.SetId("some")
				gotState = tt.args.d.State()
			}

			assert.Equal(t, wantState.Attributes, gotState.Attributes)
		})
	}
}
