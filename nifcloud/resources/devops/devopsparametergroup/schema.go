package devopsparametergroup

import (
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/validator"
	"golang.org/x/exp/slices"
)

const description = "Provides a DevOps parameter group resource."

// New returns the nifcloud_devops_parameter_group resource schema.
func New() *schema.Resource {
	return &schema.Resource{
		Description: description,
		Schema:      newSchema(),

		CreateContext: createParameterGroup,
		ReadContext:   readParameterGroup,
		UpdateContext: updateParameterGroup,
		DeleteContext: deleteParameterGroup,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func newSchema() map[string]*schema.Schema {
	sensitiveFields := []string{"smtp_password"}

	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the DevOps parameter group.",
			Required:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 63),
				validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z\-]+$`), "Enter a name within 1-63 alphanumeric lowercase characters and hyphens. Hyphens cannot be used at the beginning or end."),
			),
		},
		"description": {
			Type:             schema.TypeString,
			Description:      "Description of the DevOps parameter group.",
			Optional:         true,
			ValidateDiagFunc: validator.StringRuneCountBetween(0, 255),
		},
		"parameter": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Description: "The name of the parameter.",
						Required:    true,
						ValidateDiagFunc: func(v any, p cty.Path) diag.Diagnostics {
							value := v.(string)
							var diags diag.Diagnostics
							if slices.Contains(sensitiveFields, value) {
								diag := diag.Diagnostic{
									Severity: diag.Error,
									Summary:  "wrong value",
									Detail:   fmt.Sprintf("%q has a sensitive value. use \"sensitive_parameter\" instead.", value),
								}
								diags = append(diags, diag)
							}
							return diags
						},
					},
					"value": {
						Type:        schema.TypeString,
						Description: "The value of the parameter.",
						Required:    true,
					},
				},
			},
		},
		"sensitive_parameter": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Description: "The name of the parameter. Valid value is `smtp_password`.",
						Required:    true,
						ValidateDiagFunc: func(v any, p cty.Path) diag.Diagnostics {
							value := v.(string)
							var diags diag.Diagnostics
							if !slices.Contains(sensitiveFields, value) {
								diag := diag.Diagnostic{
									Severity: diag.Error,
									Summary:  "wrong value",
									Detail:   fmt.Sprintf("found misuse of \"sensitive_parameter\" for %q. use \"parameter\" instead.", value),
								}
								diags = append(diags, diag)
							}
							return diags
						},
					},
					"value": {
						Type:        schema.TypeString,
						Description: "The value of the parameter.",
						Required:    true,
						Sensitive:   true,
					},
				},
			},
		},
	}
}
