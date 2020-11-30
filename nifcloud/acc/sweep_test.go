package acc

import (
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

func sharedClientForRegion(region string) *client.Client {
	cfg := nifcloud.NewConfig(
		os.Getenv("NIFCLOUD_ACCESS_KEY_ID"),
		os.Getenv("NIFCLOUD_SECRET_ACCESS_KEY"),
		region,
	)
	cfg.Retryer = aws.NoOpRetryer{}

	client := client.New(cfg)
	return client
}
