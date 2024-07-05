package devopsinstance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/devops"
	"github.com/nifcloud/nifcloud-sdk-go/service/devops/types"
	"github.com/stretchr/testify/assert"
)

func TestFlatten(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_id")

	wantRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"instance_id":                    "test_id",
		"instance_type":                  "c-large",
		"firewall_group_name":            "test_name",
		"parameter_group_name":           "test_name",
		"disk_size":                      100,
		"availability_zone":              "east-11",
		"description":                    "test_description",
		"network_id":                     "test_id",
		"private_address":                "192.168.1.1/24",
		"gitlab_url":                     "test_url",
		"registry_url":                   "test_url",
		"object_storage_account":         "test_account",
		"object_storage_region":          "test_region",
		"lfs_bucket_name":                "test_name",
		"packages_bucket_name":           "test_name",
		"container_registry_bucket_name": "test_name",
		"to":                             "test@mail.com",
	})
	wantRd.SetId("test_id")

	wantNotFoundRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})

	type args struct {
		res *devops.GetInstanceOutput
		d   *schema.ResourceData
	}
	tests := []struct {
		name string
		args args
		want *schema.ResourceData
	}{
		{
			name: "flattens the response",
			args: args{
				d: rd,
				res: &devops.GetInstanceOutput{
					Instance: &types.Instance{
						InstanceId:         nifcloud.String("test_id"),
						InstanceType:       nifcloud.String("c-large"),
						FirewallGroupName:  nifcloud.String("test_name"),
						ParameterGroupName: nifcloud.String("test_name"),
						DiskSize:           nifcloud.Int32(int32(100)),
						AvailabilityZone:   nifcloud.String("east-11"),
						Description:        nifcloud.String("test_description"),
						NetworkConfig: &types.NetworkConfig{
							NetworkId:      nifcloud.String("test_id"),
							PrivateAddress: nifcloud.String("192.168.1.1/24"),
						},
						ObjectStorageConfig: &types.ObjectStorageConfig{
							Account: nifcloud.String("test_account"),
							Region:  nifcloud.String("test_region"),
							BucketUseObjects: &types.BucketUseObjects{
								Lfs:               nifcloud.String("test_name"),
								Packages:          nifcloud.String("test_name"),
								ContainerRegistry: nifcloud.String("test_name"),
							},
						},
						To:          nifcloud.String("test@mail.com"),
						GitlabUrl:   nifcloud.String("test_url"),
						RegistryUrl: nifcloud.String("test_url"),
					},
				},
			},
			want: wantRd,
		},
		{
			name: "flattens the response even when the resource has been removed externally",
			args: args{
				d:   wantNotFoundRd,
				res: nil,
			},
			want: wantNotFoundRd,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := flatten(tt.args.d, tt.args.res)
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
