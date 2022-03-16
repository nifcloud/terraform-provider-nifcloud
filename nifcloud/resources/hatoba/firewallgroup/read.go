package firewallgroup

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
	input := expandGetFirewallGroupInput(d)
	svc := meta.(*client.Client).Hatoba
	res, err := svc.GetFirewallGroup(ctx, input)
	if err != nil {
		var awsErr smithy.APIError
		if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameterNotFound.FirewallGroup" {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("failed reading Hatoba firewall group: %s", err))
	}

	if err := flatten(d, res); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
