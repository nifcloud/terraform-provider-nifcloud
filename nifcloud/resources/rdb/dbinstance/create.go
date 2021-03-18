package dbinstance

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).RDB

	if _, ok := d.GetOk("replicate_source_db"); ok {
		input := expandCreateDBInstanceReadReplicaInput(d)

		_, err := svc.CreateDBInstanceReadReplicaRequest(input).Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed creating db instance read replica: %s", err))
		}
	} else if _, ok := d.GetOk("snapshot_identifier"); ok {
		input := expandRestoreDBInstanceFromDBSnapshotInput(d)

		_, err := svc.RestoreDBInstanceFromDBSnapshotRequest(input).Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed creating db instance from snapshot: %s", err))
		}
	} else if _, ok := d.GetOk("restore_to_point_in_time"); ok {
		input := expandRestoreDBInstanceToPointInTimeInput(d)

		_, err := svc.RestoreDBInstanceToPointInTimeRequest(input).Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed creating db instance to point in time: %s", err))
		}
	} else {
		input := expandCreateDBInstanceInput(d)

		_, err := svc.CreateDBInstanceRequest(input).Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed creating db instance: %s", err))
		}
	}

	d.SetId(d.Get("identifier").(string))

	err := svc.WaitUntilDBInstanceAvailable(ctx, expandDescribeDBInstancesInput(d))
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed waiting for db instance to become ready: %s", err))
	}

	if d.Get("custom_binlog_retention_period").(bool) || d.Get("ca_cert_identifier").(string) != "" {
		input := expandModifyDBInstanceInput(d)

		req := svc.ModifyDBInstanceRequest(input)

		_, err := req.Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating db instance: %s", err))
		}

		err = svc.WaitUntilDBInstanceAvailable(ctx, expandDescribeDBInstancesInput(d))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed waiting for db instance to become ready: %s", err))
		}
	}
	return read(ctx, d, meta)
}
