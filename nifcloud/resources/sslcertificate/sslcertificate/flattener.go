package sslcertificate

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
)

type describeOutputs struct {
	describeSSLCertificatesOutput       *computing.DescribeSslCertificatesOutput
	downloadSSLCertificateOutputForCert *computing.DownloadSslCertificateOutput
	downloadSSLCertificateOutputForKey  *computing.DownloadSslCertificateOutput
	downloadSSLCertificateOutputForCA   *computing.DownloadSslCertificateOutput
}

func flatten(d *schema.ResourceData, res *describeOutputs) error {
	if res == nil {
		return fmt.Errorf("describe ssl certificate result is empty")
	}

	if res.describeSSLCertificatesOutput == nil || len(res.describeSSLCertificatesOutput.CertsSet) == 0 {
		d.SetId("")
		return nil
	}

	sslCertificate := res.describeSSLCertificatesOutput.CertsSet[0]

	if nifcloud.ToString(sslCertificate.FqdnId) != d.Id() {
		return fmt.Errorf("unable to find ssl certificate with in: %#v", res.describeSSLCertificatesOutput.CertsSet)
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

	if res.downloadSSLCertificateOutputForCert == nil {
		return fmt.Errorf("download certificate result is empty")
	}
	if err := d.Set(
		"certificate", nifcloud.ToString(res.downloadSSLCertificateOutputForCert.FileData),
	); err != nil {
		return err
	}

	if res.downloadSSLCertificateOutputForKey == nil {
		return fmt.Errorf("download private key result is empty")
	}
	if err := d.Set(
		"key", nifcloud.ToString(res.downloadSSLCertificateOutputForKey.FileData),
	); err != nil {
		return err
	}

	if res.downloadSSLCertificateOutputForCA != nil {
		if err := d.Set(
			"ca", nifcloud.ToString(res.downloadSSLCertificateOutputForCA.FileData),
		); err != nil {
			return err
		}
	}

	return nil
}
