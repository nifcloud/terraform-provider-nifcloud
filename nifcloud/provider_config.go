package nifcloud

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/smithy-go/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

type debugLogger struct{}

func (l debugLogger) Logf(classification logging.Classification, format string, v ...interface{}) {
	if len(classification) != 0 {
		format = string(classification) + " " + format
	}
	log.Printf(format, v...)
}

// configure implements schema.ConfigureContextFunc
func configure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	cfg := nifcloud.NewConfig(
		d.Get("access_key").(string),
		d.Get("secret_key").(string),
		d.Get("region").(string),
	)
	cfg.Retryer = func() aws.Retryer {
		return aws.NopRetryer{}
	}
	cfg.ClientLogMode = aws.LogRequestWithBody
	cfg.Logger = &debugLogger{}

	storageCfg := nifcloud.NewConfig(
		d.Get("storage_access_key").(string),
		d.Get("storage_secret_key").(string),
		d.Get("storage_region").(string),
	)
	storageCfg.Retryer = func() aws.Retryer {
		return aws.NopRetryer{}
	}
	storageCfg.ClientLogMode = aws.LogRequestWithBody
	storageCfg.Logger = &debugLogger{}

	client := client.New(cfg, storageCfg)
	return client, nil
}
