package cluster

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/smithy-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandGetClusterInput(d)
	svc := meta.(*client.Client).Hatoba
	res, err := svc.GetCluster(ctx, input)

	if err != nil {
		var awsErr smithy.APIError
		if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameterNotFound.Cluster" {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("failed reading Hatoba cluster: %s", err))
	}

	if err := flatten(d, res); err != nil {
		return diag.FromErr(err)
	}

	getClusterCredentialsRes, err := svc.GetClusterCredentials(ctx,
		expandGetClusterCredentialsInput(d),
	)

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed reading Hatoba cluster credentials: %s", err))
	}

	if err := flattenCredentials(d, getClusterCredentialsRes); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
