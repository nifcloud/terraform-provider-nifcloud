package nifcloud

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

type debugLogger struct{}

func (l debugLogger) Log(v ...interface{}) {
	log.Println(v...)
}

// configure implements schema.ConfigureContextFunc
func configure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	cfg := nifcloud.NewConfig(
		d.Get("access_key").(string),
		d.Get("secret_key").(string),
		d.Get("region").(string),
	)
	cfg.Retryer = aws.NoOpRetryer{}
	cfg.LogLevel = aws.LogDebugWithHTTPBody
	cfg.Logger = &debugLogger{}

	client := client.New(cfg)
	return client, nil
}
