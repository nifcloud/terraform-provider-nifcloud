package dbinstance

import (
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const description = "Provides a rdb instance resource."

// New returns the nifcloud_db_instance resource schema.
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
			Create:  schema.DefaultTimeout(120 * time.Minute),
			Update:  schema.DefaultTimeout(120 * time.Minute),
			Delete:  schema.DefaultTimeout(60 * time.Minute),
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
		"instance_class": {
			Type:        schema.TypeString,
			Description: "The instance type of the DB instance.",
			Required:    true,
		},
		"db_name": {
			Type:        schema.TypeString,
			Description: "The name of the database to create when the DB instance is created. If this parameter is not specified, no database is created.",
			Optional:    true,
			Computed:    true,
			ForceNew:    true,
		},
		"username": {
			Type:        schema.TypeString,
			Description: "Username for the master DB user.",
			Optional:    true,
			Computed:    true,
			ForceNew:    true,
		},
		"password": {
			Type:        schema.TypeString,
			Description: "Password for the master DB user.",
			Optional:    true,
			Sensitive:   true,
		},
		"engine": {
			Type:        schema.TypeString,
			Description: "The database engine. `MySQL` or `postgres` or `MariaDB`",
			Optional:    true,
			Computed:    true,
			ForceNew:    true,
			StateFunc: func(v interface{}) string {
				value := v.(string)
				return strings.ToLower(value)
			},
		},
		"engine_version": {
			Type:        schema.TypeString,
			Description: "The database engine version.",
			Optional:    true,
			Computed:    true,
			ForceNew:    true,
		},
		"ca_cert_identifier": {
			Type:        schema.TypeString,
			Description: "The identifier of the CA certificate for the DB instance.",
			Optional:    true,
			Computed:    true,
		},
		"allocated_storage": {
			Type:        schema.TypeInt,
			Description: "The allocated storage in gibibytes.",
			Optional:    true,
			Computed:    true,
		},
		"storage_type": {
			Type:        schema.TypeInt,
			Description: "One of `0` (HDD), or `1` (Flash drive). The default is `0`",
			Optional:    true,
			Default:     0,
			ForceNew:    true,
		},
		"identifier": {
			Type:        schema.TypeString,
			Description: "The name of the DB instance.",
			Required:    true,
		},

		"availability_zone": {
			Type:        schema.TypeString,
			Description: "The AZ for the DB instance.",
			Optional:    true,
			Computed:    true,
			ForceNew:    true,
		},
		"backup_retention_period": {
			Type:        schema.TypeInt,
			Description: "The days to retain backups for. If `0` automatic backup will be off",
			Optional:    true,
			Computed:    true,
		},
		"binlog_retention_period": {
			Type:        schema.TypeInt,
			Description: "The days to retain binlog for. Be sure to specify `custom_binlog_retention_period = true` as a set",
			Optional:    true,
			Computed:    true,
		},
		"custom_binlog_retention_period": {
			Type:        schema.TypeBool,
			Description: "The flag of set binary log retention period. Only MySQL can be specified",
			Optional:    true,
		},
		"backup_window": {
			Type:        schema.TypeString,
			Description: "The daily time range (in UTC) during which automated backups are created if they are enabled. Example: `09:46-10:16`",
			Optional:    true,
			Computed:    true,
		},
		"maintenance_window": {
			Type:        schema.TypeString,
			Description: "The weekly time range (in UTC) the instance maintenance window. Example: `Sun:05:00-Sun:06:00`",
			Optional:    true,
			Computed:    true,
		},
		"multi_az": {
			Description: "If the DB instance is multi AZ enabled.",
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
		},
		"multi_az_type": {
			Description: "The type of multi AZ. (0: Data priority, 1: Performance priority) default `0`",
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     0,
		},

		"port": {
			Type:        schema.TypeInt,
			Description: "The database port.",
			Optional:    true,
			ForceNew:    true,
			Computed:    true,
		},
		"publicly_accessible": {
			Type:        schema.TypeBool,
			Description: "Bool to control if instance is publicly accessible. Default is `true`",
			Optional:    true,
			Default:     true,
			ForceNew:    true,
		},
		"db_security_group_name": {
			Type:        schema.TypeString,
			Description: "The security group name to associate with; which can be managed using the nifcloud_db_security_group resource.",
			Optional:    true,
			Computed:    true,
		},
		"final_snapshot_identifier": {
			Type:        schema.TypeString,
			Description: "The name of your final DB snapshot when this DB instance is deleted. Must be provided if `skip_final_snapshot` is set to false.",
			Optional:    true,
		},
		"restore_to_point_in_time": {
			Type:        schema.TypeList,
			Description: "A configuration block for restoring a DB instance to an arbitrary point in time.",
			Optional:    true,
			MaxItems:    1,
			ForceNew:    true,
			ConflictsWith: []string{
				"snapshot_identifier",
				"replicate_source_db",
			},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"restore_time": {
						Type:          schema.TypeString,
						Description:   "The date and time to restore from. Value must be a time in Universal Coordinated Time (UTC) format.",
						Optional:      true,
						ConflictsWith: []string{"restore_to_point_in_time.0.use_latest_restorable_time"},
					},
					"source_db_instance_identifier": {
						Type:        schema.TypeString,
						Description: "The identifier of the source DB instance from which to restore.",
						Required:    true,
					},
					"use_latest_restorable_time": {
						Type:          schema.TypeBool,
						Description:   "A boolean value that indicates whether the DB instance is restored from the latest backup time. Defaults to `false`",
						Optional:      true,
						Default:       false,
						ConflictsWith: []string{"restore_to_point_in_time.0.restore_time"},
					},
				},
			},
		},
		"skip_final_snapshot": {
			Type:        schema.TypeBool,
			Description: "Determines whether a final DB snapshot is created before the DB instance is deleted. Defaults to `true` no DBSnapshot is created",
			Optional:    true,
			Default:     true,
		},
		"parameter_group_name": {
			Type:        schema.TypeString,
			Description: "Name of the DB parameter group to associate; which can be managed using the nifcloud_db_parameter_group resource.",
			Optional:    true,
			Computed:    true,
		},
		"address": {
			Type:        schema.TypeString,
			Description: "The hostname of the DB instance.",
			Computed:    true,
		},
		"replicate_source_db": {
			Type:        schema.TypeString,
			Description: "Specifies that this resource is a Replicate database, and to use this value as the source database.",
			Optional:    true,
			ForceNew:    true,
		},
		"snapshot_identifier": {
			Type:        schema.TypeString,
			Description: "Specifies whether or not to create this database from a snapshot.",
			Optional:    true,
			ForceNew:    true,
		},
		"network_id": {
			Type:        schema.TypeString,
			Description: "The id of private lan.",
			Optional:    true,
			Computed:    true,
		},
		"virtual_private_address": {
			Type:        schema.TypeString,
			Description: "Private IP address for virtual load balancer.",
			Optional:    true,
		},
		"master_private_address": {
			Type:        schema.TypeString,
			Description: "Private IP address for master DB.",
			Optional:    true,
		},
		"slave_private_address": {
			Type:        schema.TypeString,
			Description: "Private IP address for master DB.",
			Optional:    true,
		},
		"read_replica_private_address": {
			Type:        schema.TypeString,
			Description: "Private IP address for read replica.",
			Optional:    true,
		},
		"read_replica_identifier": {
			Type:        schema.TypeString,
			Description: "The DB instance name for read replica.",
			Optional:    true,
		},
		"apply_immediately": {
			Type:        schema.TypeBool,
			Description: "Specifies whether any database modifications are applied immediately, or during the next maintenance window. Default is `false`",
			Optional:    true,
			Default:     false,
		},
	}
}
