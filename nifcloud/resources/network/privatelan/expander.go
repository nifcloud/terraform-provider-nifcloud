package privatelan

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing/types"
)

func expandNiftyCreatePrivateLanInput(d *schema.ResourceData) *computing.NiftyCreatePrivateLanInput {
	at := d.Get("accounting_type").(string)
	return &computing.NiftyCreatePrivateLanInput{
		AccountingType:   types.AccountingTypeOfNiftyCreatePrivateLanRequest(at),
		AvailabilityZone: nifcloud.String(d.Get("availability_zone").(string)),
		CidrBlock:        nifcloud.String(d.Get("cidr_block").(string)),
		Description:      nifcloud.String(d.Get("description").(string)),
		PrivateLanName:   nifcloud.String(d.Get("private_lan_name").(string)),
	}
}

func expandNiftyDescribePrivateLansInput(d *schema.ResourceData) *computing.NiftyDescribePrivateLansInput {
	return &computing.NiftyDescribePrivateLansInput{
		NetworkId: []string{d.Id()},
	}
}

func expandNiftyModifyPrivateLanAttributeInputForPrivateLanName(d *schema.ResourceData) *computing.NiftyModifyPrivateLanAttributeInput {
	return &computing.NiftyModifyPrivateLanAttributeInput{
		NetworkId: nifcloud.String(d.Get("network_id").(string)),
		Attribute: "privateLanName",
		Value:     nifcloud.String(d.Get("private_lan_name").(string)),
	}
}

func expandNiftyModifyPrivateLanAttributeInputForCidrBlock(d *schema.ResourceData) *computing.NiftyModifyPrivateLanAttributeInput {
	return &computing.NiftyModifyPrivateLanAttributeInput{
		NetworkId: nifcloud.String(d.Get("network_id").(string)),
		Attribute: "cidrBlock",
		Value:     nifcloud.String(d.Get("cidr_block").(string)),
	}
}

func expandNiftyModifyPrivateLanAttributeInputForAccountingType(d *schema.ResourceData) *computing.NiftyModifyPrivateLanAttributeInput {
	return &computing.NiftyModifyPrivateLanAttributeInput{
		NetworkId: nifcloud.String(d.Get("network_id").(string)),
		Attribute: "accountingType",
		Value:     nifcloud.String(d.Get("accounting_type").(string)),
	}
}

func expandNiftyModifyPrivateLanAttributeInputForDescription(d *schema.ResourceData) *computing.NiftyModifyPrivateLanAttributeInput {
	return &computing.NiftyModifyPrivateLanAttributeInput{
		NetworkId: nifcloud.String(d.Get("network_id").(string)),
		Attribute: "description",
		Value:     nifcloud.String(d.Get("description").(string)),
	}
}
