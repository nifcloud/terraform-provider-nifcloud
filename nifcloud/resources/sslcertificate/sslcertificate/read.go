package sslcertificate

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/smithy-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
	"golang.org/x/sync/errgroup"
)

func read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).Computing

	res := describeOutputs{}

	eg, errCtx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		var err error
		res.describeSSLCertificatesOutput, err = svc.DescribeSslCertificates(errCtx, expandDescribeSSLCertificatesInput(d))
		if err != nil {
			return fmt.Errorf("failed reading SSLCertificate: %s", err.Error())
		}
		return nil
	})

	eg.Go(func() error {
		var err error
		res.downloadSSLCertificateOutputForCert, err = svc.DownloadSslCertificate(errCtx, expandDownloadSSLCertificateInputForCert(d))
		if err != nil {
			return checkNotFoundError(err)
		}
		return nil
	})

	eg.Go(func() error {
		var err error
		res.downloadSSLCertificateOutputForKey, err = svc.DownloadSslCertificate(errCtx, expandDownloadSSLCertificateInputForKey(d))
		if err != nil {
			return checkNotFoundError(err)
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return diag.FromErr(err)
	}

	if res.describeSSLCertificatesOutput != nil &&
		len(res.describeSSLCertificatesOutput.CertsSet) == 1 &&
		nifcloud.ToBool(res.describeSSLCertificatesOutput.CertsSet[0].CaState) {
		var err error
		res.downloadSSLCertificateOutputForCA, err = svc.DownloadSslCertificate(ctx, expandDownloadSSLCertificateInputForCA(d))
		if err != nil {
			return diag.FromErr(checkNotFoundError(err))
		}
	}

	if err := flatten(d, &res); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func checkNotFoundError(err error) error {
	var awserr smithy.APIError
	if errors.As(err, &awserr) && awserr.ErrorCode() == "Client.InvalidParameterNotFound.SslCertificate" {
		return nil
	}
	return fmt.Errorf("failed downloading certificate: %s", err.Error())
}
