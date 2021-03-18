package sslcertificate

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
)

type describeResponses struct {
	describeSSLCertificatesResponse       *computing.DescribeSslCertificatesResponse
	downloadSSLCertificateResponseForCert *computing.DownloadSslCertificateResponse
	downloadSSLCertificateResponseForKey  *computing.DownloadSslCertificateResponse
	downloadSSLCertificateResponseForCA   *computing.DownloadSslCertificateResponse
}

func flatten(d *schema.ResourceData, res *describeResponses) error {
	if res == nil {
		return fmt.Errorf("describe ssl certificate result is empty")
	}

	if res.describeSSLCertificatesResponse == nil || len(res.describeSSLCertificatesResponse.CertsSet) == 0 {
		d.SetId("")
		return nil
	}

	sslCertificate := res.describeSSLCertificatesResponse.CertsSet[0]

	if nifcloud.StringValue(sslCertificate.FqdnId) != d.Id() {
		return fmt.Errorf("unable to find ssl certificate with in: %#v", res.describeSSLCertificatesResponse.CertsSet)
	}

	if err := d.Set("fqdn_id", sslCertificate.FqdnId); err != nil {
		return err
	}

	if err := d.Set("fqdn", sslCertificate.Fqdn); err != nil {
		return err
	}

	if err := d.Set("description", sslCertificate.Description); err != nil {
		return err
	}

	if res.downloadSSLCertificateResponseForCert == nil {
		return fmt.Errorf("download certificate result is empty")
	}
	if err := d.Set(
		"certificate", nifcloud.StringValue(res.downloadSSLCertificateResponseForCert.FileData),
	); err != nil {
		return err
	}

	if res.downloadSSLCertificateResponseForKey == nil {
		return fmt.Errorf("download private key result is empty")
	}
	if err := d.Set(
		"key", nifcloud.StringValue(res.downloadSSLCertificateResponseForKey.FileData),
	); err != nil {
		return err
	}

	if res.downloadSSLCertificateResponseForCA != nil {
		if err := d.Set(
			"ca", nifcloud.StringValue(res.downloadSSLCertificateResponseForCA.FileData),
		); err != nil {
			return err
		}
	}

	return nil
}
