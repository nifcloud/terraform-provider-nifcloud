package sslcertificate

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).Computing

	input := expandDeleteSSLCertificateInput(d)

	req := svc.DeleteSslCertificateRequest(input)

	_, err := req.Send(ctx)
	if err != nil {
		var awsErr awserr.Error
		if errors.As(err, &awsErr) && awsErr.Code() == "Client.InvalidParameterNotFound.SslCertificate" {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("failed deleting SSLCertificate: %s", err.Error()))
	}

	d.SetId("")

	return nil
}
