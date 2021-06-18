package acc

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/nas"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
	"golang.org/x/sync/errgroup"
)

func init() {
	resource.AddTestSweepers("nifcloud_nas_instance", &resource.Sweeper{
		Name:         "nifcloud_nas_instance",
		F:            testSweepNASInstance,
		Dependencies: []string{},
	})
}

func TestAcc_NASInstance_NFS(t *testing.T) {
	var nasInstance nas.NASInstance

	resourceName := "nifcloud_nas_instance.basic"
	randName := prefix + acctest.RandStringFromCharSet(7, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccNASInstanceResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNASInstanceForNFS(t, "testdata/nas_instance_nfs.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNASInstanceForNFSExists(resourceName, &nasInstance),
					testAccCheckNASInstanceValuesForNFS(&nasInstance, randName),
					resource.TestCheckResourceAttr(resourceName, "identifier", randName),
					resource.TestCheckResourceAttr(resourceName, "description", "memo"),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "east-21"),
					resource.TestCheckResourceAttr(resourceName, "allocated_storage", "100"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "nfs"),
					resource.TestCheckResourceAttr(resourceName, "type", "0"),
					resource.TestCheckResourceAttr(resourceName, "nas_security_group_name", randName),
					resource.TestCheckResourceAttrSet(resourceName, "public_ip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "private_ip_address"),
				),
			},
			{
				Config: testAccNASInstanceForNFS(t, "testdata/nas_instance_nfs_update.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNASInstanceForNFSExists(resourceName, &nasInstance),
					testAccCheckNASInstanceValuesUpdatedForNFS(&nasInstance, randName),
					resource.TestCheckResourceAttr(resourceName, "identifier", randName+"upd"),
					resource.TestCheckResourceAttr(resourceName, "description", "memo-upd"),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "east-21"),
					resource.TestCheckResourceAttr(resourceName, "allocated_storage", "200"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "nfs"),
					resource.TestCheckResourceAttr(resourceName, "type", "0"),
					resource.TestCheckResourceAttr(resourceName, "nas_security_group_name", randName),
					resource.TestCheckResourceAttrSet(resourceName, "public_ip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "private_ip_address"),
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

func TestAcc_NASInstance_CIFS(t *testing.T) {
	var nasInstance nas.NASInstance

	resourceName := "nifcloud_nas_instance.basic"
	randName := prefix + acctest.RandStringFromCharSet(7, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccNASInstanceResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNASInstanceForCIFS(t, "testdata/nas_instance_cifs.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNASInstanceForCIFSExists(resourceName, &nasInstance),
					testAccCheckNASInstanceValuesForCIFS(&nasInstance, randName),
					resource.TestCheckResourceAttr(resourceName, "identifier", randName),
					resource.TestCheckResourceAttr(resourceName, "description", "memo"),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "east-21"),
					resource.TestCheckResourceAttr(resourceName, "allocated_storage", "100"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "cifs"),
					resource.TestCheckResourceAttr(resourceName, "type", "0"),
					resource.TestCheckResourceAttr(resourceName, "nas_security_group_name", randName),
					resource.TestCheckResourceAttr(resourceName, "master_username", "tfacc"),
					resource.TestCheckResourceAttr(resourceName, "master_user_password", "tfaccpass"),
					resource.TestCheckResourceAttr(resourceName, "authentication_type", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "public_ip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "private_ip_address"),
				),
			},
			{
				Config: testAccNASInstanceForCIFS(t, "testdata/nas_instance_cifs_update.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNASInstanceForNFSExists(resourceName, &nasInstance),
					testAccCheckNASInstanceValuesUpdatedForCIFS(&nasInstance, randName),
					resource.TestCheckResourceAttr(resourceName, "identifier", randName+"upd"),
					resource.TestCheckResourceAttr(resourceName, "description", "memo-upd"),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "east-21"),
					resource.TestCheckResourceAttr(resourceName, "allocated_storage", "200"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "cifs"),
					resource.TestCheckResourceAttr(resourceName, "type", "0"),
					resource.TestCheckResourceAttr(resourceName, "nas_security_group_name", randName),
					resource.TestCheckResourceAttr(resourceName, "master_username", "tfacc"),
					resource.TestCheckResourceAttr(resourceName, "master_user_password", "tfaccpass"),
					resource.TestCheckResourceAttr(resourceName, "authentication_type", "1"),
					resource.TestCheckResourceAttr(resourceName, "directory_service_domain_name", "tfacc.local"),
					resource.TestCheckResourceAttr(resourceName, "directory_service_administrator_name", "Administrator"),
					resource.TestCheckResourceAttr(resourceName, "directory_service_administrator_password", "tfaccpass+555"),
					resource.TestCheckResourceAttr(resourceName, "domain_controllers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "domain_controllers.0.hostname", "ad01"),
					resource.TestCheckResourceAttrSet(resourceName, "public_ip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "private_ip_address"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"directory_service_administrator_name",
					"directory_service_administrator_password",
					"master_user_password",
				},
			},
		},
	})
}

func testAccNASInstanceForNFS(t *testing.T, fileName, rName string) string {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf(string(b),
		rName,
		rName,
	)
}

func testAccNASInstanceForCIFS(t *testing.T, fileName, rName string) string {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}

	userData, err := ioutil.ReadFile("testdata/scripts/provision_ad_server.sh")
	if err != nil {
		t.Fatal(err)
	}

	return fmt.Sprintf(string(b),
		rName,
		rName,
		rName,
		string(userData),
		rName,
		rName,
	)
}

func testAccCheckNASInstanceForNFSExists(n string, nasInstance *nas.NASInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no nasInstance resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no nasInstance id is set")
		}

		svc := testAccProvider.Meta().(*client.Client).NAS
		res, err := svc.DescribeNASInstancesRequest(&nas.DescribeNASInstancesInput{
			NASInstanceIdentifier: nifcloud.String(saved.Primary.ID),
		}).Send(context.Background())
		if err != nil {
			return err
		}

		if len(res.NASInstances) == 0 {
			return fmt.Errorf("nasInstance does not found in cloud: %s", saved.Primary.ID)
		}

		foundNASInstance := res.NASInstances[0]

		if nifcloud.StringValue(foundNASInstance.NASInstanceIdentifier) != saved.Primary.ID {
			return fmt.Errorf("nasInstance does not found in cloud: %s", saved.Primary.ID)
		}

		*nasInstance = foundNASInstance

		return nil
	}
}

func testAccCheckNASInstanceForCIFSExists(n string, nasInstance *nas.NASInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no nasInstance resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no nasInstance id is set")
		}

		svc := testAccProvider.Meta().(*client.Client).NAS
		res, err := svc.DescribeNASInstancesRequest(&nas.DescribeNASInstancesInput{
			NASInstanceIdentifier: nifcloud.String(saved.Primary.ID),
		}).Send(context.Background())
		if err != nil {
			return err
		}

		if len(res.NASInstances) == 0 {
			return fmt.Errorf("nasInstance does not found in cloud: %s", saved.Primary.ID)
		}

		foundNASInstance := res.NASInstances[0]

		if nifcloud.StringValue(foundNASInstance.NASInstanceIdentifier) != saved.Primary.ID {
			return fmt.Errorf("nasInstance does not found in cloud: %s", saved.Primary.ID)
		}

		*nasInstance = foundNASInstance

		return nil
	}
}

func testAccCheckNASInstanceValuesForNFS(nasInstance *nas.NASInstance, identifier string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.StringValue(nasInstance.NASInstanceIdentifier) != identifier {
			return fmt.Errorf("bad identifier state, expected \"%s\", got: %#v", identifier, nifcloud.StringValue(nasInstance.NASInstanceIdentifier))
		}

		if nifcloud.StringValue(nasInstance.NASInstanceDescription) != "memo" {
			return fmt.Errorf("bad description state, expected \"memo\", got: %#v", nifcloud.StringValue(nasInstance.NASInstanceDescription))
		}

		if nifcloud.StringValue(nasInstance.AvailabilityZone) != "east-21" {
			return fmt.Errorf("bad availability_zone state, expected \"east-21\", got: %#v", nifcloud.StringValue(nasInstance.AvailabilityZone))
		}

		if nifcloud.StringValue(nasInstance.Protocol) != "nfs" {
			return fmt.Errorf("bad protocol state, expected \"nfs\", got: %#v", nifcloud.StringValue(nasInstance.Protocol))
		}

		if nifcloud.StringValue(nasInstance.AllocatedStorage) != "100" {
			return fmt.Errorf("bad allocated_storage state, expected \"100\", got: %#v", nifcloud.StringValue(nasInstance.AllocatedStorage))
		}

		if nifcloud.Int64Value(nasInstance.NASInstanceType) != 0 {
			return fmt.Errorf("bad type state, expected \"0\", got: %#v", nifcloud.Int64Value(nasInstance.NASInstanceType))
		}

		if nifcloud.StringValue(nasInstance.NoRootSquash) != "false" {
			return fmt.Errorf("bad no_root_squash state, expected \"false\", got: %#v", nifcloud.StringValue(nasInstance.NoRootSquash))
		}

		return nil
	}
}

func testAccCheckNASInstanceValuesUpdatedForNFS(nasInstance *nas.NASInstance, identifier string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.StringValue(nasInstance.NASInstanceIdentifier) != identifier+"upd" {
			return fmt.Errorf("bad identifier state, expected \"%s\", got: %#v", identifier+"upd", nifcloud.StringValue(nasInstance.NASInstanceIdentifier))
		}

		if nifcloud.StringValue(nasInstance.NASInstanceDescription) != "memo-upd" {
			return fmt.Errorf("bad description state, expected \"memo-upd\", got: %#v", nifcloud.StringValue(nasInstance.NASInstanceDescription))
		}

		if nifcloud.StringValue(nasInstance.AvailabilityZone) != "east-21" {
			return fmt.Errorf("bad availability_zone state, expected \"east-21\", got: %#v", nifcloud.StringValue(nasInstance.AvailabilityZone))
		}

		if nifcloud.StringValue(nasInstance.Protocol) != "nfs" {
			return fmt.Errorf("bad protocol state, expected \"nfs\", got: %#v", nifcloud.StringValue(nasInstance.Protocol))
		}

		if nifcloud.StringValue(nasInstance.AllocatedStorage) != "200" {
			return fmt.Errorf("bad allocated_storage state, expected \"200\", got: %#v", nifcloud.StringValue(nasInstance.AllocatedStorage))
		}

		if nifcloud.Int64Value(nasInstance.NASInstanceType) != 0 {
			return fmt.Errorf("bad type state, expected \"0\", got: %#v", nifcloud.Int64Value(nasInstance.NASInstanceType))
		}

		if nifcloud.StringValue(nasInstance.NoRootSquash) != "true" {
			return fmt.Errorf("bad no_root_squash state, expected \"true\", got: %#v", nifcloud.StringValue(nasInstance.NoRootSquash))
		}

		return nil
	}
}

func testAccCheckNASInstanceValuesForCIFS(nasInstance *nas.NASInstance, identifier string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.StringValue(nasInstance.NASInstanceIdentifier) != identifier {
			return fmt.Errorf("bad identifier state, expected \"%s\", got: %#v", identifier, nifcloud.StringValue(nasInstance.NASInstanceIdentifier))
		}

		if nifcloud.StringValue(nasInstance.NASInstanceDescription) != "memo" {
			return fmt.Errorf("bad description state, expected \"memo\", got: %#v", nifcloud.StringValue(nasInstance.NASInstanceDescription))
		}

		if nifcloud.StringValue(nasInstance.AvailabilityZone) != "east-21" {
			return fmt.Errorf("bad availability_zone state, expected \"east-21\", got: %#v", nifcloud.StringValue(nasInstance.AvailabilityZone))
		}

		if nifcloud.StringValue(nasInstance.Protocol) != "cifs" {
			return fmt.Errorf("bad protocol state, expected \"cifs\", got: %#v", nifcloud.StringValue(nasInstance.Protocol))
		}

		if nifcloud.StringValue(nasInstance.AllocatedStorage) != "100" {
			return fmt.Errorf("bad allocated_storage state, expected \"100\", got: %#v", nifcloud.StringValue(nasInstance.AllocatedStorage))
		}

		if nifcloud.Int64Value(nasInstance.NASInstanceType) != 0 {
			return fmt.Errorf("bad type state, expected \"0\", got: %#v", nifcloud.Int64Value(nasInstance.NASInstanceType))
		}

		if nifcloud.StringValue(nasInstance.MasterUsername) != "tfacc" {
			return fmt.Errorf("bad master_username state, expected \"tfacc\", got: %#v", nifcloud.StringValue(nasInstance.MasterUsername))
		}

		if nifcloud.Int64Value(nasInstance.AuthenticationType) != 0 {
			return fmt.Errorf("bad authentication_type state, expected 0, got: %#v", nifcloud.Int64Value(nasInstance.AuthenticationType))
		}

		return nil
	}
}

func testAccCheckNASInstanceValuesUpdatedForCIFS(nasInstance *nas.NASInstance, identifier string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.StringValue(nasInstance.NASInstanceIdentifier) != identifier+"upd" {
			return fmt.Errorf("bad identifier state, expected \"%s\", got: %#v", identifier+"upd", nifcloud.StringValue(nasInstance.NASInstanceIdentifier))
		}

		if nifcloud.StringValue(nasInstance.NASInstanceDescription) != "memo-upd" {
			return fmt.Errorf("bad description state, expected \"memo-upd\", got: %#v", nifcloud.StringValue(nasInstance.NASInstanceDescription))
		}

		if nifcloud.StringValue(nasInstance.AvailabilityZone) != "east-21" {
			return fmt.Errorf("bad availability_zone state, expected \"east-21\", got: %#v", nifcloud.StringValue(nasInstance.AvailabilityZone))
		}

		if nifcloud.StringValue(nasInstance.Protocol) != "cifs" {
			return fmt.Errorf("bad protocol state, expected \"cifs\", got: %#v", nifcloud.StringValue(nasInstance.Protocol))
		}

		if nifcloud.StringValue(nasInstance.AllocatedStorage) != "200" {
			return fmt.Errorf("bad allocated_storage state, expected \"200\", got: %#v", nifcloud.StringValue(nasInstance.AllocatedStorage))
		}

		if nifcloud.Int64Value(nasInstance.NASInstanceType) != 0 {
			return fmt.Errorf("bad type state, expected \"0\", got: %#v", nifcloud.Int64Value(nasInstance.NASInstanceType))
		}

		if nifcloud.StringValue(nasInstance.MasterUsername) != "tfacc" {
			return fmt.Errorf("bad master_username state, expected \"tfacc\", got: %#v", nifcloud.StringValue(nasInstance.MasterUsername))
		}

		if nifcloud.Int64Value(nasInstance.AuthenticationType) != 1 {
			return fmt.Errorf("bad authentication_type state, expected 1, got: %#v", nifcloud.Int64Value(nasInstance.AuthenticationType))
		}

		if nifcloud.StringValue(nasInstance.DirectoryServiceDomainName) != "tfacc.local" {
			return fmt.Errorf("bad directory_service_domain_name state, expected \"tfacc.local\" got: %#v", nifcloud.StringValue(nasInstance.DirectoryServiceDomainName))
		}

		if len(nasInstance.DomainControllers) != 1 {
			return fmt.Errorf("bad domain_controllers state, expected length is 1, got: %d", len(nasInstance.DomainControllers))
		}

		if nifcloud.StringValue(nasInstance.DomainControllers[0].Hostname) != "ad01" {
			return fmt.Errorf("bad domain_controllers.0.hostname state, expected \"ad01\" got: %#v", nifcloud.StringValue(nasInstance.DomainControllers[0].Hostname))
		}

		return nil
	}
}

func testAccNASInstanceResourceDestroy(s *terraform.State) error {
	svc := testAccProvider.Meta().(*client.Client).NAS

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nifcloud_nas_instance" {
			continue
		}

		res, err := svc.DescribeNASInstancesRequest(&nas.DescribeNASInstancesInput{
			NASInstanceIdentifier: nifcloud.String(rs.Primary.ID),
		}).Send(context.Background())

		if err != nil {
			var awsErr awserr.Error
			if errors.As(err, &awsErr) && awsErr.Code() == "Client.InvalidParameter.NotFound.NASInstanceIdentifier" {
				return nil
			}
			return fmt.Errorf("failed NiftyDescribeNatTablesRequest: %s", err)
		}

		if len(res.NASInstances) > 0 {
			return fmt.Errorf("nasInstance (%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testSweepNASInstance(region string) error {
	ctx := context.Background()
	svc := sharedClientForRegion(region).NAS

	res, err := svc.DescribeNASInstancesRequest(nil).Send(ctx)
	if err != nil {
		return err
	}

	var sweepNASInstances []string
	for _, g := range res.NASInstances {
		if strings.HasPrefix(nifcloud.StringValue(g.NASInstanceIdentifier), prefix) {
			sweepNASInstances = append(sweepNASInstances, nifcloud.StringValue(g.NASInstanceIdentifier))
		}
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, n := range sweepNASInstances {
		groupName := n
		eg.Go(func() error {
			_, err := svc.DeleteNASInstanceRequest(&nas.DeleteNASInstanceInput{
				NASInstanceIdentifier: nifcloud.String(groupName),
			}).Send(ctx)
			return err
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}
