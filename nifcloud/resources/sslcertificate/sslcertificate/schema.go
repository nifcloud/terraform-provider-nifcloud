package sslcertificate

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/validator"
)

const description = "Provides a ssl certificate resource."

// New returns the nifcloud_ssl_certificate resource schema.
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
		},
	}
}

func newSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"certificate": {
			Type:        schema.TypeString,
			Description: "The certificate material.",
			Required:    true,
			ForceNew:    true,
		},
		"key": {
			Type:        schema.TypeString,
			Description: "The private key material.",
			Required:    true,
			Sensitive:   true,
			ForceNew:    true,
		},
		"ca": {
			Type:        schema.TypeString,
			Description: "The certificate authority material.",
			ForceNew:    true,
			Optional:    true,
		},
		"description": {
			Type:             schema.TypeString,
			Description:      "The SSL certificate description.",
			Optional:         true,
			ValidateDiagFunc: validator.StringRuneCountBetween(0, 40),
		},
		"fqdn_id": {
			Type:        schema.TypeString,
			Description: "The unique identifier for the certificate.",
			Computed:    true,
		},
		"fqdn": {
			Type:        schema.TypeString,
			Description: "The name for the certificate.",
			Computed:    true,
		},
	}
}
