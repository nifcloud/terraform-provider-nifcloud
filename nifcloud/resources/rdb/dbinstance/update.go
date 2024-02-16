package dbinstance

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/service/rdb"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).RDB
	deadline, _ := ctx.Deadline()

	// lintignore:R019
	if d.HasChanges(
		"accounting_type",
		"instance_class",
		"password",
		"ca_cert_identifier",
		"allocated_storage",
		"identifier",
		"backup_retention_period",
		"binlog_retention_period",
		"backup_window",
		"maintenance_window",
		"multi_az",
		"db_security_group_name",
		"parameter_group_name",
		"slave_private_address",
		"read_replica_private_address",
		"read_replica_identifier",
	) {
		input := expandModifyDBInstanceInput(d)

		_, err := svc.ModifyDBInstance(ctx, input)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating db instance: %s", err))

		}
		d.SetId(d.Get("identifier").(string))

		err = rdb.NewDBInstanceAvailableWaiter(svc).Wait(ctx, expandDescribeDBInstancesInput(d), time.Until(deadline))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed waiting for db instance to become ready: %s", err))
		}
	}

	if d.HasChanges(
		"slave_private_address",
		"master_private_address",
		"virtual_private_address",
		"network_id",
	) {

		input := expandModifyDBInstanceNetworkInput(d)

		_, err := svc.ModifyDBInstanceNetwork(ctx, input)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating db instance network: %s", err))
		}

		err = rdb.NewDBInstanceAvailableWaiter(svc).Wait(ctx, expandDescribeDBInstancesInput(d), time.Until(deadline))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed waiting for db instance to become ready: %s", err))
		}
	}
	return read(ctx, d, meta)
}
