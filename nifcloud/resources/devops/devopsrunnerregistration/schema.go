package devopsrunnerregistration

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const description = "Provides a DevOps Runner registration resource."

// New returns the nifcloud_devops_runner_registration resource schema.
func New() *schema.Resource {
	return &schema.Resource{
		Description: description,
		Schema:      newSchema(),

		CreateContext: createRunnerRegistration,
		ReadContext:   readRunnerRegistration,
		UpdateContext: updateRunnerRegistration,
		DeleteContext: deleteRunnerRegistration,

		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				importParts, err := validateRunnerRegistrationImportString(d.Id())
				if err != nil {
					return nil, err
				}
				if err := populateRunnerRegistrationFromImport(d, importParts); err != nil {
					return nil, err
				}
				return []*schema.ResourceData{d}, nil
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func newSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Description: "The name of the DevOps Runner registration.",
			Computed:    true,
		},
		"runner_name": {
			Type:        schema.TypeString,
			Description: "The name of the DevOps Runner to register.",
			Required:    true,
			ForceNew:    true,
		},
		"gitlab_url": {
			Type:        schema.TypeString,
			Description: "GitLab URL.",
			Required:    true,
			ForceNew:    true,
			StateFunc: func(i interface{}) string {
				v := i.(string)
				// NIFCLOUD DevOps Runner API returns URLs with a trailing slash.
				if !strings.HasSuffix(v, "/") {
					return v + "/"
				}
				return v
			},
		},
		"parameter_group_name": {
			Type:        schema.TypeString,
			Description: "The name of the DevOps Runner parameter group to associate.",
			Required:    true,
		},
		"token": {
			Type:        schema.TypeString,
			Description: "GitLab Runner token.",
			Required:    true,
			ForceNew:    true,
			Sensitive:   true,
			DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
				if len(oldValue) > 14 {
					oldValue = oldValue[5:14]
				}
				if len(newValue) > 14 {
					newValue = newValue[5:14]
				}
				return oldValue == newValue
			},
			ValidateDiagFunc: func(v any, p cty.Path) diag.Diagnostics {
				value := v.(string)
				var diags diag.Diagnostics
				if !strings.HasPrefix(value, "glrt-") {
					diag := diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "wrong value",
						Detail:   fmt.Sprintf("%q does not start with \"glrt-\".", value),
					}
					diags = append(diags, diag)
				}
				return diags
			},
		},
	}
}
