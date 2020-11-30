package volume

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/stretchr/testify/assert"
)

func TestFlatten(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"size":               100,
		"volume_id":          "test_volume_id",
		"disk_type":          "High-Speed Storage A",
		"accounting_type":    "1",
		"description":        "test_description",
		"instance_id":        "test_instance_id",
		"instance_unique_id": "test_instance_unique_id",
	})
	rd.SetId("test_volume_id")

	wantNotFoundRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})

	type args struct {
		res *computing.DescribeVolumesResponse
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
				res: &computing.DescribeVolumesResponse{
					DescribeVolumesOutput: &computing.DescribeVolumesOutput{
						VolumeSet: []computing.VolumeSet{
							{
								Size:           nifcloud.String("100"),
								VolumeId:       nifcloud.String("test_volume_id"),
								DiskType:       nifcloud.String("High-Speed Storage A"),
								AccountingType: nifcloud.String("1"),
								Description:    nifcloud.String("test_description"),
								AttachmentSet: []computing.AttachmentSet{
									{
										InstanceId:       nifcloud.String("test_instance_id"),
										InstanceUniqueId: nifcloud.String("test_instance_unique_id"),
									},
								},
							},
						},
					},
				},
			},
			want: rd,
		},
		{
			name: "flattens the response even when the resource has been removed externally",
			args: args{
				d: wantNotFoundRd,
				res: &computing.DescribeVolumesResponse{
					DescribeVolumesOutput: &computing.DescribeVolumesOutput{
						VolumeSet: []computing.VolumeSet{},
					},
				},
			},
			want: wantNotFoundRd,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := flatten(tt.args.d, tt.args.res)
			assert.NoError(t, err)

			if tt.args.res == nil {
				return
			}

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
