package dbinstance

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).RDB

	if d.HasChanges(
		"accounting_type",
		"instance_class",
		"password",
		"ca_cert_identifier",
		"allocated_storage",
		"identifier",
		"backup_retention_period",
		"backup_window",
		"maintenance_window",
		"multi_az",
		"multi_az_type",
		"db_security_group_name",
		"parameter_group_name",
		"slave_private_address",
		"read_replica_private_address",
		"read_replica_identifier",
	) {
		input := expandModifyDBInstanceInput(d)

		req := svc.ModifyDBInstanceRequest(input)

		_, err := req.Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating db instance: %s", err))

		}
		d.SetId(d.Get("identifier").(string))

		err = svc.WaitUntilDBInstanceAvailable(ctx, expandDescribeDBInstancesInput(d))
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

		req := svc.ModifyDBInstanceNetworkRequest(input)

		_, err := req.Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating db instance network: %s", err))
		}

		err = svc.WaitUntilDBInstanceAvailable(ctx, expandDescribeDBInstancesInput(d))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed waiting for db instance to become ready: %s", err))
		}
	}
	return read(ctx, d, meta)
}
