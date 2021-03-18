package image

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const description = "Use this data source to get the ID of a image for use in nifcloud_instance resources."

// New returns the nifcloud_image data source schema.
func New() *schema.Resource {
	return &schema.Resource{
		Description: description,
		Schema:      newSchema(),

		ReadContext: read,

		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func newSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"image_name": {
			Type:        schema.TypeString,
			Description: "The name of image.",
			Required:    true,
		},
		"owner": {
			Type:         schema.TypeString,
			Description:  "The image owner; valid values: `niftycloud` (standard image) `self` (current account) `other` (other user).",
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"niftycloud", "self", "other"}, false),
		},
		"image_id": {
			Type:        schema.TypeString,
			Description: "The id of image.",
			Computed:    true,
		},
	}
}
