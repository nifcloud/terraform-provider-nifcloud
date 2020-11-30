package instance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/stretchr/testify/assert"
)

func TestFlatten(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"accounting_type":         "test_accounting_type",
		"admin":                   "test_admin",
		"availability_zone":       "test_availability_zone",
		"description":             "test_description",
		"disable_api_termination": false,
		"image_id":                "test_image_id",
		"instance_id":             "test_instance_id",
		"instance_type":           "test_instance_type",
		"key_name":                "test_key_name",
		"license_name":            "test_license_name",
		"license_num":             200,
		"network_interface": []interface{}{map[string]interface{}{
			"network_id":   "test_network_id",
			"network_name": "test_network_name",
			"ip_address":   "test_ip_address",
		}},
		"password":       "test_password",
		"security_group": "test_security_group",
		"user_data":      "test_user_data",
		"public_ip":      "test_public_ip",
		"private_ip":     "test_private_ip",
		"unique_id":      "test_unique_id",
		"instance_state": "test_instance_state",
	})
	rd.SetId("test_instance_id")

	wantNotFoundRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})

	type args struct {
		res *computing.DescribeInstancesResponse
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
				res: &computing.DescribeInstancesResponse{
					DescribeInstancesOutput: &computing.DescribeInstancesOutput{
						ReservationSet: []computing.ReservationSet{
							{
								InstancesSet: []computing.InstancesSet{
									{
										Description: nifcloud.String("test_description"),
										Placement: &computing.Placement{
											AvailabilityZone: nifcloud.String("test_availability_zone"),
										},
										NextMonthAccountingType: nifcloud.String("test_accountingT_type"),
										ImageId:                 nifcloud.String("test_image_id"),
										InstanceId:              nifcloud.String("test_instance_id"),
										InstanceType:            nifcloud.String("test_instance_type"),
										KeyName:                 nifcloud.String("test_key_name"),
										InstanceState: &computing.InstanceState{
											Name: nifcloud.String("test_instance_state"),
										},
										PrivateIpAddress: nifcloud.String("test_private_ip"),
										IpAddress:        nifcloud.String("test_public_ip"),
										InstanceUniqueId: nifcloud.String("test_unique_id"),
										NetworkInterfaceSet: []computing.NetworkInterfaceSetOfDescribeInstances{
											{
												NiftyNetworkId:   nifcloud.String("test_network_id"),
												NiftyNetworkName: nifcloud.String("test_network_name"),
												PrivateIpAddress: nifcloud.String("test_ip_address"),
											},
										},
									},
								},
								GroupSet: []computing.GroupSet{
									{
										GroupId: nifcloud.String("test_group_id"),
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
				res: &computing.DescribeInstancesResponse{
					DescribeInstancesOutput: &computing.DescribeInstancesOutput{
						ReservationSet: []computing.ReservationSet{},
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

func TestFlattenDisableAPITermination(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"disable_api_termination": false,
	})
	rd.SetId("test_instance_id")

	type args struct {
		res *computing.DescribeInstanceAttributeResponse
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
				res: &computing.DescribeInstanceAttributeResponse{
					DescribeInstanceAttributeOutput: &computing.DescribeInstanceAttributeOutput{
						DisableApiTermination: &computing.DisableApiTermination{
							Value: nifcloud.String("false"),
						},
					},
				},
			},
			want: rd,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := flattenDisableAPITermination(tt.args.d, tt.args.res)
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
