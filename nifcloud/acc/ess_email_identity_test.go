package acc

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/ess"
	"github.com/nifcloud/nifcloud-sdk-go/service/ess/types"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
	"golang.org/x/sync/errgroup"
)

func init() {
	resource.AddTestSweepers("nifcloud_ess_email_identity", &resource.Sweeper{
		Name: "nifcloud_ess__email_identity",
		F:    testSweepESSEmail,
	})
}

func TestAcc_ESSEmail(t *testing.T) {
	var entry types.VerificationAttributes

	resourceName := "nifcloud_ess_email_identity.basic"
	randName := prefix + acctest.RandString(7)
	email := fmt.Sprintf("no-reply@%s.example.com", randName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccESSEmailResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccESSEmail(t, "testdata/ess_email.tf", email),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckESSEmailExists(resourceName, &entry),
					testAccCheckESSEmailValues(&entry, email),
					resource.TestCheckResourceAttr(resourceName, "email", email),
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

func testAccESSEmail(t *testing.T, fileName, email string) string {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf(string(b),
		email,
	)
}

func testAccCheckESSEmailExists(n string, entry *types.VerificationAttributes) resource.TestCheckFunc {
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
		res, err := svc.GetIdentityVerificationAttributes(context.Background(), &ess.GetIdentityVerificationAttributesInput{
			Identities: []string{saved.Primary.ID},
		})
		if err != nil {
			return err
		}

		if len(res.VerificationAttributes) == 0 {
			return fmt.Errorf("ess email does not found in cloud: %s", saved.Primary.ID)
		}

		foundEntry := res.VerificationAttributes[0]

		if nifcloud.ToString(foundEntry.Key) != saved.Primary.ID {
			return fmt.Errorf("ess email does not found in cloud: %s", saved.Primary.ID)
		}

		*entry = foundEntry
		return nil
	}
}

func testAccCheckESSEmailValues(entry *types.VerificationAttributes, email string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(entry.Key) != email {
			return fmt.Errorf("bad email state, expected %s, got: %#v", email, entry.Key)
		}
		return nil
	}
}

func testAccESSEmailResourceDestroy(s *terraform.State) error {
	svc := testAccProvider.Meta().(*client.Client).ESS

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nifcloud_ess_email_identity" {
			continue
		}

		res, err := svc.GetIdentityVerificationAttributes(context.Background(), &ess.GetIdentityVerificationAttributesInput{
			Identities: []string{rs.Primary.ID},
		})

		if err != nil {
			return fmt.Errorf("failed GetIdentityVerificationAttributesRequest: %s", err)
		}

		if len(res.VerificationAttributes) > 0 {
			return fmt.Errorf("ess email (%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testSweepESSEmail(region string) error {
	ctx := context.Background()
	svc := sharedClientForRegion(region).ESS

	res, err := svc.ListIdentities(ctx, nil)
	if err != nil {
		return err
	}

	var sweepIdentities []string
	for _, identity := range res.Identities {
		if strings.Contains(identity, prefix) {
			sweepIdentities = append(sweepIdentities, identity)
		}
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, n := range sweepIdentities {
		identity := n
		eg.Go(func() error {
			input := &ess.DeleteIdentityInput{
				Identity: nifcloud.String(identity),
			}
			_, err := svc.DeleteIdentity(ctx, input)
			return err
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}

	return err

}
