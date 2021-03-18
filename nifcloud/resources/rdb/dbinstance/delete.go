package dbinstance

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/rdb"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).RDB

	if v, ok := d.GetOk("read_replica_identifier"); ok {
		input := &rdb.DeleteDBInstanceInput{
			DBInstanceIdentifier: nifcloud.String(v.(string)),
			SkipFinalSnapshot:    nifcloud.Bool(true),
		}
		req := svc.DeleteDBInstanceRequest(input)

		_, err := req.Send(ctx)
		if err != nil {
			var awsErr awserr.Error
			if errors.As(err, &awsErr) && awsErr.Code() == "Client.InvalidParameterNotFound.DBInstance" {
				d.SetId("")
				return nil
			}
			return diag.FromErr(fmt.Errorf("failed deleting: %s", err))
		}

		err = svc.WaitUntilDBInstanceAvailable(ctx, expandDescribeDBInstancesInput(d))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed deleting for wait until available after read replica deleted error: %s", err))
		}
	}

	input := expandDeleteDBInstanceInput(d)

	req := svc.DeleteDBInstanceRequest(input)

	_, err := req.Send(ctx)
	if err != nil {
		var awsErr awserr.Error
		if errors.As(err, &awsErr) && awsErr.Code() == "Client.InvalidParameterNotFound.DBInstance" {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("failed deleting: %s", err))
	}

	err = svc.WaitUntilDBInstanceDeleted(ctx, expandDescribeDBInstancesInput(d))
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed deleting for wait until deleted db instance error: %s", err))
	}

	d.SetId("")
	return nil
}
