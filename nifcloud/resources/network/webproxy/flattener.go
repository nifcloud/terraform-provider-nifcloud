package webproxy

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
)

func flatten(d *schema.ResourceData, res *computing.NiftyDescribeWebProxiesOutput) error {
	if res == nil || len(res.WebProxy) == 0 {
		d.SetId("")
		return nil
	}

	webProxy := res.WebProxy[0]

	if nifcloud.ToString(webProxy.RouterId) != d.Id() {
		return fmt.Errorf("unable to find web proxy within: %#v", res.WebProxy)
	}

	if err := d.Set("router_name", webProxy.RouterName); err != nil {
		return err
	}

	if err := d.Set("router_id", webProxy.RouterId); err != nil {
		return err
	}

	if err := d.Set("description", webProxy.Description); err != nil {
		return err
	}

	if err := d.Set("listen_interface_network_id", webProxy.ListenInterface.NetworkId); err != nil {
		return err
	}

	if err := d.Set("listen_interface_network_name", webProxy.ListenInterface.NetworkName); err != nil {
		return err
	}

	if err := d.Set("listen_port", webProxy.ListenPort); err != nil {
		return err
	}
	if err := d.Set("bypass_interface_network_id", webProxy.BypassInterface.NetworkId); err != nil {
		return err
	}

	if err := d.Set("bypass_interface_network_name", webProxy.BypassInterface.NetworkName); err != nil {
		return err
	}

	if err := d.Set("name_server", webProxy.Option.NameServer); err != nil {
		return err
	}
	return nil
}
