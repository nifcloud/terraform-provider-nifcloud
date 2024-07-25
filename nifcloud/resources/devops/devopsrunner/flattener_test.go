package devopsrunner

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/devopsrunner"
	"github.com/nifcloud/nifcloud-sdk-go/service/devopsrunner/types"
	"github.com/stretchr/testify/assert"
)

func TestFlatten(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_name")

	wantRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"name":              "test_name",
		"instance_type":     "c-small",
		"availability_zone": "east-11",
		"concurrent":        1,
		"description":       "test_description",
		"network_id":        "test_id",
		"private_address":   "192.168.1.1/24",
		"public_ip_address": "198.51.100.1",
		"system_id":         "test_id",
	})
	wantRd.SetId("test_name")

	wantNotFoundRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})

	type args struct {
		res *devopsrunner.GetRunnerOutput
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
				res: &devopsrunner.GetRunnerOutput{
					Runner: &types.Runner{
						RunnerName:       nifcloud.String("test_name"),
						InstanceType:     nifcloud.String("c-small"),
						AvailabilityZone: nifcloud.String("east-11"),
						Concurrent:       nifcloud.Int32(1),
						Description:      nifcloud.String("test_description"),
						NetworkConfig: &types.NetworkConfig{
							NetworkId:      nifcloud.String("test_id"),
							PrivateAddress: nifcloud.String("192.168.1.1/24"),
						},
						PublicIpAddress: nifcloud.String("198.51.100.1"),
						SystemId:        nifcloud.String("test_id"),
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
