package cluster

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/hatoba"
	"github.com/stretchr/testify/assert"
)

func TestFlatten(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"name":               "test_name",
		"description":        "test_description",
		"kubernetes_version": "test_kubernetes_version",
		"locations":          []interface{}{"test_location"},
		"addons_config": []interface{}{map[string]interface{}{
			"http_load_balancing": []interface{}{map[string]interface{}{
				"disabled": true,
			}},
		}},
		"network_config": []interface{}{map[string]interface{}{
			"network_id": "test_network_id",
		}},
		"node_pools": []interface{}{map[string]interface{}{
			"name":          "test_node_pool_name",
			"instance_type": "test_instance_type",
			"node_count":    3,
			"nodes": []interface{}{
				map[string]interface{}{
					"name":               "test_node01",
					"availability_zone":  "test_availability_zone",
					"public_ip_address":  "203.0.113.1",
					"private_ip_address": "192.168.0.1",
				},
				map[string]interface{}{
					"name":               "test_node02",
					"availability_zone":  "test_availability_zone",
					"public_ip_address":  "203.0.113.2",
					"private_ip_address": "192.168.0.2",
				},
				map[string]interface{}{
					"name":               "test_node03",
					"availability_zone":  "test_availability_zone",
					"public_ip_address":  "203.0.113.3",
					"private_ip_address": "192.168.0.3",
				},
			}},
		},
	})
	rd.SetId("test_name")

	wantNotFoundRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})

	type args struct {
		res *hatoba.GetClusterResponse
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
				res: &hatoba.GetClusterResponse{
					GetClusterOutput: &hatoba.GetClusterOutput{
						Cluster: &hatoba.Cluster{
							Name:              nifcloud.String("test_name"),
							Description:       nifcloud.String("test_description"),
							KubernetesVersion: nifcloud.String("test_kubernetes_version"),
							Locations:         []string{"test_location"},
							AddonsConfig: &hatoba.AddonsConfig{
								HttpLoadBalancing: &hatoba.HttpLoadBalancing{
									Disabled: nifcloud.Bool(true),
								},
							},
							NetworkConfig: &hatoba.NetworkConfig{
								NetworkId: nifcloud.String("test_network_id"),
							},
							NodePools: []hatoba.NodePool{
								{
									Name:         nifcloud.String("test_node_pool_name"),
									InstanceType: nifcloud.String("test_instance_type"),
									NodeCount:    nifcloud.Int64(3),
									Nodes: []hatoba.Node{
										{
											Name:             nifcloud.String("test_node01"),
											AvailabilityZone: nifcloud.String("test_availability_zone"),
											PublicIpAddress:  nifcloud.String("203.0.113.1"),
											PrivateIpAddress: nifcloud.String("192.168.0.1"),
										},
										{
											Name:             nifcloud.String("test_node02"),
											AvailabilityZone: nifcloud.String("test_availability_zone"),
											PublicIpAddress:  nifcloud.String("203.0.113.2"),
											PrivateIpAddress: nifcloud.String("192.168.0.2"),
										},
										{
											Name:             nifcloud.String("test_node03"),
											AvailabilityZone: nifcloud.String("test_availability_zone"),
											PublicIpAddress:  nifcloud.String("203.0.113.3"),
											PrivateIpAddress: nifcloud.String("192.168.0.3"),
										},
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
				res: &hatoba.GetClusterResponse{
					GetClusterOutput: &hatoba.GetClusterOutput{
						Cluster: &hatoba.Cluster{},
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
