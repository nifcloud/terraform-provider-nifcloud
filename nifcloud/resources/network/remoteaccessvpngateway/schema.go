package remoteaccessvpngateway

import (
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/validator"
)

const description = "Provides a remote access vpn gateway resource."

// New returns the nifcloud_remote_access_vpn_gateway resource schema.
func New() *schema.Resource {
	return &schema.Resource{
		Description: description,
		Schema:      newSchema(),

		CreateContext: create,
		ReadContext:   read,
		UpdateContext: update,
		DeleteContext: delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(5 * time.Minute),
			Create:  schema.DefaultTimeout(30 * time.Minute),
			Update:  schema.DefaultTimeout(60 * time.Minute),
			Delete:  schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func newSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"accounting_type": {
			Type:         schema.TypeString,
			Description:  "Accounting type. (1: monthly, 2: pay per use).",
			Optional:     true,
			Default:      "2",
			ValidateFunc: validation.StringInSlice([]string{"1", "2"}, false),
		},
		"availability_zone": {
			Type:        schema.TypeString,
			Description: "The availability zone.",
			Optional:    true,
			Computed:    true,
			ForceNew:    true,
		},
		"ca_certificate_id": {
			Type:        schema.TypeString,
			Description: "The ID of ca certificate.",
			Optional:    true,
		},
		"cipher_suite": {
			Type:        schema.TypeList,
			Description: "The Cipher suite; can be specified one of `AES128-GCM-SHA256` `AES256-GCM-SHA384` `ECDHE-RSA-AES128-GCM-SHA256` `ECDHE-RSA-AES256-GCM-SHA384`",
			Required:    true,
			ForceNew:    true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{
					"AES128-GCM-SHA256",
					"AES256-GCM-SHA384",
					"ECDHE-RSA-AES128-GCM-SHA256",
					"ECDHE-RSA-AES256-GCM-SHA384",
				}, false),
			},
		},
		"client_config": {
			Type:        schema.TypeString,
			Description: "The base64 encoding remote access vpn gateway client config.",
			Computed:    true,
		},
		"description": {
			Type:             schema.TypeString,
			Description:      "The remote access vpn gateway description.",
			Optional:         true,
			ValidateDiagFunc: validator.StringRuneCountBetween(0, 500),
		},
		"name": {
			Type:        schema.TypeString,
			Description: "The remote access vpn gateway name.",
			Optional:    true,
			Computed:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 15),
				validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z]+$`), "Enter the name within 1-15 characters [0-9a-zA-Z]."),
			),
		},
		"network_interface": {
			Type:     schema.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"ip_address": {
						Type:         schema.TypeString,
						Description:  "The IP address of the network interface.",
						Optional:     true,
						ValidateFunc: validation.IsIPAddress,
					},
					"network_id": {
						Type:        schema.TypeString,
						Description: "The ID of the network to attach private lan network.",
						Optional:    true,
					},
				},
			},
		},
		"pool_network_cidr": {
			Type:             schema.TypeString,
			Description:      "The cidr of pool network; can be specified in the range of /16 to /27.",
			ForceNew:         true,
			Required:         true,
			ValidateDiagFunc: validator.CIDRNetworkAddress,
		},
		"remote_access_vpn_gateway_id": {
			Type:        schema.TypeString,
			Description: "The unique ID of the remote access vpn gateway.",
			Computed:    true,
		},
		"ssl_certificate_id": {
			Type:        schema.TypeString,
			Description: "The ID of ssl certificate.",
			Required:    true,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "The type of the remote access vpn gateway.",
			Optional:    true,
			Default:     "small",
			ValidateFunc: validation.StringInSlice([]string{
				"small",
				"medium",
				"large",
			}, false),
		},
		"user": {
			Type:        schema.TypeSet,
			Description: "List of the remote access vpn gateway user.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Description: "The name of remote access vpn gateway user.",
						Required:    true,
					},
					"password": {
						Type:        schema.TypeString,
						Description: "The password of remote access vpn gateway user.",
						Sensitive:   true,
						Required:    true,
					},
					"description": {
						Type:        schema.TypeString,
						Description: "The remote access vpn gateway user description.",
						Optional:    true,
					},
				},
			},
		},
	}
}
