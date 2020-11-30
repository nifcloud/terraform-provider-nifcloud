package nifcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/datasources/image"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/elasticip"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/instance"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/keypair"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/privatelan"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/securitygroup"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/securitygrouprule"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/sslcertificate"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/resources/volume"
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
			"nifcloud_elastic_ip":          elasticip.New(),
			"nifcloud_instance":            instance.New(),
			"nifcloud_key_pair":            keypair.New(),
			"nifcloud_private_lan":         privatelan.New(),
			"nifcloud_security_group":      securitygroup.New(),
			"nifcloud_security_group_rule": securitygrouprule.New(),
			"nifcloud_ssl_certificate":     sslcertificate.New(),
			"nifcloud_volume":              volume.New(),
		},
	}
}
