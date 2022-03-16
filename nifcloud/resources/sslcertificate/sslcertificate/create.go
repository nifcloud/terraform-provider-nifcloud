package sslcertificate

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := expandUploadSSLCertificateInput(d)

	svc := meta.(*client.Client).Computing

	res, err := svc.UploadSslCertificate(ctx, input)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed uploading SSLCertificate: %s", err.Error()))
	}

	d.SetId(nifcloud.ToString(res.FqdnId))

	_, err = svc.ModifySslCertificateAttribute(ctx, expandModifySSLCertificateAttributeInput(d))
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed updating SSLCertificate: %s", err.Error()))
	}

	return read(ctx, d, meta)
}
