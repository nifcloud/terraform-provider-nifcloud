package sslcertificate

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing/types"
	"github.com/stretchr/testify/assert"
)

func TestFlatten(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"fqdn_id":     "test_fqdn_id",
		"fqdn":        "test_fqdn",
		"description": "test_description",
		"certificate": "test_certificate",
		"key":         "test_key",
		"ca":          "test_ca",
	})
	rd.SetId("test_fqdn_id")

	rdWithoutCA := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"fqdn_id":     "test_fqdn_id",
		"fqdn":        "test_fqdn",
		"description": "test_description",
		"certificate": "test_certificate",
		"key":         "test_key",
	})
	rdWithoutCA.SetId("test_fqdn_id")

	wantNotFoundRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})

	type args struct {
		res *describeOutputs
		d   *schema.ResourceData
	}
	tests := []struct {
		name string
		args args
		want *schema.ResourceData
	}{
		{
			name: "flattens the response",
			args: args{
				d: rd,
				res: &describeOutputs{
					describeSSLCertificatesOutput: &computing.DescribeSslCertificatesOutput{
						CertsSet: []types.CertsSet{
							{
								FqdnId:      nifcloud.String("test_fqdn_id"),
								Fqdn:        nifcloud.String("test_fqdn"),
								Description: nifcloud.String("test_description"),
							},
						},
					},
					downloadSSLCertificateOutputForCert: &computing.DownloadSslCertificateOutput{
						FileData: nifcloud.String("test_certificate"),
					},
					downloadSSLCertificateOutputForKey: &computing.DownloadSslCertificateOutput{
						FileData: nifcloud.String("test_key"),
					},
					downloadSSLCertificateOutputForCA: &computing.DownloadSslCertificateOutput{
						FileData: nifcloud.String("test_ca"),
					},
				},
			},
			want: rd,
		},
		{
			name: "flattens the response without ca",
			args: args{
				d: rdWithoutCA,
				res: &describeOutputs{
					describeSSLCertificatesOutput: &computing.DescribeSslCertificatesOutput{
						CertsSet: []types.CertsSet{
							{
								FqdnId:      nifcloud.String("test_fqdn_id"),
								Fqdn:        nifcloud.String("test_fqdn"),
								Description: nifcloud.String("test_description"),
							},
						},
					},
					downloadSSLCertificateOutputForCert: &computing.DownloadSslCertificateOutput{
						FileData: nifcloud.String("test_certificate"),
					},
					downloadSSLCertificateOutputForKey: &computing.DownloadSslCertificateOutput{
						FileData: nifcloud.String("test_key"),
					},
				},
			},
			want: rdWithoutCA,
		},
		{
			name: "flattens the response even when the resource has been removed externally",
			args: args{
				d: wantNotFoundRd,
				res: &describeOutputs{
					describeSSLCertificatesOutput: &computing.DescribeSslCertificatesOutput{
						CertsSet: []types.CertsSet{},
					},
				},
			},
			want: wantNotFoundRd,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := flatten(tt.args.d, tt.args.res)
			assert.NoError(t, err)

			if tt.args.res == nil {
				return
			}

			wantState := tt.want.State()
			if wantState == nil {
				tt.want.SetId("some")
				wantState = tt.want.State()
			}

			gotState := tt.args.d.State()
			if gotState == nil {
				tt.args.d.SetId("some")
				gotState = tt.args.d.State()
			}

			assert.Equal(t, wantState.Attributes, gotState.Attributes)
		})
	}
}
