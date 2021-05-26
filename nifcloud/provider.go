package nifcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/datasources/computing/image"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/computing/elasticip"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/computing/instance"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/computing/keypair"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/computing/networkinterface"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/computing/securitygroup"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/computing/securitygrouprule"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/computing/volume"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/hatoba/firewallgroup"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/nas/nasinstance"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/nas/nassecuritygroup"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/network/customergateway"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/network/dhcpconfig"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/network/dhcpoption"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/network/elb"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/network/elblistener"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/network/loadbalancer"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/network/loadbalancerlistener"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/network/nattable"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/network/privatelan"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/network/router"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/network/routetable"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/network/vpnconnection"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/network/vpngateway"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/network/webproxy"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/rdb/dbinstance"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/rdb/dbparametergroup"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/rdb/dbsecuritygroup"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/sslcertificate/sslcertificate"
)

// Provider returns a schema.Provider.
func Provider() *schema.Provider {
	return &schema.Provider{
		ConfigureContextFunc: configure,
		Schema: map[string]*schema.Schema{
			"access_key": {
				Description: "This is the NIFCLOUD access key. It must be provided, but it can also be sourced from the `NIFCLOUD_ACCESS_KEY_ID` env var.",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("NIFCLOUD_ACCESS_KEY_ID", nil),
			},
			"secret_key": {
				Description: "This is the NIFCLOUD secret key. It must be provided, but it can also be sourced from the `NIFCLOUD_SECRET_ACCESS_KEY` env var.",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("NIFCLOUD_SECRET_ACCESS_KEY", nil),
			},
			"region": {
				Description: "This is the NIFCLOUD region. It must be provided, but it can also be sourced from the `NIFCLOUD_DEFAULT_REGION` env var.",
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NIFCLOUD_DEFAULT_REGION", nil),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"nifcloud_image": image.New(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"nifcloud_customer_gateway":       customergateway.New(),
			"nifcloud_db_instance":            dbinstance.New(),
			"nifcloud_db_parameter_group":     dbparametergroup.New(),
			"nifcloud_db_security_group":      dbsecuritygroup.New(),
			"nifcloud_dhcp_config":            dhcpconfig.New(),
			"nifcloud_dhcp_option":            dhcpoption.New(),
			"nifcloud_elastic_ip":             elasticip.New(),
			"nifcloud_elb":                    elb.New(),
			"nifcloud_elb_listener":           elblistener.New(),
			"nifcloud_instance":               instance.New(),
			"nifcloud_key_pair":               keypair.New(),
			"nifcloud_nas_instance":           nasinstance.New(),
			"nifcloud_nas_security_group":     nassecuritygroup.New(),
			"nifcloud_nat_table":              nattable.New(),
			"nifcloud_network_interface":      networkinterface.New(),
			"nifcloud_load_balancer":          loadbalancer.New(),
			"nifcloud_load_balancer_listener": loadbalancerlistener.New(),
			"nifcloud_private_lan":            privatelan.New(),
			"nifcloud_router":                 router.New(),
			"nifcloud_route_table":            routetable.New(),
			"nifcloud_security_group":         securitygroup.New(),
			"nifcloud_security_group_rule":    securitygrouprule.New(),
			"nifcloud_ssl_certificate":        sslcertificate.New(),
			"nifcloud_volume":                 volume.New(),
			"nifcloud_vpn_connection":         vpnconnection.New(),
			"nifcloud_vpn_gateway":            vpngateway.New(),
			"nifcloud_web_proxy":              webproxy.New(),
			"nifcloud_hatoba_firewall_group":  firewallgroup.New(),
		},
	}
}
