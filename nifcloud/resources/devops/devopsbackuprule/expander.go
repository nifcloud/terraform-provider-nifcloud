package devopsbackuprule

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/devops"
)

func expandCreateBackupRuleInput(d *schema.ResourceData) *devops.CreateBackupRuleInput {
	return &devops.CreateBackupRuleInput{
		BackupRuleName: nifcloud.String(d.Get("name").(string)),
		InstanceId:     nifcloud.String(d.Get("instance_id").(string)),
		Description:    nifcloud.String(d.Get("description").(string)),
	}
}

func expandUpdateBackupRuleInput(d *schema.ResourceData) *devops.UpdateBackupRuleInput {
	return &devops.UpdateBackupRuleInput{
		BackupRuleName:        nifcloud.String(d.Id()),
		ChangedBackupRuleName: nifcloud.String(d.Get("name").(string)),
		Description:           nifcloud.String(d.Get("description").(string)),
	}
}

func expandGetBackupRuleInput(d *schema.ResourceData) *devops.GetBackupRuleInput {
	return &devops.GetBackupRuleInput{
		BackupRuleName: nifcloud.String(d.Id()),
	}
}

func expandDeleteBackupRuleInput(d *schema.ResourceData) *devops.DeleteBackupRuleInput {
	return &devops.DeleteBackupRuleInput{
		BackupRuleName: nifcloud.String(d.Id()),
	}
}
