package record

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandListResourceRecordSets(d)
	svc := meta.(*client.Client).DNS
	req := svc.ListResourceRecordSetsRequest(input)

	res, err := req.Send(ctx)
	if err != nil {
		var awsErr awserr.Error
		if errors.As(err, &awsErr) && awsErr.Code() == "NoSuchHostedZone" {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("failed reading dns record: %s", err))
	}

	if d.IsNewResource() {
		for _, s := range res.ResourceRecordSets {
			for _, r := range s.ResourceRecords {
				if nifcloud.StringValue(r.Value) == d.Get("record").(string) {
					d.SetId(nifcloud.StringValue(s.SetIdentifier))
				}
			}
		}
	}

	if err := flatten(d, res); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
