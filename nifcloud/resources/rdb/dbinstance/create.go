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

func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).RDB
	deadline, _ := ctx.Deadline()

	if _, ok := d.GetOk("replicate_source_db"); ok {
		input := expandCreateDBInstanceReadReplicaInput(d)

		_, err := svc.CreateDBInstanceReadReplica(ctx, input)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed creating db instance read replica: %s", err))
		}
	} else if _, ok := d.GetOk("snapshot_identifier"); ok {
		input := expandRestoreDBInstanceFromDBSnapshotInput(d)

		_, err := svc.RestoreDBInstanceFromDBSnapshot(ctx, input)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed creating db instance from snapshot: %s", err))
		}
	} else if _, ok := d.GetOk("restore_to_point_in_time"); ok {
		input := expandRestoreDBInstanceToPointInTimeInput(d)

		_, err := svc.RestoreDBInstanceToPointInTime(ctx, input)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed creating db instance to point in time: %s", err))
		}
	} else {
		input := expandCreateDBInstanceInput(d)

		_, err := svc.CreateDBInstance(ctx, input)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed creating db instance: %s", err))
		}
	}

	d.SetId(d.Get("identifier").(string))

	err := rdb.NewDBInstanceAvailableWaiter(svc).Wait(ctx, expandDescribeDBInstancesInput(d), time.Until(deadline))
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed waiting for db instance to become ready: %s", err))
	}

	if d.Get("custom_binlog_retention_period").(bool) || d.Get("ca_cert_identifier").(string) != "" {
		input := expandModifyDBInstanceInput(d)

		_, err := svc.ModifyDBInstance(ctx, input)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating db instance: %s", err))
		}

		err = rdb.NewDBInstanceAvailableWaiter(svc).Wait(ctx, expandDescribeDBInstancesInput(d), time.Until(deadline))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed waiting for db instance to become ready: %s", err))
		}
	}
	return read(ctx, d, meta)
}
