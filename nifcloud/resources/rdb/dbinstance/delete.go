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

	describeDBInstancesInput := expandDescribeDBInstancesInput(d)
	res, err := svc.DescribeDBInstancesRequest(describeDBInstancesInput).Send(ctx)
	if err != nil {
		var awsErr awserr.Error
		if errors.As(err, &awsErr) && awsErr.Code() == "Client.InvalidParameterNotFound.DBInstance" {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("failed deleting: %s", err))
	}
	if len(res.DBInstances) == 0 {
		d.SetId("")
		return nil
	}

	// Delete Read Replica
	for _, rr := range res.DBInstances[0].ReadReplicaDBInstanceIdentifiers {
		input := &rdb.DeleteDBInstanceInput{
			DBInstanceIdentifier: nifcloud.String(rr),
			SkipFinalSnapshot:    nifcloud.Bool(true),
		}

		_, err := svc.DeleteDBInstanceRequest(input).Send(ctx)
		if err != nil {
			var awsErr awserr.Error
			if errors.As(err, &awsErr) && awsErr.Code() == "Client.InvalidParameterNotFound.DBInstance" {
				continue
			}
			return diag.FromErr(fmt.Errorf("failed deleting: %s", err))
		}

		err = svc.WaitUntilDBInstanceAvailable(ctx, describeDBInstancesInput)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed deleting for wait until available after read replica deleted error: %s", err))
		}
	}

	// Delete DB Instance
	input := expandDeleteDBInstanceInput(d)
	_, err = svc.DeleteDBInstanceRequest(input).Send(ctx)
	if err != nil {
		var awsErr awserr.Error
		if errors.As(err, &awsErr) && awsErr.Code() == "Client.InvalidParameterNotFound.DBInstance" {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("failed deleting: %s", err))
	}

	err = svc.WaitUntilDBInstanceDeleted(ctx, describeDBInstancesInput)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed deleting for wait until deleted db instance error: %s", err))
	}

	d.SetId("")
	return nil
}
