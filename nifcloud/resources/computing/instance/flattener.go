package instance

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
)

func flatten(d *schema.ResourceData, res *computing.DescribeInstancesOutput) error {
	if res == nil || len(res.ReservationSet) == 0 || len(res.ReservationSet[0].InstancesSet) == 0 {
		d.SetId("")
		return nil
	}

	instance := res.ReservationSet[0].InstancesSet[0]

	if nifcloud.ToString(instance.InstanceId) != d.Id() {
		return fmt.Errorf("unable to find instance within: %#v", res.ReservationSet)
	}

	if err := d.Set("accounting_type", instance.NextMonthAccountingType); err != nil {
		return err
	}

	if err := d.Set("availability_zone", instance.Placement.AvailabilityZone); err != nil {
		return err
	}

	if err := d.Set("description", instance.Description); err != nil {
		return err
	}

	if err := d.Set("image_id", instance.ImageId); err != nil {
		return err
	}

	if err := d.Set("instance_id", instance.InstanceId); err != nil {
		return err
	}

	if err := d.Set("instance_type", instance.InstanceType); err != nil {
		return err
	}

	if err := d.Set("key_name", instance.KeyName); err != nil {
		return err
	}

	var networkInterfaces []map[string]interface{}
	for _, n := range instance.NetworkInterfaceSet {
		ni := map[string]interface{}{
			"network_interface_id":            n.NetworkInterfaceId,
			"network_interface_attachment_id": n.Attachment.AttachmentId,
		}

		switch nifcloud.ToString(n.NiftyNetworkId) {
		case "net-COMMON_GLOBAL":
			if nifcloud.ToString(instance.IpType) == "elastic" {
				ni["ip_address"] = nifcloud.ToString(instance.IpAddress)
			}
			ni["network_id"] = nifcloud.ToString(n.NiftyNetworkId)
		case "net-COMMON_PRIVATE":
			if nifcloud.ToString(instance.NiftyPrivateIpType) == "elastic" {
				ni["ip_address"] = nifcloud.ToString(n.PrivateIpAddress)
			}
			ni["network_id"] = nifcloud.ToString(n.NiftyNetworkId)
		case "net-MULTI_IP_ADDRESS":
			ni["network_id"] = nifcloud.ToString(n.NiftyNetworkId)
			ni["multi_ip_address_group_id"] = nifcloud.ToString(instance.MultiIpAddressGroup.MultiIpAddressGroupId)
		default:
			var findElm map[string]interface{}
			for _, dn := range d.Get("network_interface").(*schema.Set).List() {
				elm := dn.(map[string]interface{})

				if elm["network_id"] != nil && n.NiftyNetworkId != nil && elm["network_id"] == nifcloud.ToString(n.NiftyNetworkId) {
					findElm = elm
					break
				}

				if elm["network_name"] != nil && n.NiftyNetworkName != nil && elm["network_name"] == nifcloud.ToString(n.NiftyNetworkName) {
					findElm = elm
					break
				}

				if elm["network_interface_id"] != nil && n.NetworkInterfaceId != nil && elm["network_interface_id"] == nifcloud.ToString(n.NetworkInterfaceId) {
					findElm = elm
					break
				}
			}

			if findElm != nil {
				if findElm["ip_address"] == "static" {
					ni["ip_address"] = "static"
				} else if findElm["ip_address"] != nil && findElm["ip_address"] != "" {
					ni["ip_address"] = nifcloud.ToString(n.PrivateIpAddress)
				}

				if findElm["network_id"] != nil && findElm["network_id"] != "" {
					ni["network_id"] = nifcloud.ToString(n.NiftyNetworkId)
				}

				if findElm["network_name"] != nil && findElm["network_name"] != "" {
					ni["network_name"] = nifcloud.ToString(n.NiftyNetworkName)
				}

			} else {
				ni["network_id"] = nifcloud.ToString(n.NiftyNetworkId)
				ni["ip_address"] = nifcloud.ToString(n.PrivateIpAddress)
			}
		}

		networkInterfaces = append(networkInterfaces, ni)
	}

	if err := d.Set("network_interface", networkInterfaces); err != nil {
		return err
	}

	if len(res.ReservationSet[0].GroupSet) > 0 {
		if err := d.Set("security_group", res.ReservationSet[0].GroupSet[0].GroupId); err != nil {
			return err
		}
	}

	if err := d.Set("instance_state", instance.InstanceState.Name); err != nil {
		return err
	}

	if err := d.Set("private_ip", instance.PrivateIpAddress); err != nil {
		return err
	}

	if err := d.Set("public_ip", instance.IpAddress); err != nil {
		return err
	}

	if err := d.Set("unique_id", instance.InstanceUniqueId); err != nil {
		return err
	}
	return nil
}

func flattenDisableAPITermination(d *schema.ResourceData, res *computing.DescribeInstanceAttributeOutput) error {
	if err := d.Set("disable_api_termination", res.DisableApiTermination.Value); err != nil {
		return err
	}
	return nil
}
