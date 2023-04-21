package acc

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/ess"
	"github.com/nifcloud/nifcloud-sdk-go/service/ess/types"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func TestAcc_ESSDomainDkim(t *testing.T) {
	var entry types.DkimAttributes

	resourceName := "nifcloud_ess_domain_dkim.basic"
	randName := prefix + acctest.RandString(7)
	domain := fmt.Sprintf("%s.example.com", randName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccESSDomainDkimResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccESSDomainDkim(t, "testdata/ess_domain_dkim.tf", domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckESSDomainDkimExists(resourceName, &entry),
					testAccCheckESSDomainDkimValues(domain, &entry),
					resource.TestCheckResourceAttr(resourceName, "domain", domain),
					resource.TestCheckResourceAttrSet(resourceName, "dkim_tokens.#"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccESSDomainDkim(t *testing.T, fileName, domain string) string {
	b, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf(string(b),
		domain,
	)
}

func testAccCheckESSDomainDkimExists(n string, entry *types.DkimAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// lintignore:R018
		time.Sleep(10 * time.Second)
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no id is set")
		}

		svc := testAccProvider.Meta().(*client.Client).ESS
		res, err := svc.GetIdentityDkimAttributes(context.Background(), &ess.GetIdentityDkimAttributesInput{
			Identities: []string{saved.Primary.ID},
		})
		if err != nil {
			return err
		}

		if len(res.DkimAttributes) == 0 {
			return fmt.Errorf("ess domain dkim does not found in cloud: %s", saved.Primary.ID)
		}

		foundEntry := res.DkimAttributes[0]

		if nifcloud.ToString(foundEntry.Key) != saved.Primary.ID {
			return fmt.Errorf("ess domain dkim does not found in cloud: %s", saved.Primary.ID)
		}

		*entry = foundEntry
		return nil
	}
}

func testAccCheckESSDomainDkimValues(domain string, entry *types.DkimAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(entry.Key) != domain {
			return fmt.Errorf("bad domain state, expected %s, got: %#v", domain, nifcloud.ToString(entry.Key))
		}
		if len(entry.Value.DkimTokens) != 3 {
			return fmt.Errorf("bad dkim_tokens state, expected 3 tokens, got: %#v", len(entry.Value.DkimTokens))
		}
		return nil
	}
}

func testAccESSDomainDkimResourceDestroy(s *terraform.State) error {
	svc := testAccProvider.Meta().(*client.Client).ESS

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nifcloud_ess_domain_dkim" {
			continue
		}

		res, err := svc.GetIdentityDkimAttributes(context.Background(), &ess.GetIdentityDkimAttributesInput{
			Identities: []string{rs.Primary.ID},
		})

		if err != nil {
			return fmt.Errorf("failed GetIdentityDkimAttributesRequest: %s", err)
		}

		if len(res.DkimAttributes) > 0 {
			return fmt.Errorf("ess domain dkim (%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}
