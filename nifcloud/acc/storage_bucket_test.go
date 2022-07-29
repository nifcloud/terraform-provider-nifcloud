package acc

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/aws/smithy-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/storage"
	"github.com/nifcloud/nifcloud-sdk-go/service/storage/types"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
	"golang.org/x/sync/errgroup"
)

func init() {
	resource.AddTestSweepers("nifcloud_storage_bucket", &resource.Sweeper{
		Name: "nifcloud_storage_bucket",
		F:    testSweepStorageBucket,
	})
}

func TestAcc_StorageBucket(t *testing.T) {
	var bucket types.Buckets

	resourceName := "nifcloud_storage_bucket.basic"
	randName := prefix + acctest.RandString(7)
	policy := fmt.Sprintf(`{"Statement":[{"Effect":"Allow","Principal":"*","Action":"s3:GetObject","Resource":"urn:sgws:s3:::%s/index.html"}]}`, randName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccStorageBucketResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket(t, "testdata/storage_bucket.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(resourceName, &bucket),
					testAccCheckStorageBucketValues(&bucket, randName),
					resource.TestCheckResourceAttr(resourceName, "bucket", randName),
					resource.TestCheckResourceAttr(resourceName, "versioning.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "versioning.0.enabled", "false"),
				),
			},
			{
				Config: testAccStorageBucketUpdated(t, "testdata/storage_bucket_update.tf", randName, policy),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(resourceName, &bucket),
					testAccCheckStorageBucketValuesUpdated(&bucket, randName, policy),
					resource.TestCheckResourceAttr(resourceName, "bucket", randName),
					resource.TestCheckResourceAttr(resourceName, "versioning.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "versioning.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "policy", policy),
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

func testAccStorageBucket(t *testing.T, fileName, rName string) string {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf(string(b),
		rName,
	)
}

func testAccStorageBucketUpdated(t *testing.T, fileName, rName, policy string) string {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf(string(b),
		rName,
		policy,
	)
}

func testAccCheckStorageBucketExists(n string, bucket *types.Buckets) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no storage bucket resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no storage bucket id is set")
		}

		svc := testAccProvider.Meta().(*client.Client).Storage
		res, err := svc.GetService(context.Background(), nil)
		if err != nil {
			return fmt.Errorf("failed GetService request: %w", err)
		}

		if res == nil || len(res.Buckets) == 0 {
			return fmt.Errorf("storage bucket does not found in cloud: %s", saved.Primary.ID)
		}

		var foundBucket *types.Buckets
		for i := range res.Buckets {
			if nifcloud.ToString(res.Buckets[i].Name) == saved.Primary.ID {
				foundBucket = &res.Buckets[i]
				break
			}
		}

		if foundBucket == nil {
			return fmt.Errorf("storage bucket does not found in cloud: %s", saved.Primary.ID)
		}
		*bucket = *foundBucket

		return nil
	}
}

func testAccCheckStorageBucketValues(bucket *types.Buckets, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(bucket.Name) != rName {
			return fmt.Errorf("bad bucket state, expected \"%s\", got: %#v", rName, nifcloud.ToString(bucket.Name))
		}

		svc := testAccProvider.Meta().(*client.Client).Storage
		versioningRes, err := svc.GetBucketVersioning(context.Background(), &storage.GetBucketVersioningInput{
			Bucket: nifcloud.String(rName),
		})
		if err != nil {
			return fmt.Errorf("failed GetBucketVersioning request: %w", err)
		}

		if nifcloud.ToString(versioningRes.Status) != "" {
			return fmt.Errorf("bad versioning.0.enabled state, expected empty string, got: %#v", nifcloud.ToString(versioningRes.Status))
		}

		policyRes, err := svc.GetBucketPolicy(context.Background(), &storage.GetBucketPolicyInput{
			Bucket: nifcloud.String(rName),
		})
		if err != nil {
			var awsErr smithy.APIError
			if errors.As(err, &awsErr) {
				if awsErr.ErrorCode() != "NoSuchBucketPolicy" {
					return fmt.Errorf("bad policy state, expected no bucket policy, got %#v", awsErr.ErrorCode())
				}
			} else {
				return fmt.Errorf("failed GetBucketPolicy request: %w", err)
			}
		} else {
			return fmt.Errorf("bad policy state, expected no bucket policy, got: %#v", nifcloud.ToString(policyRes.Policy))
		}

		return nil
	}
}

func testAccCheckStorageBucketValuesUpdated(bucket *types.Buckets, rName, policy string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(bucket.Name) != rName {
			return fmt.Errorf("bad bucket state, expected \"%s\", got: %#v", rName, nifcloud.ToString(bucket.Name))
		}

		svc := testAccProvider.Meta().(*client.Client).Storage
		versioningRes, err := svc.GetBucketVersioning(context.Background(), &storage.GetBucketVersioningInput{
			Bucket: nifcloud.String(rName),
		})
		if err != nil {
			return fmt.Errorf("failed GetBucketVersioning request: %w", err)
		}

		if nifcloud.ToString(versioningRes.Status) != "Enabled" {
			return fmt.Errorf("bad versioning.0.enabled state,  expected \"Enabled\", got: %#v", nifcloud.ToString(versioningRes.Status))
		}

		policyRes, err := svc.GetBucketPolicy(context.Background(), &storage.GetBucketPolicyInput{
			Bucket: nifcloud.String(rName),
		})
		if err != nil {
			return fmt.Errorf("failed GetBucketPolicy request: %w", err)
		}

		if nifcloud.ToString(policyRes.Policy) != policy {
			return fmt.Errorf("bad policy state, expected %s, got: %#v", policy, nifcloud.ToString(policyRes.Policy))
		}

		return nil
	}

}

func testAccStorageBucketResourceDestroy(s *terraform.State) error {
	svc := testAccProvider.Meta().(*client.Client).Storage

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nifcloud_storage_bucket" {
			continue
		}

		res, err := svc.GetService(context.Background(), nil)
		if err != nil {
			return fmt.Errorf("failed GetService request: %w", err)
		}

		for _, b := range res.Buckets {
			if nifcloud.ToString(b.Name) == rs.Primary.ID {
				return fmt.Errorf("instance (%s) still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testSweepStorageBucket(region string) error {
	ctx := context.Background()
	svc := sharedClientForRegion(region).Storage

	res, err := svc.GetService(ctx, nil)
	if err != nil {
		return err
	}

	var sweepBuckets []string
	for _, b := range res.Buckets {
		if strings.HasPrefix(nifcloud.ToString(b.Name), prefix) {
			sweepBuckets = append(sweepBuckets, nifcloud.ToString(b.Name))
		}
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, b := range sweepBuckets {
		bucketName := b
		eg.Go(func() error {
			_, err = svc.DeleteBucket(ctx, &storage.DeleteBucketInput{
				Bucket: nifcloud.String(bucketName),
			})
			if err != nil {
				return err
			}

			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}
