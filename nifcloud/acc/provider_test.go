package acc

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud"
)

const (
	prefix = "tfacc"
)

var testAccExternalProviders = map[string]resource.ExternalProvider{
	"tls": {
		Source: "registry.terraform.io/hashicorp/tls",
	},
	"randam": {
		Source: "registry.terraform.io/hashicorp/random",
	},
}

var testAccProviderFactory = map[string]func() (*schema.Provider, error){
	"nifcloud": providerFactory,
}
var testAccProvider = nifcloud.Provider()

func providerFactory() (*schema.Provider, error) {
	return testAccProvider, nil
}

func TestProvider(t *testing.T) {
	if err := nifcloud.Provider().InternalValidate(); err != nil {
		t.Fatal(err.Error())
	}
}

func testAccPreCheck(t *testing.T) {
	var region, accessKey, secretKey string
	if k := os.Getenv("NIFCLOUD_DEFAULT_REGION"); k != "" {
		region = k
	}

	if k := os.Getenv("NIFCLOUD_ACCESS_KEY_ID"); k != "" {
		accessKey = k
	}

	if k := os.Getenv("NIFCLOUD_SECRET_ACCESS_KEY"); k != "" {
		secretKey = k
	}

	if region == "" || accessKey == "" || secretKey == "" {
		t.Fatal("No valid credentials found to execute acceptance tests")
	}
}
