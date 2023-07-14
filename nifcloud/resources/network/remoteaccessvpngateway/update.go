package remoteaccessvpngateway

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

const waiterInitialDelay = 3

func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).Computing

	if d.IsNewResource() {
		if d := waitForRemoteAccessVpnGatewayAvailable(ctx, d, svc); d != nil {
			return d
		}
	}

	if d.HasChange("accounting_type") {
		input := expandModifyRemoteAccessVpnGatewayAttributeInputForAccountingType(d)

		_, err := svc.ModifyRemoteAccessVpnGatewayAttribute(ctx, input)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating remote access vpn Gateway accounting_type: %s", err))
		}

		if d := waitForRemoteAccessVpnGatewayAvailable(ctx, d, svc); d != nil {
			return d
		}
	}

	if d.HasChange("description") {
		input := expandModifyRemoteAccessVpnGatewayAttributeInputForDescription(d)

		_, err := svc.ModifyRemoteAccessVpnGatewayAttribute(ctx, input)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating remote access vpn Gateway description: %s", err))
		}

		if d := waitForRemoteAccessVpnGatewayAvailable(ctx, d, svc); d != nil {
			return d
		}
	}

	if d.HasChange("name") {
		input := expandModifyRemoteAccessVpnGatewayAttributeInputForRemoteAccessVpnGatewayName(d)

		_, err := svc.ModifyRemoteAccessVpnGatewayAttribute(ctx, input)

		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating remote access vpn Gateway name %s", err))
		}

		if d := waitForRemoteAccessVpnGatewayAvailable(ctx, d, svc); d != nil {
			return d
		}
	}

	if d.HasChange("type") {
		input := expandModifyRemoteAccessVpnGatewayAttributeInputForType(d)

		_, err := svc.ModifyRemoteAccessVpnGatewayAttribute(ctx, input)

		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating remote access vpn Gateway type: %s", err))
		}

		if d := waitForRemoteAccessVpnGatewayAvailable(ctx, d, svc); d != nil {
			return d
		}
	}

	if d.HasChange("ssl_certificate_id") {
		input := expandSetRemoteAccessVpnGatewaySSLCertificateInput(d)

		_, err := svc.SetRemoteAccessVpnGatewaySSLCertificate(ctx, input)

		if err != nil {
			return diag.FromErr(fmt.Errorf("failed updating remote access vpn gateway ssl certificate: %s", err))
		}

		if d := waitForRemoteAccessVpnGatewayAvailable(ctx, d, svc); d != nil {
			return d
		}
	}

	if d.HasChange("ca_certificate_id") {
		before, after := d.GetChange("ca_certificate_id")
		if before != "" && after == "" {
			input := expandUnsetRemoteAccessVpnGatewayCACertificateInput(d)

			_, err := svc.UnsetRemoteAccessVpnGatewayCACertificate(ctx, input)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed unset ca certificate: %s", err))
			}
		} else if _, ok := d.GetOk("ca_certificate_id"); ok {
			input := expandSetRemoteAccessVpnGatewayCACertificateInput(d)

			_, err := svc.SetRemoteAccessVpnGatewayCACertificate(ctx, input)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed updating remote access vpn gateway ca certificate: %s", err))
			}
		}

		if d := waitForRemoteAccessVpnGatewayAvailable(ctx, d, svc); d != nil {
			return d
		}
	}

	if d.HasChange("user") {
		o, n := d.GetChange("user")
		ors := o.(*schema.Set).Difference(n.(*schema.Set))
		nrs := n.(*schema.Set).Difference(o.(*schema.Set))

		// Now first loop through all the old users and delete any obsolete ones
		for _, user := range ors.List() {
			input := expandDeleteRemoteAccessVpnGatewayUsersInput(d, user.(map[string]interface{}))

			_, err := svc.DeleteRemoteAccessVpnGatewayUsers(ctx, input)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed updating remote access vpn gateway to delete user: %s", err))
			}

			if d := waitForRemoteAccessVpnGatewayAvailable(ctx, d, svc); d != nil {
				return d
			}
		}

		// Make sure we save the state of the currently configured rules
		users := o.(*schema.Set).Intersection(n.(*schema.Set))
		if err := d.Set("user", users); err != nil {
			return diag.FromErr(err)
		}

		// Then loop through all the newly configured users and create them
		for _, user := range nrs.List() {
			input := expandCreateRemoteAccessVpnGatewayUsersInput(d, user.(map[string]interface{}))

			_, err := svc.CreateRemoteAccessVpnGatewayUsers(ctx, input)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed updatig remote access vpn gateway to create user: %s", err))
			}

			if d := waitForRemoteAccessVpnGatewayAvailable(ctx, d, svc); d != nil {
				return d
			}

			users.Add(user)
			if err := d.Set("user", users); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return read(ctx, d, meta)
}

func waitForRemoteAccessVpnGatewayAvailable(ctx context.Context, d *schema.ResourceData, svc *computing.Client) diag.Diagnostics {
	// lintignore:R018
	time.Sleep(waiterInitialDelay * time.Second)
	deadline, _ := ctx.Deadline()

	if err := computing.NewRemoteAccessVpnGatewayAvailableWaiter(svc).Wait(ctx, expandDescribeRemoteAccessVpnGatewaysInput(d), time.Until(deadline)); err != nil {
		return diag.FromErr(fmt.Errorf("failed waiting for remote access vpn Gateway available: %s", err))
	}

	return nil
}
