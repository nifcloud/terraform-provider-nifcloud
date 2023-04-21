package acc

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/aws/smithy-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing/types"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
	"golang.org/x/sync/errgroup"
)

func init() {
	resource.AddTestSweepers("nifcloud_key_pair", &resource.Sweeper{
		Name: "nifcloud_key_pair",
		F:    testSweepKeyPair,
	})
}

func TestAcc_KeyPair(t *testing.T) {
	var keyPair types.KeySet

	resourceName := "nifcloud_key_pair.basic"
	randName := prefix + acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccKeyPairResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyPair(t, "testdata/key_pair.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyPairExists(resourceName, &keyPair),
					testAccCheckKeyPairValues(&keyPair, randName),
					resource.TestCheckResourceAttr(resourceName, "key_name", randName),
					resource.TestCheckResourceAttr(resourceName, "description", "memo"),
					resource.TestCheckResourceAttrSet(resourceName, "fingerprint"),
				),
			},
			{
				Config: testAccKeyPair(t, "testdata/key_pair_update.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyPairExists(resourceName, &keyPair),
					testAccCheckKeyPairValuesUpdated(&keyPair, randName),
					resource.TestCheckResourceAttr(resourceName, "key_name", randName),
					resource.TestCheckResourceAttr(resourceName, "description", "memo-upd"),
					resource.TestCheckResourceAttrSet(resourceName, "fingerprint"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"public_key",
				},
			},
		},
	})
}

func testAccKeyPair(t *testing.T, fileName, keyName string) string {
	b, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf(string(b),
		keyName,
	)
}

func testAccCheckKeyPairExists(n string, keyPair *types.KeySet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no keyPair resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no keyPair id is set")
		}

		svc := testAccProvider.Meta().(*client.Client).Computing
		res, err := svc.DescribeKeyPairs(context.Background(), &computing.DescribeKeyPairsInput{
			KeyName: []string{saved.Primary.ID},
		})

		if err != nil {
			return err
		}

		if len(res.KeySet) == 0 {
			return fmt.Errorf("keyPair does not found in cloud: %s", saved.Primary.ID)
		}

		foundKeyPair := res.KeySet[0]

		if nifcloud.ToString(foundKeyPair.KeyName) != saved.Primary.ID {
			return fmt.Errorf("keyPair does not found in cloud: %s", saved.Primary.ID)
		}

		*keyPair = foundKeyPair
		return nil
	}
}

func testAccCheckKeyPairValues(keyPair *types.KeySet, keyName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(keyPair.KeyName) != keyName {
			return fmt.Errorf("bad key_name state, expected \"%s\", got: %#v", keyName, keyPair.KeyName)
		}

		if nifcloud.ToString(keyPair.Description) != "memo" {
			return fmt.Errorf("bad description state, expected \"memo\", got: %#v", keyPair.Description)
		}

		if nifcloud.ToString(keyPair.KeyFingerprint) == "" {
			return fmt.Errorf("bad fingerprint state, expected not nil, got: nil")
		}
		return nil
	}
}

func testAccCheckKeyPairValuesUpdated(keyPair *types.KeySet, keyName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(keyPair.KeyName) != keyName {
			return fmt.Errorf("bad key_name state, expected \"%s\", got: %#v", keyName, keyPair.KeyName)
		}

		if nifcloud.ToString(keyPair.Description) != "memo-upd" {
			return fmt.Errorf("bad description state, expected \"memo-upd\", got: %#v", keyPair.Description)
		}

		if nifcloud.ToString(keyPair.KeyFingerprint) == "" {
			return fmt.Errorf("bad fingerprint state, expected not nil, got: nil")
		}
		return nil
	}
}

func testAccKeyPairResourceDestroy(s *terraform.State) error {
	svc := testAccProvider.Meta().(*client.Client).Computing

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nifcloud_key_pair" {
			continue
		}

		res, err := svc.DescribeKeyPairs(context.Background(), &computing.DescribeKeyPairsInput{
			KeyName: []string{rs.Primary.ID},
		})

		if err != nil {
			var awsErr smithy.APIError
			if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameterNotFound.KeyPair" {
				return nil
			}
			return fmt.Errorf("failed DescribeKeyPairsRequest: %s", err)
		}

		if len(res.KeySet) > 0 {
			return fmt.Errorf("keyPair (%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testSweepKeyPair(region string) error {
	ctx := context.Background()
	svc := sharedClientForRegion(region).Computing

	res, err := svc.DescribeKeyPairs(ctx, nil)
	if err != nil {
		return err
	}

	var sweepKeyPairs []string
	for _, k := range res.KeySet {
		if strings.HasPrefix(nifcloud.ToString(k.KeyName), prefix) {
			sweepKeyPairs = append(sweepKeyPairs, nifcloud.ToString(k.KeyName))
		}
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, n := range sweepKeyPairs {
		keyName := n
		eg.Go(func() error {
			_, err := svc.DeleteKeyPair(ctx, &computing.DeleteKeyPairInput{
				KeyName: nifcloud.String(keyName),
			})
			return err
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}
