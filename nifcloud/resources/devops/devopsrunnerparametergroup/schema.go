package devopsrunnerparametergroup

import (
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/validator"
)

const description = "Provides a DevOps Runner parameter group resource."

// New returns the nifcloud_devops_runner_parameter_group resource schema.
func New() *schema.Resource {
	return &schema.Resource{
		Description: description,
		Schema:      newSchema(),

		CreateContext: createRunnerParameterGroup,
		ReadContext:   readRunnerParameterGroup,
		UpdateContext: updateRunnerParameterGroup,
		DeleteContext: deleteRunnerParameterGroup,

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
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the DevOps Runner parameter group.",
			Required:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 63),
				validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z\-]+$`), "Enter a name within 1-63 alphanumeric lowercase characters and hyphens. Hyphens cannot be used at the beginning or end."),
			),
		},
		"description": {
			Type:             schema.TypeString,
			Description:      "Description of the DevOps Runner parameter group.",
			Optional:         true,
			ValidateDiagFunc: validator.StringRuneCountBetween(0, 255),
		},
		"docker_disable_cache": {
			Type:        schema.TypeBool,
			Description: "The Docker executor has two levels of caching: a global one (like any other executor) and a local cache based on Docker volumes. This configuration flag acts only on the local one which disables the use of automatically created (not mapped to a host directory) cache volumes. In other words, it only prevents creating a container that holds temporary files of builds, it does not disable the cache if the runner is configured in distributed cache mode.",
			Optional:    true,
			Computed:    true,
		},
		"docker_disable_entrypoint_overwrite": {
			Type:        schema.TypeBool,
			Description: "Disable the image entrypoint overwriting.",
			Optional:    true,
			Computed:    true,
		},
		"docker_extra_host": {
			Type:        schema.TypeSet,
			Description: "Hosts that should be defined in container environment.",
			Optional:    true,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"host_name": {
						Type:        schema.TypeString,
						Description: "Host name.",
						Required:    true,
					},
					"ip_address": {
						Type:        schema.TypeString,
						Description: "IPv4 address.",
						Required:    true,
					},
				},
			},
		},
		"docker_image": {
			Type:        schema.TypeString,
			Description: "The image to run jobs with.",
			Optional:    true,
			Computed:    true,
		},
		"docker_oom_kill_disable": {
			Type:        schema.TypeBool,
			Description: "If an out-of-memory (OOM) error occurs, do not kill processes in a container.",
			Optional:    true,
			Computed:    true,
		},
		"docker_privileged": {
			Type:        schema.TypeBool,
			Description: "Run all containers with the privileged flag enabled.",
			Optional:    true,
			Computed:    true,
		},
		"docker_shm_size": {
			Type:        schema.TypeInt,
			Description: "Shared memory size for images (in bytes).",
			Optional:    true,
			Computed:    true,
		},
		"docker_tls_verify": {
			Type:        schema.TypeBool,
			Description: "Enable or disable TLS verification of connections to Docker daemon. Disabled by default.",
			Optional:    true,
			Computed:    true,
		},
		"docker_volume": {
			Type:        schema.TypeSet,
			Description: "Additional volumes that should be mounted. Same syntax as the Docker -v flag.",
			Optional:    true,
			Computed:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}
