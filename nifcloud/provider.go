package nifcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/datasources/image"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/customergateway"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/dhcpconfig"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/dhcpoption"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/elasticip"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/elb"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/elblistener"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/instance"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/keypair"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/loadbalancer"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/nattable"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/privatelan"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/router"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/routetable"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/securitygroup"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/securitygrouprule"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/sslcertificate"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/volume"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/webproxy"
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
			"nifcloud_customer_gateway":    customergateway.New(),
			"nifcloud_dhcp_config":         dhcpconfig.New(),
			"nifcloud_dhcp_option":         dhcpoption.New(),
			"nifcloud_elastic_ip":          elasticip.New(),
			"nifcloud_elb":                 elb.New(),
			"nifcloud_elb_listener":        elblistener.New(),
			"nifcloud_instance":            instance.New(),
			"nifcloud_key_pair":            keypair.New(),
			"nifcloud_nat_table":           nattable.New(),
			"nifcloud_load_balancer":       loadbalancer.New(),
			"nifcloud_private_lan":         privatelan.New(),
			"nifcloud_router":              router.New(),
			"nifcloud_route_table":         routetable.New(),
			"nifcloud_security_group":      securitygroup.New(),
			"nifcloud_security_group_rule": securitygrouprule.New(),
			"nifcloud_ssl_certificate":     sslcertificate.New(),
			"nifcloud_volume":              volume.New(),
			"nifcloud_web_proxy":           webproxy.New(),
		},
	}
}
