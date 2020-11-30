package image

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).Computing

	req := svc.DescribeImagesRequest(&computing.DescribeImagesInput{
		ImageName: []string{d.Get("image_name").(string)},
		Owner:     []string{d.Get("owner").(string)},
	})

	res, err := req.Send(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed reading: %s", err))
	}

	images := res.ImagesSet[:]

	if len(images) < 1 {
		return diag.FromErr(fmt.Errorf("your query returned no results. Please change your search criteria and try again"))
	}

	if len(images) > 1 {
		return diag.FromErr(fmt.Errorf("your query returned more than one result. Please try a more specific search criteria"))
	}

	image := images[0]

	d.SetId(nifcloud.StringValue(image.ImageId))

	if err := d.Set("image_id", image.ImageId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("image_name", image.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("owner", image.ImageOwnerId); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
