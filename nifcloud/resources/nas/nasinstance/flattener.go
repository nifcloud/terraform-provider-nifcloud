package nasinstance

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/nas"
)

func flatten(d *schema.ResourceData, res *nas.DescribeNASInstancesOutput) error {
	if res == nil || len(res.NASInstances) == 0 {
		d.SetId("")
		return nil
	}

	nasInstance := res.NASInstances[0]

	if nifcloud.ToString(nasInstance.NASInstanceIdentifier) != d.Id() {
		return fmt.Errorf("unable to find NAS instance within: %#v", res.NASInstances)
	}

	if err := d.Set("identifier", nasInstance.NASInstanceIdentifier); err != nil {
		return err
	}

	if err := d.Set("description", nasInstance.NASInstanceDescription); err != nil {
		return err
	}

	if err := d.Set("allocated_storage", nasInstance.AllocatedStorage); err != nil {
		return err
	}

	if err := d.Set("availability_zone", nasInstance.AvailabilityZone); err != nil {
		return err
	}

	if err := d.Set("nas_security_group_name", nasInstance.NASSecurityGroups[0].NASSecurityGroupName); err != nil {
		return err
	}

	if err := d.Set("private_ip_address", nasInstance.Endpoint.PrivateAddress); err != nil {
		return err
	}

	if err := d.Set("public_ip_address", nasInstance.Endpoint.Address); err != nil {
		return err
	}

	if err := d.Set("protocol", nasInstance.Protocol); err != nil {
		return err
	}

	if err := d.Set("master_username", nasInstance.MasterUsername); err != nil {
		return err
	}

	if err := d.Set("type", nasInstance.NASInstanceType); err != nil {
		return err
	}

	if err := d.Set("no_root_squash", nasInstance.NoRootSquash); err != nil {
		return err
	}

	if err := d.Set("network_id", nasInstance.NetworkId); err != nil {
		return err
	}

	return nil
}
