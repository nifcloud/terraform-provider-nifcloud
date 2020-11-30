package sslcertificate

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
)

func expandUploadSSLCertificateInput(d *schema.ResourceData) *computing.UploadSslCertificateInput {
	return &computing.UploadSslCertificateInput{
		Certificate:          nifcloud.String(d.Get("certificate").(string)),
		Key:                  nifcloud.String(d.Get("key").(string)),
		CertificateAuthority: nifcloud.String(d.Get("ca").(string)),
	}
}

func expandDescribeSSLCertificatesInput(d *schema.ResourceData) *computing.DescribeSslCertificatesInput {
	return &computing.DescribeSslCertificatesInput{
		FqdnId: []string{d.Id()},
	}
}

func expandDownloadSSLCertificateInputForKey(d *schema.ResourceData) *computing.DownloadSslCertificateInput {
	return &computing.DownloadSslCertificateInput{
		FqdnId:   nifcloud.String(d.Id()),
		FileType: computing.FileTypeOfDownloadSslCertificateRequest1,
	}
}

func expandDownloadSSLCertificateInputForCA(d *schema.ResourceData) *computing.DownloadSslCertificateInput {
	return &computing.DownloadSslCertificateInput{
		FqdnId:   nifcloud.String(d.Id()),
		FileType: computing.FileTypeOfDownloadSslCertificateRequest2,
	}
}

func expandDownloadSSLCertificateInputForCert(d *schema.ResourceData) *computing.DownloadSslCertificateInput {
	return &computing.DownloadSslCertificateInput{
		FqdnId:   nifcloud.String(d.Id()),
		FileType: computing.FileTypeOfDownloadSslCertificateRequest3,
	}
}

func expandModifySSLCertificateAttributeInput(d *schema.ResourceData) *computing.ModifySslCertificateAttributeInput {
	return &computing.ModifySslCertificateAttributeInput{
		FqdnId: nifcloud.String(d.Id()),
		Description: &computing.RequestDescription{
			Value: nifcloud.String(d.Get("description").(string)),
		},
	}
}

func expandDeleteSSLCertificateInput(d *schema.ResourceData) *computing.DeleteSslCertificateInput {
	return &computing.DeleteSslCertificateInput{
		FqdnId: nifcloud.String(d.Id()),
	}
}
