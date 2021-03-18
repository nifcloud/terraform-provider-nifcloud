package volume

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/stretchr/testify/assert"
)

func TestExpandCreateVolumeInputForInstanceId(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"size":            100,
		"volume_id":       "test_volume_id",
		"disk_type":       "High-Speed Storage A",
		"accounting_type": "1",
		"description":     "test_description",
		"instance_id":     "test_instance_id",
	})

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.CreateVolumeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.CreateVolumeInput{
				Size:           nifcloud.Int64(100),
				VolumeId:       nifcloud.String("test_volume_id"),
				DiskType:       computing.DiskTypeOfCreateVolumeRequest3,
				AccountingType: computing.AccountingTypeOfCreateVolumeRequest1,
				Description:    nifcloud.String("test_description"),
				InstanceId:     nifcloud.String("test_instance_id"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandCreateVolumeInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandCreateVolumeInputForInstanceUniqueId(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"size":               100,
		"volume_id":          "test_volume_id",
		"disk_type":          "High-Speed Storage A",
		"accounting_type":    "1",
		"description":        "test_description",
		"instance_unique_id": "test_instance_unique_id",
	})

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.CreateVolumeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.CreateVolumeInput{
				Size:             nifcloud.Int64(100),
				VolumeId:         nifcloud.String("test_volume_id"),
				DiskType:         computing.DiskTypeOfCreateVolumeRequest3,
				AccountingType:   computing.AccountingTypeOfCreateVolumeRequest1,
				Description:      nifcloud.String("test_description"),
				InstanceUniqueId: nifcloud.String("test_instance_unique_id"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandCreateVolumeInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandModifyVolumeAttributeInputForAccountingType(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"volume_id":       "test_volume_id",
		"accounting_type": "1",
	})
	rd.SetId("test_volume_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.ModifyVolumeAttributeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.ModifyVolumeAttributeInput{
				VolumeId:  nifcloud.String("test_volume_id"),
				Attribute: computing.AttributeOfModifyVolumeAttributeRequestAccountingType,
				Value:     nifcloud.String("1"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandModifyVolumeAttributeInputForAccountingType(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandModifyVolumeAttributeInputForVolumeName(t *testing.T) {
	r := New()
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"volume_id": "test_volume_id",
	})
	rd.SetId("test_volume_id")
	dn := r.Data(rd.State())

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.ModifyVolumeAttributeInput
	}{
		{
			name: "expands the resource data",
			args: dn,
			want: &computing.ModifyVolumeAttributeInput{
				VolumeId:  nifcloud.String("test_volume_id"),
				Attribute: computing.AttributeOfModifyVolumeAttributeRequestVolumeName,
				Value:     nifcloud.String("test_volume_id"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandModifyVolumeAttributeInputForVolumeName(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandModifyVolumeAttributeInputForDescription(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"volume_id":   "test_volume_id",
		"description": "test_description",
	})
	rd.SetId("test_volume_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.ModifyVolumeAttributeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.ModifyVolumeAttributeInput{
				VolumeId:  nifcloud.String("test_volume_id"),
				Attribute: computing.AttributeOfModifyVolumeAttributeRequestDescription,
				Value:     nifcloud.String("test_description"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandModifyVolumeAttributeInputForDescription(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandExtendVolumeSizeInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"volume_id": "test_volume_id",
		"reboot":    "true",
	})
	rd.SetId("test_volume_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.ExtendVolumeSizeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.ExtendVolumeSizeInput{
				VolumeId:    nifcloud.String("test_volume_id"),
				NiftyReboot: computing.NiftyRebootOfExtendVolumeSizeRequest("true"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandExtendVolumeSizeInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDescribeVolumesInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"volume_id": "test_volume_id",
	})
	rd.SetId("test_volume_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.DescribeVolumesInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.DescribeVolumesInput{
				VolumeId: []string{"test_volume_id"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDescribeVolumesInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDetachVolumeInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"volume_id":   "test_volume_id",
		"instance_id": "test_instance_id",
	})
	rd.SetId("test_volume_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.DetachVolumeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.DetachVolumeInput{
				VolumeId:   nifcloud.String("test_volume_id"),
				InstanceId: nifcloud.String("test_instance_id"),
				Agreement:  nifcloud.Bool(true),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDetachVolumeInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDeleteVolumeInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"volume_id": "test_volume_id",
	})
	rd.SetId("test_volume_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.DeleteVolumeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.DeleteVolumeInput{
				VolumeId: nifcloud.String("test_volume_id"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDeleteVolumeInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}
