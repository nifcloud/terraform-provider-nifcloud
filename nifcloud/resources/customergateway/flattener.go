package customergateway

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
)

func flatten(d *schema.ResourceData, res *computing.DescribeCustomerGatewaysResponse) error {
	if res == nil || len(res.CustomerGatewaySet) == 0 {
		d.SetId("")
		return nil
	}

	customerGateway := res.CustomerGatewaySet[0]

	if nifcloud.StringValue(customerGateway.CustomerGatewayId) != d.Id() {
		return fmt.Errorf("unable to find customer gateway within: %#v", res.CustomerGatewaySet)
	}

	if err := d.Set("customer_gateway_id", customerGateway.CustomerGatewayId); err != nil {
		return err
	}

	if err := d.Set("nifty_customer_gateway_name", customerGateway.NiftyCustomerGatewayName); err != nil {
		return err
	}

	if err := d.Set("state", customerGateway.State); err != nil {
		return err
	}

	if err := d.Set("ip_address", customerGateway.IpAddress); err != nil {
		return err
	}

	if err := d.Set("nifty_lan_side_ip_address", customerGateway.NiftyLanSideIpAddress); err != nil {
		return err
	}

	if err := d.Set("nifty_lan_side_cidr_block", customerGateway.NiftyLanSideCidrBlock); err != nil {
		return err
	}

	if err := d.Set("nifty_customer_gateway_description", customerGateway.NiftyCustomerGatewayDescription); err != nil {
		return err
	}
	return nil
}
