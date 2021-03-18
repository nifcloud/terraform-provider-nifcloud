package vpnconnection

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/validator"
)

const description = "Provides a vpn connection resource."

// New returns the nifcloud_vpn_connection resource schema.
func New() *schema.Resource {
	return &schema.Resource{
		Description: description,
		Schema:      newSchema(),

		CreateContext: create,
		ReadContext:   read,
		DeleteContext: delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func newSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:         schema.TypeString,
			Description:  "The type of vpn connection.",
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"IPsec", "L2TPv3 / IPsec", "IPsec VTI"}, false),
		},
		"vpn_gateway_id": {
			Type:          schema.TypeString,
			Description:   "The id for the vpn gateway. Cannot be specified with vpn_gateway_name.",
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"vpn_gateway_name"},
		},
		"vpn_gateway_name": {
			Type:          schema.TypeString,
			Description:   "The name for the vpn gateway. Cannot be specified with vpn_gateway_id.",
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"vpn_gateway_id"},
		},
		"customer_gateway_id": {
			Type:          schema.TypeString,
			Description:   "The id for the customer gateway. Cannot be specified with customer_gateway_name.",
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"customer_gateway_name"},
		},
		"customer_gateway_name": {
			Type:          schema.TypeString,
			Description:   "The name for the customer gateway. Cannot be specified with customer_gateway_id.",
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"customer_gateway_id"},
		},
		"tunnel_type": {
			Type:         schema.TypeString,
			Description:  "The type of vpn connection tunnel.",
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"L2TPv3"}, false),
		},
		"tunnel_mode": {
			Type:         schema.TypeString,
			Description:  "The mode of vpn connection tunnel; Unmanaged or Managed.",
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"Unmanaged", "Managed"}, false),
		},
		"tunnel_encapsulation": {
			Type:         schema.TypeString,
			Description:  "The encapsulation of vpn connection.",
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"IP", "UDP"}, false),
		},
		"tunnel_id": {
			Type:        schema.TypeString,
			Description: "The id for the vpn gateway tunnel.",
			Optional:    true,
			ForceNew:    true,
		},
		"tunnel_peer_id": {
			Type:        schema.TypeString,
			Description: "The id for the customer gateway tunnel.",
			Optional:    true,
			ForceNew:    true,
		},
		"tunnel_session_id": {
			Type:        schema.TypeString,
			Description: "The session id for the vpn gateway tunnel.",
			Optional:    true,
			ForceNew:    true,
		},
		"tunnel_peer_session_id": {
			Type:        schema.TypeString,
			Description: "The session id for the customer gateway tunnel.",
			Optional:    true,
			ForceNew:    true,
		},
		"tunnel_source_port": {
			Type:        schema.TypeString,
			Description: "The port for the vpn gateway tunnel.",
			Optional:    true,
			ForceNew:    true,
		},
		"tunnel_destination_port": {
			Type:        schema.TypeString,
			Description: "The port for the customer gateway tunnel.",
			Optional:    true,
			ForceNew:    true,
		},
		"ipsec_config_encryption_algorithm": {
			Type:         schema.TypeString,
			Description:  "The encryption algorithm for IPsec config.",
			Optional:     true,
			ForceNew:     true,
			Default:      "AES128",
			ValidateFunc: validation.StringInSlice([]string{"AES128", "AES256", "3DES"}, false),
		},
		"ipsec_config_hash_algorithm": {
			Type:         schema.TypeString,
			Description:  "The hash algorithm for IPsec config.",
			Optional:     true,
			ForceNew:     true,
			Default:      "SHA1",
			ValidateFunc: validation.StringInSlice([]string{"SHA1", "MD5", "SHA256", "SHA384", "SHA512"}, false),
		},
		"ipsec_config_pre_shared_key": {
			Type:         schema.TypeString,
			Description:  "The pre shared key for IPsec config.",
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(1, 64),
		},
		"ipsec_config_internet_key_exchange": {
			Type:         schema.TypeString,
			Description:  "The IKE protocol for IPsec config.",
			Optional:     true,
			ForceNew:     true,
			Default:      "IKEv1",
			ValidateFunc: validation.StringInSlice([]string{"IKEv1", "IKEv2"}, false),
		},
		"ipsec_config_internet_key_exchange_lifetime": {
			Type:         schema.TypeInt,
			Description:  "The IKE SA expiration seconds for IPsec config.",
			Optional:     true,
			ForceNew:     true,
			Default:      28800,
			ValidateFunc: validation.IntBetween(30, 86400),
		},
		"ipsec_config_encapsulating_security_payload_lifetime": {
			Type:         schema.TypeInt,
			Description:  "The ESP SA expiration seconds for IPsec config.",
			Optional:     true,
			ForceNew:     true,
			Default:      3600,
			ValidateFunc: validation.IntBetween(30, 86400),
		},
		"ipsec_config_diffie_hellman_group": {
			Type:         schema.TypeInt,
			Description:  "The Diffie-Hellman Group for IKE and PFS.",
			Optional:     true,
			ForceNew:     true,
			Default:      2,
			ValidateFunc: validation.IntBetween(2, 26),
		},
		"mtu": {
			Type:        schema.TypeString,
			Description: "The MTU size for vpn connection.",
			Optional:    true,
			ForceNew:    true,
			Default:     "1500",
		},
		"description": {
			Type:             schema.TypeString,
			Description:      "The vpn connection description.",
			Optional:         true,
			ForceNew:         true,
			ValidateDiagFunc: validator.StringRuneCountBetween(0, 500),
		},
		"vpn_connection_id": {
			Type:        schema.TypeString,
			Description: "The id of the vpn connection.",
			Computed:    true,
		},
	}
}
