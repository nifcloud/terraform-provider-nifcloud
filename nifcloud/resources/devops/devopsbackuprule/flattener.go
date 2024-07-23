package devopsbackuprule

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/devops"
)

func flatten(d *schema.ResourceData, res *devops.GetBackupRuleOutput) error {
	if res == nil || res.BackupRule == nil {
		d.SetId("")
		return nil
	}

	group := res.BackupRule

	if nifcloud.ToString(group.BackupRuleName) != d.Id() {
		return fmt.Errorf("unable to find the DevOps backup rule within: %#v", group)
	}

	if err := d.Set("name", group.BackupRuleName); err != nil {
		return err
	}

	if err := d.Set("instance_id", group.InstanceId); err != nil {
		return err
	}

	if err := d.Set("description", group.Description); err != nil {
		return err
	}

	if err := d.Set("backup_time", group.BackupTime); err != nil {
		return err
	}

	return nil
}
