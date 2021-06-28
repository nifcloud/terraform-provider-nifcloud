package separateinstancerule

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
)

func flatten(d *schema.ResourceData, res *computing.NiftyDescribeSeparateInstanceRulesResponse) error {
	if res == nil || len(res.SeparateInstanceRulesInfo) == 0 {
		d.SetId("")
		return nil
	}

	separateInstanceRule := res.SeparateInstanceRulesInfo[0]

	if err := d.Set("name", separateInstanceRule.SeparateInstanceRuleName); err != nil {
		return err
	}

	if err := d.Set("description", separateInstanceRule.SeparateInstanceRuleDescription); err != nil {
		return err
	}

	if err := d.Set("availability_zone", separateInstanceRule.AvailabilityZone); err != nil {
		return err
	}

	if _, ok := d.GetOk("instance_id"); ok {
		if err := d.Set("instance_id", flattenInstanceID(separateInstanceRule.InstancesSet)); err != nil {
			return err
		}
	}

	if _, ok := d.GetOk("instance_unique_id"); ok {
		if err := d.Set("instance_unique_id", flattenInstanceUniqueID(separateInstanceRule.InstancesSet)); err != nil {
			return err
		}
	}

	return nil
}

func flattenInstanceID(instancesSet []computing.InstancesSetOfNiftyDescribeSeparateInstanceRules) []string {
	res := make([]string, len(instancesSet))

	for i, InstancesSetOfNiftyDescribeSeparateInstanceRules := range instancesSet {
		res[i] = nifcloud.StringValue(InstancesSetOfNiftyDescribeSeparateInstanceRules.InstanceId)
	}
	return res
}

func flattenInstanceUniqueID(instancesSet []computing.InstancesSetOfNiftyDescribeSeparateInstanceRules) []string {
	res := make([]string, len(instancesSet))

	for i, InstancesSetOfNiftyDescribeSeparateInstanceRules := range instancesSet {
		res[i] = nifcloud.StringValue(InstancesSetOfNiftyDescribeSeparateInstanceRules.InstanceUniqueId)
	}
	return res
}
