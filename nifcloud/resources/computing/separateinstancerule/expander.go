package separateinstancerule

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
)

func expandNiftyCreateSeparateInstanceRuleInput(d *schema.ResourceData) *computing.NiftyCreateSeparateInstanceRuleInput {
	return &computing.NiftyCreateSeparateInstanceRuleInput{
		SeparateInstanceRuleName:        nifcloud.String(d.Get("name").(string)),
		SeparateInstanceRuleDescription: nifcloud.String(d.Get("description").(string)),
		InstanceId:                      expandInstanceIds(d.Get("instance_id").([]interface{})),
		InstanceUniqueId:                expandInstanceIds(d.Get("instance_unique_id").([]interface{})),
		Placement: &computing.RequestPlacementOfNiftyCreateSeparateInstanceRule{
			AvailabilityZone: nifcloud.String(d.Get("availability_zone").(string)),
		},
	}
}

func expandInstanceIds(raw []interface{}) []string {
	if len(raw) == 0 {
		return nil
	}

	ids := make([]string, len(raw))
	for i, l := range raw {
		ids[i] = l.(string)
	}

	return ids
}

func expandNiftyDescribeSeparateInstanceRulesInput(d *schema.ResourceData) *computing.NiftyDescribeSeparateInstanceRulesInput {
	return &computing.NiftyDescribeSeparateInstanceRulesInput{
		SeparateInstanceRuleName: []string{d.Id()},
	}
}

func expandNiftyDeleteSeparateInstanceRuleInput(d *schema.ResourceData) *computing.NiftyDeleteSeparateInstanceRuleInput {
	return &computing.NiftyDeleteSeparateInstanceRuleInput{
		SeparateInstanceRuleName: nifcloud.String(d.Id()),
	}
}

func expandNiftyUpdateSeparateInstanceRuleInputForName(d *schema.ResourceData) *computing.NiftyUpdateSeparateInstanceRuleInput {
	before, after := d.GetChange("name")
	return &computing.NiftyUpdateSeparateInstanceRuleInput{
		SeparateInstanceRuleName:       nifcloud.String(before.(string)),
		SeparateInstanceRuleNameUpdate: nifcloud.String(after.(string)),
	}
}

func expandNiftyUpdateSeparateInstanceRuleInputForDescription(d *schema.ResourceData) *computing.NiftyUpdateSeparateInstanceRuleInput {
	return &computing.NiftyUpdateSeparateInstanceRuleInput{
		SeparateInstanceRuleName:              nifcloud.String(d.Id()),
		SeparateInstanceRuleDescriptionUpdate: nifcloud.String(d.Get("description").(string)),
	}
}

func expandNiftyRegisterInstancesWithSeparateInstanceRuleInstanceIDInput(
	d *schema.ResourceData,
	instanceIDs []string,
) *computing.NiftyRegisterInstancesWithSeparateInstanceRuleInput {
	return &computing.NiftyRegisterInstancesWithSeparateInstanceRuleInput{
		SeparateInstanceRuleName: nifcloud.String(d.Get("name").(string)),
		InstanceId:               instanceIDs,
	}
}

func expandNiftyDeregisterInstancesFromSeparateInstanceRuleInstanceIDInput(
	d *schema.ResourceData,
	instanceIDs []string,
) *computing.NiftyDeregisterInstancesFromSeparateInstanceRuleInput {
	return &computing.NiftyDeregisterInstancesFromSeparateInstanceRuleInput{
		SeparateInstanceRuleName: nifcloud.String(d.Get("name").(string)),
		InstanceId:               instanceIDs,
	}
}

func expandNiftyRegisterInstancesWithSeparateInstanceRuleInstanceUniqueIDInput(
	d *schema.ResourceData,
	instanceUniqueIDs []string,
) *computing.NiftyRegisterInstancesWithSeparateInstanceRuleInput {
	return &computing.NiftyRegisterInstancesWithSeparateInstanceRuleInput{
		SeparateInstanceRuleName: nifcloud.String(d.Get("name").(string)),
		InstanceUniqueId:         instanceUniqueIDs,
	}
}

func expandNiftyDeregisterInstancesFromSeparateInstanceRuleInstanceUniqueIDInput(
	d *schema.ResourceData,
	instanceUniqueIDs []string,
) *computing.NiftyDeregisterInstancesFromSeparateInstanceRuleInput {
	return &computing.NiftyDeregisterInstancesFromSeparateInstanceRuleInput{
		SeparateInstanceRuleName: nifcloud.String(d.Get("name").(string)),
		InstanceUniqueId:         instanceUniqueIDs,
	}
}
