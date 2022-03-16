package separateinstancerule

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing/types"
	"github.com/stretchr/testify/assert"
)

func TestExpandNiftyCreateSeparateInstanceRuleInputForInstanceId(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"name":              "test_separate_name",
		"availability_zone": "test_availability_zone",
		"description":       "test_description",
		"instance_id":       []interface{}{"test_instance_id1", "test_instance_id2"},
	})
	rd.SetId("test_separate_name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyCreateSeparateInstanceRuleInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyCreateSeparateInstanceRuleInput{
				SeparateInstanceRuleName: nifcloud.String("test_separate_name"),
				Placement: &types.RequestPlacementOfNiftyCreateSeparateInstanceRule{
					AvailabilityZone: nifcloud.String(("test_availability_zone")),
				},
				SeparateInstanceRuleDescription: nifcloud.String("test_description"),
				InstanceId:                      []string{"test_instance_id1", "test_instance_id2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyCreateSeparateInstanceRuleInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyCreateSeparateInstanceRuleInputForInstanceUniqueId(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"name":               "test_separate_name",
		"availability_zone":  "test_availability_zone",
		"description":        "test_description",
		"instance_unique_id": []interface{}{"test_instance_unique_id1", "test_instance_unique_id2"},
	})
	rd.SetId("test_separate_name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyCreateSeparateInstanceRuleInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyCreateSeparateInstanceRuleInput{
				SeparateInstanceRuleName: nifcloud.String("test_separate_name"),
				Placement: &types.RequestPlacementOfNiftyCreateSeparateInstanceRule{
					AvailabilityZone: nifcloud.String(("test_availability_zone")),
				},
				SeparateInstanceRuleDescription: nifcloud.String("test_description"),
				InstanceUniqueId:                []string{"test_instance_unique_id1", "test_instance_unique_id2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyCreateSeparateInstanceRuleInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyUpdateSeparateInstanceRuleInputForName(t *testing.T) {
	r := New()
	rd := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"name": "test_separate_name_new",
	})
	rd.SetId("test_separate_name")
	dn := r.Data(rd.State())

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyUpdateSeparateInstanceRuleInput
	}{
		{
			name: "expands the resource data",
			args: dn,
			want: &computing.NiftyUpdateSeparateInstanceRuleInput{
				SeparateInstanceRuleName:       nifcloud.String("test_separate_name"),
				SeparateInstanceRuleNameUpdate: nifcloud.String("test_separate_name_new"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyUpdateSeparateInstanceRuleInputForName(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyUpdateSeparateInstanceRuleInputForDescription(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"name":        "test_separate_name",
		"description": "test_description",
	})
	rd.SetId("test_separate_name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyUpdateSeparateInstanceRuleInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyUpdateSeparateInstanceRuleInput{
				SeparateInstanceRuleName:              nifcloud.String("test_separate_name"),
				SeparateInstanceRuleDescriptionUpdate: nifcloud.String("test_description"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyUpdateSeparateInstanceRuleInputForDescription(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyDescribeSeparateInstanceRulesInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"name": "test_separate_name",
	})
	rd.SetId("test_separate_name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyDescribeSeparateInstanceRulesInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyDescribeSeparateInstanceRulesInput{
				SeparateInstanceRuleName: []string{"test_separate_name"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyDescribeSeparateInstanceRulesInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyDeleteSeparateInstanceRuleInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"name": "test_separate_name",
	})
	rd.SetId("test_separate_name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyDeleteSeparateInstanceRuleInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyDeleteSeparateInstanceRuleInput{
				SeparateInstanceRuleName: nifcloud.String("test_separate_name"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyDeleteSeparateInstanceRuleInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyRegisterInstancesWithSeparateInstanceRuleInstanceIDInput(t *testing.T) {
	instanceId := []string{"test_instance_id1", "test_instance_id2"}
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_separate_name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyRegisterInstancesWithSeparateInstanceRuleInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyRegisterInstancesWithSeparateInstanceRuleInput{
				SeparateInstanceRuleName: nifcloud.String("test_separate_name"),
				InstanceId:               []string{"test_instance_id1", "test_instance_id2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyRegisterInstancesWithSeparateInstanceRuleInstanceIDInput(tt.args, instanceId)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyDeregisterInstancesFromSeparateInstanceRuleInstanceIDInput(t *testing.T) {
	instanceId := []string{"test_instance_id1", "test_instance_id2"}
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_separate_name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyDeregisterInstancesFromSeparateInstanceRuleInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyDeregisterInstancesFromSeparateInstanceRuleInput{
				SeparateInstanceRuleName: nifcloud.String("test_separate_name"),
				InstanceId:               []string{"test_instance_id1", "test_instance_id2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyDeregisterInstancesFromSeparateInstanceRuleInstanceIDInput(tt.args, instanceId)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyRegisterInstancesWithSeparateInstanceRuleInstanceUniqueIDInput(t *testing.T) {
	instanceUniqueId := []string{"test_instance_unique_id1", "test_instance_unique_id2"}
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_separate_name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyRegisterInstancesWithSeparateInstanceRuleInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyRegisterInstancesWithSeparateInstanceRuleInput{
				SeparateInstanceRuleName: nifcloud.String("test_separate_name"),
				InstanceUniqueId:         []string{"test_instance_unique_id1", "test_instance_unique_id2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyRegisterInstancesWithSeparateInstanceRuleInstanceUniqueIDInput(tt.args, instanceUniqueId)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandNiftyDeregisterInstancesFromSeparateInstanceRuleInstanceUniqueIDInput(t *testing.T) {
	instanceUniqueId := []string{"test_instance_unique_id1", "test_instance_unique_id2"}
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_separate_name")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.NiftyDeregisterInstancesFromSeparateInstanceRuleInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.NiftyDeregisterInstancesFromSeparateInstanceRuleInput{
				SeparateInstanceRuleName: nifcloud.String("test_separate_name"),
				InstanceUniqueId:         []string{"test_instance_unique_id1", "test_instance_unique_id2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandNiftyDeregisterInstancesFromSeparateInstanceRuleInstanceUniqueIDInput(tt.args, instanceUniqueId)
			assert.Equal(t, tt.want, got)
		})
	}
}
