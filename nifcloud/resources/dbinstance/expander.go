package dbinstance

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/rdb"
)

func expandCreateDBInstanceInput(d *schema.ResourceData) *rdb.CreateDBInstanceInput {
	input := &rdb.CreateDBInstanceInput{
		AccountingType:                       nifcloud.String(d.Get("accounting_type").(string)),
		AllocatedStorage:                     nifcloud.Int64(int64(d.Get("allocated_storage").(int))),
		AvailabilityZone:                     nifcloud.String(d.Get("availability_zone").(string)),
		BackupRetentionPeriod:                nifcloud.Int64(int64(d.Get("backup_retention_period").(int))),
		DBInstanceClass:                      nifcloud.String(d.Get("instance_class").(string)),
		DBInstanceIdentifier:                 nifcloud.String(d.Get("identifier").(string)),
		DBName:                               nifcloud.String(d.Get("db_name").(string)),
		DBParameterGroupName:                 nifcloud.String(d.Get("parameter_group_name").(string)),
		DBSecurityGroups:                     []string{d.Get("db_security_group_name").(string)},
		Engine:                               nifcloud.String(d.Get("engine").(string)),
		EngineVersion:                        nifcloud.String(d.Get("engine_version").(string)),
		MasterUserPassword:                   nifcloud.String(d.Get("password").(string)),
		MasterUsername:                       nifcloud.String(d.Get("username").(string)),
		MultiAZ:                              nifcloud.Bool(d.Get("multi_az").(bool)),
		NiftyMultiAZType:                     nifcloud.Int64(int64(d.Get("multi_az_type").(int))),
		Port:                                 nifcloud.Int64(int64(d.Get("port").(int))),
		PreferredBackupWindow:                nifcloud.String(d.Get("backup_window").(string)),
		PreferredMaintenanceWindow:           nifcloud.String(d.Get("maintenance_window").(string)),
		PubliclyAccessible:                   nifcloud.Bool(d.Get("publicly_accessible").(bool)),
		NiftyStorageType:                     nifcloud.Int64(int64(d.Get("storage_type").(int))),
		NiftyNetworkId:                       nifcloud.String(d.Get("network_id").(string)),
		NiftyVirtualPrivateAddress:           nifcloud.String(d.Get("virtual_private_address").(string)),
		NiftyMasterPrivateAddress:            nifcloud.String(d.Get("master_private_address").(string)),
		NiftySlavePrivateAddress:             nifcloud.String(d.Get("slave_private_address").(string)),
		NiftyReadReplicaPrivateAddress:       nifcloud.String(d.Get("read_replica_private_address").(string)),
		NiftyReadReplicaDBInstanceIdentifier: nifcloud.String(d.Get("read_replica_identifier").(string)),
		ReadReplicaAccountingType:            nifcloud.String(d.Get("accounting_type").(string)),
	}
	return input
}

func expandCreateDBInstanceReadReplicaInput(d *schema.ResourceData) *rdb.CreateDBInstanceReadReplicaInput {
	input := &rdb.CreateDBInstanceReadReplicaInput{
		AccountingType:                 nifcloud.String(d.Get("accounting_type").(string)),
		DBInstanceClass:                nifcloud.String(d.Get("instance_class").(string)),
		DBInstanceIdentifier:           nifcloud.String(d.Get("identifier").(string)),
		SourceDBInstanceIdentifier:     nifcloud.String(d.Get("replicate_source_db").(string)),
		NiftyStorageType:               nifcloud.Int64(int64(d.Get("storage_type").(int))),
		NiftyReadReplicaPrivateAddress: nifcloud.String(d.Get("read_replica_private_address").(string)),
	}
	return input
}

func expandRestoreDBInstanceFromDBSnapshotInput(d *schema.ResourceData) *rdb.RestoreDBInstanceFromDBSnapshotInput {
	input := &rdb.RestoreDBInstanceFromDBSnapshotInput{
		AccountingType:                       nifcloud.String(d.Get("accounting_type").(string)),
		AvailabilityZone:                     nifcloud.String(d.Get("availability_zone").(string)),
		DBInstanceClass:                      nifcloud.String(d.Get("instance_class").(string)),
		DBInstanceIdentifier:                 nifcloud.String(d.Get("identifier").(string)),
		NiftyDBParameterGroupName:            nifcloud.String(d.Get("parameter_group_name").(string)),
		NiftyDBSecurityGroups:                []string{d.Get("db_security_group_name").(string)},
		DBSnapshotIdentifier:                 nifcloud.String(d.Get("snapshot_identifier").(string)),
		MultiAZ:                              nifcloud.Bool(d.Get("multi_az").(bool)),
		NiftyMultiAZType:                     nifcloud.Int64(int64(d.Get("multi_az_type").(int))),
		Port:                                 nifcloud.Int64(int64(d.Get("port").(int))),
		PubliclyAccessible:                   nifcloud.Bool(d.Get("publicly_accessible").(bool)),
		NiftyStorageType:                     nifcloud.Int64(int64(d.Get("storage_type").(int))),
		NiftyNetworkId:                       nifcloud.String(d.Get("network_id").(string)),
		NiftyVirtualPrivateAddress:           nifcloud.String(d.Get("virtual_private_address").(string)),
		NiftyMasterPrivateAddress:            nifcloud.String(d.Get("master_private_address").(string)),
		NiftySlavePrivateAddress:             nifcloud.String(d.Get("slave_private_address").(string)),
		NiftyReadReplicaPrivateAddress:       nifcloud.String(d.Get("read_replica_private_address").(string)),
		NiftyReadReplicaDBInstanceIdentifier: nifcloud.String(d.Get("read_replica_identifier").(string)),
		ReadReplicaAccountingType:            nifcloud.String(d.Get("accounting_type").(string)),
	}
	return input
}

func expandRestoreDBInstanceToPointInTimeInput(d *schema.ResourceData) *rdb.RestoreDBInstanceToPointInTimeInput {
	tfMap := d.Get("restore_to_point_in_time").([]interface{})[0].(map[string]interface{})

	input := &rdb.RestoreDBInstanceToPointInTimeInput{
		AccountingType:                       nifcloud.String(d.Get("accounting_type").(string)),
		AvailabilityZone:                     nifcloud.String(d.Get("availability_zone").(string)),
		DBInstanceClass:                      nifcloud.String(d.Get("instance_class").(string)),
		NiftyDBParameterGroupName:            nifcloud.String(d.Get("parameter_group_name").(string)),
		NiftyDBSecurityGroups:                []string{d.Get("db_security_group_name").(string)},
		MultiAZ:                              nifcloud.Bool(d.Get("multi_az").(bool)),
		NiftyMultiAZType:                     nifcloud.Int64(int64(d.Get("multi_az_type").(int))),
		Port:                                 nifcloud.Int64(int64(d.Get("port").(int))),
		PubliclyAccessible:                   nifcloud.Bool(d.Get("publicly_accessible").(bool)),
		TargetDBInstanceIdentifier:           nifcloud.String(d.Get("identifier").(string)),
		NiftyStorageType:                     nifcloud.Int64(int64(d.Get("storage_type").(int))),
		NiftyNetworkId:                       nifcloud.String(d.Get("network_id").(string)),
		NiftyVirtualPrivateAddress:           nifcloud.String(d.Get("virtual_private_address").(string)),
		NiftyMasterPrivateAddress:            nifcloud.String(d.Get("master_private_address").(string)),
		NiftySlavePrivateAddress:             nifcloud.String(d.Get("slave_private_address").(string)),
		NiftyReadReplicaPrivateAddress:       nifcloud.String(d.Get("read_replica_private_address").(string)),
		NiftyReadReplicaDBInstanceIdentifier: nifcloud.String(d.Get("read_replica_identifier").(string)),
		ReadReplicaAccountingType:            nifcloud.String(d.Get("accounting_type").(string)),
	}

	if v, ok := tfMap["restore_time"].(string); ok && v != "" {
		parsedTime, err := time.Parse(time.RFC3339, v)
		if err == nil {
			input.RestoreTime = nifcloud.Time(parsedTime)
		}
	}

	if v, ok := tfMap["source_db_instance_identifier"].(string); ok && v != "" {
		input.SourceDBInstanceIdentifier = nifcloud.String(v)
	}

	if v, ok := tfMap["use_latest_restorable_time"].(bool); ok && v {
		input.UseLatestRestorableTime = nifcloud.Bool(v)
	}
	return input
}

func expandDescribeDBInstancesInput(d *schema.ResourceData) *rdb.DescribeDBInstancesInput {
	input := &rdb.DescribeDBInstancesInput{
		DBInstanceIdentifier: nifcloud.String(d.Id()),
	}
	return input
}

func expandDeleteDBInstanceInput(d *schema.ResourceData) *rdb.DeleteDBInstanceInput {
	input := &rdb.DeleteDBInstanceInput{
		DBInstanceIdentifier:      nifcloud.String(d.Id()),
		FinalDBSnapshotIdentifier: nifcloud.String(d.Get("final_snapshot_identifier").(string)),
		SkipFinalSnapshot:         nifcloud.Bool(d.Get("skip_final_snapshot").(bool)),
	}
	return input
}

func expandModifyDBInstanceInput(d *schema.ResourceData) *rdb.ModifyDBInstanceInput {
	input := &rdb.ModifyDBInstanceInput{
		DBInstanceIdentifier:                 nifcloud.String(d.Id()),
		AccountingType:                       nifcloud.String(d.Get("accounting_type").(string)),
		ApplyImmediately:                     nifcloud.Bool(d.Get("apply_immediately").(bool)),
		AllocatedStorage:                     nifcloud.Int64(int64(d.Get("allocated_storage").(int))),
		BackupRetentionPeriod:                nifcloud.Int64(int64(d.Get("backup_retention_period").(int))),
		BinlogRetentionPeriod:                nifcloud.Int64(int64(d.Get("binlog_retention_period").(int))),
		CustomBinlogRetentionPeriod:          nifcloud.Bool(d.Get("custom_binlog_retention_period").(bool)),
		DBInstanceClass:                      nifcloud.String(d.Get("instance_class").(string)),
		DBParameterGroupName:                 nifcloud.String(d.Get("parameter_group_name").(string)),
		DBSecurityGroups:                     []string{d.Get("db_security_group_name").(string)},
		MasterUserPassword:                   nifcloud.String(d.Get("password").(string)),
		MultiAZ:                              nifcloud.Bool(d.Get("multi_az").(bool)),
		NiftyMultiAZType:                     nifcloud.Int64(int64(d.Get("multi_az_type").(int))),
		PreferredBackupWindow:                nifcloud.String(d.Get("backup_window").(string)),
		PreferredMaintenanceWindow:           nifcloud.String(d.Get("maintenance_window").(string)),
		NiftySlavePrivateAddress:             nifcloud.String(d.Get("slave_private_address").(string)),
		NiftyReadReplicaDBInstanceIdentifier: nifcloud.String(d.Get("read_replica_identifier").(string)),
		NiftyReadReplicaPrivateAddress:       nifcloud.String(d.Get("read_replica_private_address").(string)),
		ReadReplicaAccountingType:            nifcloud.String(d.Get("accounting_type").(string)),
	}

	if d.HasChange("identifier") && !d.IsNewResource() {
		input.NewDBInstanceIdentifier = nifcloud.String(d.Get("identifier").(string))
	}

	return input
}

func expandModifyDBInstanceNetworkInput(d *schema.ResourceData) *rdb.ModifyDBInstanceNetworkInput {
	input := &rdb.ModifyDBInstanceNetworkInput{
		DBInstanceIdentifier:       nifcloud.String(d.Id()),
		NiftySlavePrivateAddress:   nifcloud.String(d.Get("slave_private_address").(string)),
		NiftyNetworkId:             nifcloud.String(d.Get("network_id").(string)),
		NiftyVirtualPrivateAddress: nifcloud.String(d.Get("virtual_private_address").(string)),
		NiftyMasterPrivateAddress:  nifcloud.String(d.Get("master_private_address").(string)),
	}
	return input
}
