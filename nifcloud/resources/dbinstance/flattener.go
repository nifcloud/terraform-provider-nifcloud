package dbinstance

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/rdb"
)

func flatten(d *schema.ResourceData, res *rdb.DescribeDBInstancesResponse) error {
	if res == nil || len(res.DBInstances) == 0 {
		d.SetId("")
		return nil
	}

	dbInstance := res.DBInstances[0]

	if nifcloud.StringValue(dbInstance.DBInstanceIdentifier) != d.Id() {
		return fmt.Errorf("unable to find DB instance within: %#v", res.DBInstances)
	}

	if err := d.Set("accounting_type", strconv.FormatInt(nifcloud.Int64Value(dbInstance.NextMonthAccountingType), 10)); err != nil {
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

	if allocatedStorage, err := strconv.Atoi(nifcloud.StringValue(dbInstance.AllocatedStorage)); err == nil {
		if err := d.Set("allocated_storage", allocatedStorage); err != nil {
			return err
		}
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

	if backupRetentionPeriod, err := strconv.Atoi(nifcloud.StringValue(dbInstance.BackupRetentionPeriod)); err == nil {
		if err := d.Set("backup_retention_period", backupRetentionPeriod); err != nil {
			return err
		}
	}

	if binlogRetentionPeriod, err := strconv.Atoi(nifcloud.StringValue(dbInstance.BinlogRetentionPeriod)); err == nil {
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

	if multiAZ, err := strconv.ParseBool(nifcloud.StringValue(dbInstance.MultiAZ)); err == nil {
		if err := d.Set("multi_az", multiAZ); err != nil {
			return err
		}
	}

	if multiAZType, err := strconv.Atoi(nifcloud.StringValue(dbInstance.NiftyMultiAZType)); err == nil {
		if err := d.Set("multi_az_type", multiAZType); err != nil {
			return err
		}
	}

	if port, err := strconv.Atoi(nifcloud.StringValue(dbInstance.Endpoint.Port)); err == nil {
		if err := d.Set("port", port); err != nil {
			return err
		}
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
