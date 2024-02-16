package dbinstance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/rdb"
	"github.com/nifcloud/nifcloud-sdk-go/service/rdb/types"
	"github.com/stretchr/testify/assert"
)

func TestFlatten(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"accounting_type":                1,
		"availability_zone":              "test_availability_zone",
		"instance_class":                 "instance_class",
		"db_name":                        "test_db_name",
		"username":                       "test_username",
		"password":                       "test_password",
		"engine":                         "test_engine",
		"engine_version":                 "test_engine_version",
		"ca_cert_identifier":             "test_ca_cert_identifier",
		"allocated_storage":              1,
		"storage_type":                   1,
		"identifier":                     "test_identifier",
		"backup_retention_period":        1,
		"binlog_retention_period":        1,
		"custom_binlog_retention_period": true,
		"backup_window":                  "test_backup_window",
		"maintenance_window":             "test_maintenance_window",
		"multi_az":                       true,
		"port":                           1,
		"publicly_accessible":            true,
		"db_security_group_name":         "test_db_security_group_name",
		"final_snapshot_identifier":      "test_final_snapshot_identifier",
		"skip_final_snapshot":            true,
		"parameter_group_name":           "test_parameter_group_name",
		"address":                        "test_address",
		"replicate_source_db":            "test_replicate_source_db",
		"snapshot_identifier":            "test_snapshot_identifier",
		"network_id":                     "test_network_id",
		"virtual_private_address":        "test_virtual_private_address",
		"master_private_address":         "test_master_private_address",
		"slave_private_address":          "test_slave_private_address",
		"read_replica_private_address":   "test_read_replica_private_address",
		"apply_immediately":              true,
	})
	rd.SetId("test_identifier")

	wantNotFoundRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})

	type args struct {
		res *rdb.DescribeDBInstancesOutput
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
				res: &rdb.DescribeDBInstancesOutput{
					DBInstances: []types.DBInstances{
						{
							AccountingType:          nifcloud.String("1"),
							AllocatedStorage:        nifcloud.Int32(1),
							AvailabilityZone:        nifcloud.String("test_availability_zone"),
							BackupRetentionPeriod:   nifcloud.Int32(1),
							BinlogRetentionPeriod:   nifcloud.Int32(1),
							CACertificateIdentifier: nifcloud.String("test_ca_cert_identifier"),
							DBInstanceClass:         nifcloud.String("test_instance_class"),
							DBInstanceIdentifier:    nifcloud.String("test_identifier"),
							DBName:                  nifcloud.String("test_db_name"),
							DBParameterGroups: []types.DBParameterGroups{{
								DBParameterGroupName: nifcloud.String("test_parameter_group_name")},
							},
							DBSecurityGroups: []types.DBSecurityGroups{{
								DBSecurityGroupName: nifcloud.String("test_db_security_group_name")},
							},
							Endpoint:                              &types.Endpoint{Address: nifcloud.String("address")},
							Engine:                                nifcloud.String("test_engine"),
							EngineVersion:                         nifcloud.String("test_engine_version"),
							MasterUsername:                        nifcloud.String("test_username"),
							MultiAZ:                               nifcloud.Bool(true),
							NextMonthAccountingType:               nifcloud.String("1"),
							NiftyMasterPrivateAddress:             nifcloud.String("test_master_private_address"),
							NiftyMultiAZType:                      nifcloud.String("0"),
							NiftyNetworkId:                        nifcloud.String("test_network_id"),
							NiftySlavePrivateAddress:              nifcloud.String("test_slave_private_address"),
							NiftyStorageType:                      nifcloud.Int32(1),
							PreferredBackupWindow:                 nifcloud.String("test_backup_window"),
							PreferredMaintenanceWindow:            nifcloud.String("test_maintenance_window"),
							PubliclyAccessible:                    nifcloud.Bool(true),
							ReadReplicaSourceDBInstanceIdentifier: nifcloud.String("test_replicate_source_db"),
						},
					},
				},
			},
			want: rd,
		},
		{
			name: "flattens the response even when the resource has been removed externally",
			args: args{
				d:   wantNotFoundRd,
				res: &rdb.DescribeDBInstancesOutput{},
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
