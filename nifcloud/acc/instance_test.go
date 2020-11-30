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
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
	"golang.org/x/sync/errgroup"
)

func init() {
	resource.AddTestSweepers("nifcloud_instance", &resource.Sweeper{
		Name: "nifcloud_instance",
		F:    testSweepInstance,
		Dependencies: []string{
			"nifcloud_volume",
		},
	})
}

func TestAcc_Instance(t *testing.T) {
	var instance computing.InstancesSet

	resourceName := "nifcloud_instance.basic"
	randName := prefix + acctest.RandStringFromCharSet(7, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccInstanceResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInstance(t, "testdata/instance.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(resourceName, &instance),
					testAccCheckInstanceValues(&instance, randName),
					resource.TestCheckResourceAttr(resourceName, "instance_id", randName),
					resource.TestCheckResourceAttr(resourceName, "description", "memo"),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "east-21"),
					resource.TestCheckResourceAttr(resourceName, "accounting_type", "2"),
					resource.TestCheckResourceAttr(resourceName, "disable_api_termination", "true"),
					resource.TestCheckResourceAttr(resourceName, "image_id", "221"),
					resource.TestCheckResourceAttr(resourceName, "instance_type", "mini"),
					resource.TestCheckResourceAttr(resourceName, "key_name", randName),
					resource.TestCheckResourceAttr(resourceName, "network_interface.1.network_id", "net-COMMON_GLOBAL"),
					resource.TestCheckResourceAttrSet(resourceName, "network_interface.1.ip_address"),
					resource.TestCheckResourceAttr(resourceName, "network_interface.0.network_id", "net-COMMON_PRIVATE"),
					resource.TestCheckResourceAttrSet(resourceName, "network_interface.0.ip_address"),
					resource.TestCheckResourceAttr(resourceName, "security_group", randName),
					resource.TestCheckResourceAttr(resourceName, "user_data", "#!/bin/bash"),
					resource.TestCheckResourceAttrSet(resourceName, "instance_state"),
					resource.TestCheckResourceAttrSet(resourceName, "public_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "private_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "unique_id"),
				),
			},
			{
				Config: testAccInstance(t, "testdata/instance_update.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(resourceName, &instance),
					testAccCheckInstanceValuesUpdated(&instance, randName),
					resource.TestCheckResourceAttr(resourceName, "instance_id", randName+"upd"),
					resource.TestCheckResourceAttr(resourceName, "description", "memo"+"-upd"),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "east-21"),
					resource.TestCheckResourceAttr(resourceName, "accounting_type", "1"),
					resource.TestCheckResourceAttr(resourceName, "disable_api_termination", "false"),
					resource.TestCheckResourceAttr(resourceName, "image_id", "221"),
					resource.TestCheckResourceAttr(resourceName, "instance_type", "small"),
					resource.TestCheckResourceAttr(resourceName, "key_name", randName),
					resource.TestCheckResourceAttr(resourceName, "network_interface.0.network_id", "net-COMMON_PRIVATE"),
					resource.TestCheckResourceAttrSet(resourceName, "network_interface.0.ip_address"),
					resource.TestCheckResourceAttr(resourceName, "security_group", randName+"upd"),
					resource.TestCheckResourceAttr(resourceName, "user_data", "#!/bin/bash"),
					resource.TestCheckResourceAttrSet(resourceName, "instance_state"),
					resource.TestCheckResourceAttrSet(resourceName, "private_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "unique_id"),
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

func TestAcc_Instance_windows(t *testing.T) {
	var instance computing.InstancesSet

	resourceName := "nifcloud_instance.basic"
	randName := prefix + acctest.RandStringFromCharSet(7, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccInstanceResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceWindows(t, "testdata/instance_windows.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "instance_id", randName),
					resource.TestCheckResourceAttr(resourceName, "image_id", "189"),
					resource.TestCheckResourceAttr(resourceName, "admin", "testadmin"),
					resource.TestCheckResourceAttr(resourceName, "password", "testpassword"),
					resource.TestCheckResourceAttr(resourceName, "license_name", "RDS"),
					resource.TestCheckResourceAttr(resourceName, "license_num", "1"),
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

func testAccInstance(t *testing.T, fileName, rName string) string {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf(string(b),
		rName,
		rName,
		rName,
	)
}

func testAccInstanceWindows(t *testing.T, fileName, rName string) string {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf(string(b),
		rName,
	)
}

func testAccCheckInstanceExists(n string, instance *computing.InstancesSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no instance resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no instance id is set")
		}

		svc := testAccProvider.Meta().(*client.Client).Computing
		res, err := svc.DescribeInstancesRequest(&computing.DescribeInstancesInput{
			InstanceId: []string{saved.Primary.ID},
		}).Send(context.Background())

		if err != nil {
			return err
		}

		if res == nil || len(res.ReservationSet) == 0 || len(res.ReservationSet[0].InstancesSet) == 0 {
			return fmt.Errorf("instance does not found in cloud: %s", saved.Primary.ID)
		}

		foundInstance := res.ReservationSet[0].InstancesSet[0]

		if nifcloud.StringValue(foundInstance.InstanceId) != saved.Primary.ID {
			return fmt.Errorf("instance does not found in cloud: %s", saved.Primary.ID)
		}

		*instance = foundInstance
		return nil
	}
}

func testAccCheckInstanceValues(instance *computing.InstancesSet, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.StringValue(instance.NextMonthAccountingType) != "2" {
			return fmt.Errorf("bad accounting_type state, expected \"2\", got: %#v", instance.NextMonthAccountingType)
		}

		if nifcloud.StringValue(instance.Placement.AvailabilityZone) != "east-21" {
			return fmt.Errorf("bad availability_zone state,  expected \"east-21\", got: %#v", instance.Placement.AvailabilityZone)
		}

		if nifcloud.StringValue(instance.Description) != "memo" {
			return fmt.Errorf("bad description state,  expected \"memo\", got: %#v", instance.Description)
		}

		if nifcloud.StringValue(instance.ImageId) != "221" {
			return fmt.Errorf("bad image_id state,  expected \"221\", got: %#v", instance.ImageId)
		}

		if nifcloud.StringValue(instance.InstanceId) != rName {
			return fmt.Errorf("bad instance_id state,  expected \"%s\", got: %#v", rName, instance.InstanceId)
		}

		if nifcloud.StringValue(instance.InstanceType) != "mini" {
			return fmt.Errorf("bad instance_type state,  expected \"mini\", got: %#v", instance.InstanceType)
		}

		if nifcloud.StringValue(instance.KeyName) != rName {
			return fmt.Errorf("bad key_name state,  expected \"%s\", got: %#v", rName, instance.KeyName)
		}

		if nifcloud.StringValue(instance.PrivateIpAddress) == "" {
			return fmt.Errorf("bad private_ip state,  expected not nil, got: nil")
		}

		if nifcloud.StringValue(instance.IpAddress) == "" {
			return fmt.Errorf("bad public_ip state,  expected not nil, got: nil")
		}

		if nifcloud.StringValue(instance.InstanceUniqueId) == "" {
			return fmt.Errorf("bad unique_id state,  expected not nil, got: nil")
		}

		if nifcloud.StringValue(instance.InstanceState.Name) == "" {
			return fmt.Errorf("bad instance_state state,  expected not nil, got: nil")
		}

		if nifcloud.StringValue(instance.NetworkInterfaceSet[0].NiftyNetworkId) != "net-COMMON_GLOBAL" {
			return fmt.Errorf("bad network_interface.0.network_id state,  expected net-COMMON_GLOBAL, got: %#v", instance.NetworkInterfaceSet[0].NiftyNetworkId)
		}

		if nifcloud.StringValue(instance.NetworkInterfaceSet[1].NiftyNetworkId) != "net-COMMON_PRIVATE" {
			return fmt.Errorf("bad network_interface.1.network_id state,  expected net-COMMON_GLOBAL, got: %#v", instance.NetworkInterfaceSet[1].NiftyNetworkId)
		}

		if nifcloud.StringValue(instance.NetworkInterfaceSet[1].PrivateIpAddress) == "" {
			return fmt.Errorf("bad network_interface.1.ip_address state,  expected not nil, got: nil")
		}
		return nil
	}
}

func testAccCheckInstanceValuesUpdated(instance *computing.InstancesSet, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.StringValue(instance.NextMonthAccountingType) != "1" {
			return fmt.Errorf("bad accounting_type state, expected \"1\", got: %#v", instance.NextMonthAccountingType)
		}

		if nifcloud.StringValue(instance.Placement.AvailabilityZone) != "east-21" {
			return fmt.Errorf("bad availability_zone state,  expected \"east-21\", got: %#v", instance.Placement.AvailabilityZone)
		}

		if nifcloud.StringValue(instance.Description) != "memo-upd" {
			return fmt.Errorf("bad description state,  expected \"memo-upd\", got: %#v", instance.Description)
		}

		if nifcloud.StringValue(instance.ImageId) != "221" {
			return fmt.Errorf("bad image_id state,  expected \"221\", got: %#v", instance.ImageId)
		}

		if nifcloud.StringValue(instance.InstanceId) != rName+"upd" {
			return fmt.Errorf("bad instance_id state,  expected \"%s\", got: %#v", rName+"upd", instance.InstanceId)
		}

		if nifcloud.StringValue(instance.InstanceType) != "small" {
			return fmt.Errorf("bad instance_type state,  expected \"small\", got: %#v", instance.InstanceType)
		}

		if nifcloud.StringValue(instance.KeyName) != rName {
			return fmt.Errorf("bad key_name state,  expected \"%s\", got: %#v", rName, instance.KeyName)
		}

		if nifcloud.StringValue(instance.PrivateIpAddress) == "" {
			return fmt.Errorf("bad private_ip state,  expected not nil, got: nil")
		}

		if nifcloud.StringValue(instance.IpAddress) != "" {
			return fmt.Errorf("bad public_ip state,  expected nil, got: not nil")
		}

		if nifcloud.StringValue(instance.InstanceUniqueId) == "" {
			return fmt.Errorf("bad unique_id state,  expected not nil, got: nil")
		}

		if nifcloud.StringValue(instance.InstanceState.Name) == "" {
			return fmt.Errorf("bad instance_state state,  expected not nil, got: nil")
		}

		if nifcloud.StringValue(instance.NetworkInterfaceSet[0].NiftyNetworkId) != "net-COMMON_PRIVATE" {
			return fmt.Errorf("bad network_interface.1.network_id state,  expected net-COMMON_GLOBAL, got: %#v", instance.NetworkInterfaceSet[1].NiftyNetworkId)
		}

		if nifcloud.StringValue(instance.NetworkInterfaceSet[0].PrivateIpAddress) == "" {
			return fmt.Errorf("bad network_interface.1.ip_address state,  expected not nil, got: nil")
		}

		return nil
	}

}

func testAccInstanceResourceDestroy(s *terraform.State) error {
	svc := testAccProvider.Meta().(*client.Client).Computing

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nifcloud_instance" {
			continue
		}

		res, err := svc.DescribeInstancesRequest(&computing.DescribeInstancesInput{
			InstanceId: []string{rs.Primary.ID},
		}).Send(context.Background())

		if err != nil {
			var awsErr awserr.Error
			if errors.As(err, &awsErr) && awsErr.Code() != "Client.InvalidParameterNotFound.Instance" {
				return fmt.Errorf("failed DescribeSInstancesRequest: %s", err)
			}
		}

		if len(res.ReservationSet) > 0 {
			return fmt.Errorf("instance (%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testSweepInstance(region string) error {
	ctx := context.Background()
	svc := sharedClientForRegion(region).Computing

	res, err := svc.DescribeInstancesRequest(nil).Send(ctx)
	if err != nil {
		return err
	}

	var sweepInstances []string
	for _, r := range res.ReservationSet {
		for _, i := range r.InstancesSet {
			if strings.HasPrefix(nifcloud.StringValue(i.InstanceId), prefix) {
				sweepInstances = append(sweepInstances, nifcloud.StringValue(i.InstanceId))
			}
		}
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, n := range sweepInstances {
		instanceID := n
		eg.Go(func() error {
			_, err = svc.StopInstancesRequest(&computing.StopInstancesInput{
				InstanceId: []string{instanceID},
			}).Send(ctx)
			if err != nil {
				return err
			}

			err = svc.WaitUntilInstanceStopped(ctx, &computing.DescribeInstancesInput{
				InstanceId: []string{instanceID},
			})
			if err != nil {
				return err
			}
			_, err = svc.TerminateInstancesRequest(&computing.TerminateInstancesInput{
				InstanceId: []string{instanceID},
			}).Send(ctx)
			if err != nil {
				return err
			}

			err = svc.WaitUntilInstanceDeleted(ctx, &computing.DescribeInstancesInput{
				InstanceId: []string{instanceID},
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
