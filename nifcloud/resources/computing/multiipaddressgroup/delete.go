package multiipaddressgroup

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).Computing

	input := expandDeleteMultiIPAddressGroupInput(d)
	if _, err := svc.DeleteMultiIpAddressGroup(ctx, input); err != nil {
		return diag.FromErr(fmt.Errorf("failed deleting multi IP address group: %s", err))
	}

	waiter := computing.NewMultiIpAddressGroupDeletedWaiter(svc)
	if err := waiter.Wait(ctx, expandDescribeMultiIPAddressGroupsInput(d), 5*time.Minute); err != nil {
		return diag.Errorf("failed waiting multi IP address group to be deleted: %s", err)
	}

	d.SetId("")

	return nil
}
