package devopsrunner

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/devopsrunner"
	"github.com/nifcloud/nifcloud-sdk-go/service/devopsrunner/types"
	"github.com/stretchr/testify/assert"
)

func TestExpandCreateRunnerInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"name":              "test_name",
		"instance_type":     "c-small",
		"availability_zone": "east-11",
		"concurrent":        1,
		"description":       "test_description",
		"network_id":        "test_id",
		"private_address":   "192.168.1.1/24",
	})

	tests := []struct {
		name string
		args *schema.ResourceData
		want *devopsrunner.CreateRunnerInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &devopsrunner.CreateRunnerInput{
				RunnerName:       nifcloud.String("test_name"),
				InstanceType:     types.InstanceTypeOfCreateRunnerRequest("c-small"),
				AvailabilityZone: types.AvailabilityZoneOfCreateRunnerRequest("east-11"),
				Concurrent:       nifcloud.Int32(1),
				Description:      nifcloud.String("test_description"),
				NetworkConfig: &types.RequestNetworkConfig{
					NetworkId:      nifcloud.String("test_id"),
					PrivateAddress: nifcloud.String("192.168.1.1/24"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandCreateRunnerInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandUpdateInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"name":              "test_name_changed",
		"instance_type":     "c-small",
		"availability_zone": "east-11",
		"concurrent":        1,
		"description":       "test_description",
		"network_id":        "test_id",
		"private_address":   "192.168.1.1/24",
	})
	rd.SetId("test_name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *devopsrunner.UpdateRunnerInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &devopsrunner.UpdateRunnerInput{
				RunnerName:        nifcloud.String("test_name"),
				ChangedRunnerName: nifcloud.String("test_name_changed"),
				Concurrent:        nifcloud.Int32(int32(1)),
				Description:       nifcloud.String("test_description"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandUpdateRunnerInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandGetsInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *devopsrunner.GetRunnerInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &devopsrunner.GetRunnerInput{
				RunnerName: nifcloud.String("test_name"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandGetRunnerInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDeleteInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *devopsrunner.DeleteRunnerInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &devopsrunner.DeleteRunnerInput{
				RunnerName: nifcloud.String("test_name"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDeleteRunnerInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}
