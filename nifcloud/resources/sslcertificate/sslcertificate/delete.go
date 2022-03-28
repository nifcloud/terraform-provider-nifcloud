package sslcertificate

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/smithy-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).Computing

	input := expandDeleteSSLCertificateInput(d)

	_, err := svc.DeleteSslCertificate(ctx, input)
	if err != nil {
		var awsErr smithy.APIError
		if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameterNotFound.SslCertificate" {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("failed deleting SSLCertificate: %s", err.Error()))
	}

	d.SetId("")

	return nil
}
