package record

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/smithy-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandListResourceRecordSets(d)
	svc := meta.(*client.Client).DNS

	res, err := svc.ListResourceRecordSets(ctx, input)
	if err != nil {
		var awsErr smithy.APIError
		if errors.As(err, &awsErr) && awsErr.ErrorCode() == "NoSuchHostedZone" {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("failed reading dns record: %s", err))
	}

	if d.IsNewResource() {
		for _, s := range res.ResourceRecordSets {
			for _, r := range s.ResourceRecords {
				if nifcloud.ToString(r.Value) == d.Get("record").(string) {
					d.SetId(nifcloud.ToString(s.SetIdentifier))
				}
			}
		}
	}

	if err := flatten(d, res); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
