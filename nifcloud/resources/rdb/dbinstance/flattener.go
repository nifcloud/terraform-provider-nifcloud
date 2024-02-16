package dbinstance

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/rdb"
)

func flatten(d *schema.ResourceData, res *rdb.DescribeDBInstancesOutput) error {
	if res == nil || len(res.DBInstances) == 0 {
		d.SetId("")
		return nil
	}

	dbInstance := res.DBInstances[0]

	if nifcloud.ToString(dbInstance.DBInstanceIdentifier) != d.Id() {
		return fmt.Errorf("unable to find DB instance within: %#v", res.DBInstances)
	}

	if err := d.Set("accounting_type", dbInstance.NextMonthAccountingType); err != nil {
		return err
	}

	if err := d.Set("instance_class", dbInstance.DBInstanceClass); err != nil {
		return err
	}

	if err := d.Set("db_name", dbInstance.DBName); err != nil {
		return err
	}

	if err := d.Set("username", dbInstance.MasterUsername); err != nil {
		return err
	}

	if err := d.Set("engine", dbInstance.Engine); err != nil {
		return err
	}

	if err := d.Set("engine_version", dbInstance.EngineVersion); err != nil {
		return err
	}

	if err := d.Set("ca_cert_identifier", dbInstance.CACertificateIdentifier); err != nil {
		return err
	}

	if err := d.Set("allocated_storage", dbInstance.AllocatedStorage); err != nil {
		return err
	}

	if err := d.Set("storage_type", dbInstance.NiftyStorageType); err != nil {
		return err
	}

	if err := d.Set("identifier", dbInstance.DBInstanceIdentifier); err != nil {
		return err
	}

	if err := d.Set("availability_zone", dbInstance.AvailabilityZone); err != nil {
		return err
	}

	if err := d.Set("backup_retention_period", dbInstance.BackupRetentionPeriod); err != nil {
		return err
	}

	if binlogRetentionPeriod := dbInstance.BinlogRetentionPeriod; binlogRetentionPeriod != nil {
		if err := d.Set("custom_binlog_retention_period", true); err != nil {
			return err
		}
		if err := d.Set("binlog_retention_period", binlogRetentionPeriod); err != nil {
			return err
		}
	}

	if err := d.Set("backup_window", dbInstance.PreferredBackupWindow); err != nil {
		return err
	}

	if err := d.Set("maintenance_window", dbInstance.PreferredMaintenanceWindow); err != nil {
		return err
	}

	if err := d.Set("multi_az", dbInstance.MultiAZ); err != nil {
		return err
	}

	if err := d.Set("port", dbInstance.Endpoint.Port); err != nil {
		return err
	}

	if err := d.Set("publicly_accessible", dbInstance.PubliclyAccessible); err != nil {
		return err
	}

	if err := d.Set("db_security_group_name", dbInstance.DBSecurityGroups[0].DBSecurityGroupName); err != nil {
		return err
	}

	if err := d.Set("parameter_group_name", dbInstance.DBParameterGroups[0].DBParameterGroupName); err != nil {
		return err
	}

	if err := d.Set("address", dbInstance.Endpoint.Address); err != nil {
		return err
	}

	if err := d.Set("network_id", dbInstance.NiftyNetworkId); err != nil {
		return err
	}

	if err := d.Set("replicate_source_db", dbInstance.ReadReplicaSourceDBInstanceIdentifier); err != nil {
		return err
	}
	return nil
}
