package dbinstance

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/rdb"
	"github.com/stretchr/testify/assert"
)

func TestExpandCreateInstanceInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"accounting_type":                1,
		"availability_zone":              "test_availability_zone",
		"instance_class":                 "test_instance_class",
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
		"multi_az_type":                  1,
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
		"read_replica_identifier":        "test_read_replica_identifier",
		"apply_immediately":              true,
	})
	rd.SetId("test_identifier")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *rdb.CreateDBInstanceInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &rdb.CreateDBInstanceInput{
				AccountingType:                       nifcloud.String("1"),
				AllocatedStorage:                     nifcloud.Int64(1),
				AvailabilityZone:                     nifcloud.String("test_availability_zone"),
				BackupRetentionPeriod:                nifcloud.Int64(1),
				DBInstanceClass:                      nifcloud.String("test_instance_class"),
				DBInstanceIdentifier:                 nifcloud.String("test_identifier"),
				DBName:                               nifcloud.String("test_db_name"),
				DBParameterGroupName:                 nifcloud.String("test_parameter_group_name"),
				DBSecurityGroups:                     []string{"test_db_security_group_name"},
				Engine:                               nifcloud.String("test_engine"),
				EngineVersion:                        nifcloud.String("test_engine_version"),
				MasterUserPassword:                   nifcloud.String("test_password"),
				MasterUsername:                       nifcloud.String("test_username"),
				MultiAZ:                              nifcloud.Bool(true),
				NiftyMultiAZType:                     nifcloud.Int64(1),
				Port:                                 nifcloud.Int64(1),
				PreferredBackupWindow:                nifcloud.String("test_backup_window"),
				PreferredMaintenanceWindow:           nifcloud.String("test_maintenance_window"),
				PubliclyAccessible:                   nifcloud.Bool(true),
				NiftyStorageType:                     nifcloud.Int64(1),
				NiftyNetworkId:                       nifcloud.String("test_network_id"),
				NiftyVirtualPrivateAddress:           nifcloud.String("test_virtual_private_address"),
				NiftyMasterPrivateAddress:            nifcloud.String("test_master_private_address"),
				NiftySlavePrivateAddress:             nifcloud.String("test_slave_private_address"),
				NiftyReadReplicaPrivateAddress:       nifcloud.String("test_read_replica_private_address"),
				NiftyReadReplicaDBInstanceIdentifier: nifcloud.String("test_read_replica_identifier"),
				ReadReplicaAccountingType:            nifcloud.String("1"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandCreateDBInstanceInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandCreateDBInstanceReadReplicaInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"accounting_type":                1,
		"availability_zone":              "test_availability_zone",
		"instance_class":                 "test_instance_class",
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
		"multi_az_type":                  1,
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
		"read_replica_identifier":        "test_read_replica_identifier",
		"apply_immediately":              true,
	})
	rd.SetId("test_identifier")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *rdb.CreateDBInstanceReadReplicaInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &rdb.CreateDBInstanceReadReplicaInput{
				AccountingType:                 nifcloud.String("1"),
				DBInstanceClass:                nifcloud.String("test_instance_class"),
				DBInstanceIdentifier:           nifcloud.String("test_identifier"),
				SourceDBInstanceIdentifier:     nifcloud.String("test_replicate_source_db"),
				NiftyStorageType:               nifcloud.Int64(1),
				NiftyReadReplicaPrivateAddress: nifcloud.String("test_read_replica_private_address"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandCreateDBInstanceReadReplicaInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandRestoreDBInstanceFromDBSnapshotInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"accounting_type":                1,
		"availability_zone":              "test_availability_zone",
		"instance_class":                 "test_instance_class",
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
		"multi_az_type":                  1,
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
		"read_replica_identifier":        "test_read_replica_identifier",
		"apply_immediately":              true,
	})
	rd.SetId("test_identifier")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *rdb.RestoreDBInstanceFromDBSnapshotInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &rdb.RestoreDBInstanceFromDBSnapshotInput{
				AccountingType:                       nifcloud.String("1"),
				AvailabilityZone:                     nifcloud.String("test_availability_zone"),
				DBInstanceClass:                      nifcloud.String("test_instance_class"),
				DBInstanceIdentifier:                 nifcloud.String("test_identifier"),
				NiftyDBParameterGroupName:            nifcloud.String("test_parameter_group_name"),
				NiftyDBSecurityGroups:                []string{"test_db_security_group_name"},
				DBSnapshotIdentifier:                 nifcloud.String("test_snapshot_identifier"),
				MultiAZ:                              nifcloud.Bool(true),
				NiftyMultiAZType:                     nifcloud.Int64(1),
				Port:                                 nifcloud.Int64(1),
				PubliclyAccessible:                   nifcloud.Bool(true),
				NiftyStorageType:                     nifcloud.Int64(1),
				NiftyNetworkId:                       nifcloud.String("test_network_id"),
				NiftyVirtualPrivateAddress:           nifcloud.String("test_virtual_private_address"),
				NiftyMasterPrivateAddress:            nifcloud.String("test_master_private_address"),
				NiftySlavePrivateAddress:             nifcloud.String("test_slave_private_address"),
				NiftyReadReplicaPrivateAddress:       nifcloud.String("test_read_replica_private_address"),
				NiftyReadReplicaDBInstanceIdentifier: nifcloud.String("test_read_replica_identifier"),
				ReadReplicaAccountingType:            nifcloud.String("1"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandRestoreDBInstanceFromDBSnapshotInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandRestoreDBInstanceToPointInTimeInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"accounting_type":                1,
		"availability_zone":              "test_availability_zone",
		"instance_class":                 "test_instance_class",
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
		"multi_az_type":                  1,
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
		"read_replica_identifier":        "test_read_replica_identifier",
		"apply_immediately":              true,
		"restore_to_point_in_time": []interface{}{map[string]interface{}{
			"restore_time":                  "2001-02-03T04:05:06Z",
			"source_db_instance_identifier": "test_source_db_instance_identifier",
			"use_latest_restorable_time":    true,
		}},
	})
	rd.SetId("test_identifier")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *rdb.RestoreDBInstanceToPointInTimeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &rdb.RestoreDBInstanceToPointInTimeInput{
				AccountingType:                       nifcloud.String("1"),
				AvailabilityZone:                     nifcloud.String("test_availability_zone"),
				DBInstanceClass:                      nifcloud.String("test_instance_class"),
				NiftyDBParameterGroupName:            nifcloud.String("test_parameter_group_name"),
				NiftyDBSecurityGroups:                []string{"test_db_security_group_name"},
				MultiAZ:                              nifcloud.Bool(true),
				NiftyMultiAZType:                     nifcloud.Int64(1),
				Port:                                 nifcloud.Int64(1),
				PubliclyAccessible:                   nifcloud.Bool(true),
				TargetDBInstanceIdentifier:           nifcloud.String("test_identifier"),
				NiftyStorageType:                     nifcloud.Int64(1),
				NiftyNetworkId:                       nifcloud.String("test_network_id"),
				NiftyVirtualPrivateAddress:           nifcloud.String("test_virtual_private_address"),
				NiftyMasterPrivateAddress:            nifcloud.String("test_master_private_address"),
				NiftySlavePrivateAddress:             nifcloud.String("test_slave_private_address"),
				NiftyReadReplicaPrivateAddress:       nifcloud.String("test_read_replica_private_address"),
				NiftyReadReplicaDBInstanceIdentifier: nifcloud.String("test_read_replica_identifier"),
				ReadReplicaAccountingType:            nifcloud.String("1"),
				RestoreTime:                          nifcloud.Time(time.Date(2001, 2, 3, 4, 5, 6, 0, time.UTC)),
				SourceDBInstanceIdentifier:           nifcloud.String("test_source_db_instance_identifier"),
				UseLatestRestorableTime:              nifcloud.Bool(true),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandRestoreDBInstanceToPointInTimeInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDescribeInstanceInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"accounting_type":                1,
		"availability_zone":              "test_availability_zone",
		"instance_class":                 "test_instance_class",
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
		"multi_az_type":                  1,
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
		"read_replica_identifier":        "test_read_replica_identifier",
		"apply_immediately":              true,
	})
	rd.SetId("test_identifier")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *rdb.DescribeDBInstancesInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &rdb.DescribeDBInstancesInput{
				DBInstanceIdentifier: nifcloud.String("test_identifier"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDescribeDBInstancesInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDeleteInstanceInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"accounting_type":                1,
		"availability_zone":              "test_availability_zone",
		"instance_class":                 "test_instance_class",
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
		"multi_az_type":                  1,
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
		"read_replica_identifier":        "test_read_replica_identifier",
		"apply_immediately":              true,
	})
	rd.SetId("test_identifier")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *rdb.DeleteDBInstanceInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &rdb.DeleteDBInstanceInput{
				DBInstanceIdentifier:      nifcloud.String("test_identifier"),
				FinalDBSnapshotIdentifier: nifcloud.String("test_final_snapshot_identifier"),
				SkipFinalSnapshot:         nifcloud.Bool(true),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDeleteDBInstanceInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandModifyInstanceInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"accounting_type":                1,
		"availability_zone":              "test_availability_zone",
		"instance_class":                 "test_instance_class",
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
		"multi_az_type":                  1,
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
		"read_replica_identifier":        "test_read_replica_identifier",
		"apply_immediately":              true,
	})
	rd.SetId("test_identifier")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *rdb.ModifyDBInstanceInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &rdb.ModifyDBInstanceInput{
				AccountingType:                       nifcloud.String("1"),
				AllocatedStorage:                     nifcloud.Int64(1),
				BackupRetentionPeriod:                nifcloud.Int64(1),
				DBInstanceClass:                      nifcloud.String("test_instance_class"),
				DBInstanceIdentifier:                 nifcloud.String("test_identifier"),
				BinlogRetentionPeriod:                nifcloud.Int64(1),
				CustomBinlogRetentionPeriod:          nifcloud.Bool(true),
				DBParameterGroupName:                 nifcloud.String("test_parameter_group_name"),
				DBSecurityGroups:                     []string{"test_db_security_group_name"},
				MasterUserPassword:                   nifcloud.String("test_password"),
				MultiAZ:                              nifcloud.Bool(true),
				NiftyMultiAZType:                     nifcloud.Int64(1),
				PreferredBackupWindow:                nifcloud.String("test_backup_window"),
				PreferredMaintenanceWindow:           nifcloud.String("test_maintenance_window"),
				NiftySlavePrivateAddress:             nifcloud.String("test_slave_private_address"),
				NiftyReadReplicaPrivateAddress:       nifcloud.String("test_read_replica_private_address"),
				NiftyReadReplicaDBInstanceIdentifier: nifcloud.String("test_read_replica_identifier"),
				ReadReplicaAccountingType:            nifcloud.String("1"),
				NewDBInstanceIdentifier:              nifcloud.String("test_identifier"),
				ApplyImmediately:                     nifcloud.Bool(true),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandModifyDBInstanceInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandModifyDBInstanceNetworkInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"accounting_type":                1,
		"availability_zone":              "test_availability_zone",
		"instance_class":                 "test_instance_class",
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
		"multi_az_type":                  1,
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
		"read_replica_identifier":        "test_read_replica_identifier",
		"apply_immediately":              true,
	})
	rd.SetId("test_identifier")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *rdb.ModifyDBInstanceNetworkInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &rdb.ModifyDBInstanceNetworkInput{
				DBInstanceIdentifier:       nifcloud.String("test_identifier"),
				NiftyNetworkId:             nifcloud.String("test_network_id"),
				NiftySlavePrivateAddress:   nifcloud.String("test_slave_private_address"),
				NiftyMasterPrivateAddress:  nifcloud.String("test_master_private_address"),
				NiftyVirtualPrivateAddress: nifcloud.String("test_virtual_private_address"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandModifyDBInstanceNetworkInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}
