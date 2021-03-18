package sslcertificate

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/stretchr/testify/assert"
)

func TestExpandUploadSSLCertificateInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"certificate": "test_certificate",
		"key":         "test_key",
		"ca":          "test_ca",
	})

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.UploadSslCertificateInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.UploadSslCertificateInput{
				Certificate:          nifcloud.String("test_certificate"),
				Key:                  nifcloud.String("test_key"),
				CertificateAuthority: nifcloud.String("test_ca"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandUploadSSLCertificateInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDescribeSslCertificatesInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"fqdn_id": "test_fqdn_id",
	})
	rd.SetId("test_fqdn_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.DescribeSslCertificatesInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.DescribeSslCertificatesInput{
				FqdnId: []string{"test_fqdn_id"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDescribeSSLCertificatesInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDownloadSSLCertificateInputForKey(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"fqdn_id": "test_fqdn_id",
	})
	rd.SetId("test_fqdn_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.DownloadSslCertificateInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.DownloadSslCertificateInput{
				FqdnId:   nifcloud.String("test_fqdn_id"),
				FileType: computing.FileTypeOfDownloadSslCertificateRequest1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDownloadSSLCertificateInputForKey(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDownloadSSLCertificateInputForCA(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"fqdn_id": "test_fqdn_id",
	})
	rd.SetId("test_fqdn_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.DownloadSslCertificateInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.DownloadSslCertificateInput{
				FqdnId:   nifcloud.String("test_fqdn_id"),
				FileType: computing.FileTypeOfDownloadSslCertificateRequest2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDownloadSSLCertificateInputForCA(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDownloadSSLCertificateInputForCert(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"fqdn_id": "test_fqdn_id",
	})
	rd.SetId("test_fqdn_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.DownloadSslCertificateInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.DownloadSslCertificateInput{
				FqdnId:   nifcloud.String("test_fqdn_id"),
				FileType: computing.FileTypeOfDownloadSslCertificateRequest3,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDownloadSSLCertificateInputForCert(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandModifySSLCertificateAttributeInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"fqdn_id":     "test_fqdn_id",
		"description": "test_description",
	})
	rd.SetId("test_fqdn_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.ModifySslCertificateAttributeInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.ModifySslCertificateAttributeInput{
				FqdnId: nifcloud.String("test_fqdn_id"),
				Description: &computing.RequestDescription{
					Value: nifcloud.String("test_description"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandModifySSLCertificateAttributeInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandDeleteSSLCertificateInput(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"fqdn_id": "test_fqdn_id",
	})
	rd.SetId("test_fqdn_id")

	tests := []struct {
		name string
		args *schema.ResourceData
		want *computing.DeleteSslCertificateInput
	}{
		{
			name: "expands the resource data",
			args: rd,
			want: &computing.DeleteSslCertificateInput{
				FqdnId: nifcloud.String("test_fqdn_id"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDeleteSSLCertificateInput(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}
