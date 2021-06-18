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
	resource.AddTestSweepers("nifcloud_volume", &resource.Sweeper{
		Name: "nifcloud_volume",
		F:    testSweepVolume,
	})
}

func TestAcc_Volume(t *testing.T) {
	var volume computing.VolumeSet

	resourceName := "nifcloud_volume.basic"
	randName := prefix + acctest.RandStringFromCharSet(7, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccVolumeResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVolume(t, "testdata/volume.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists(resourceName, &volume),
					testAccCheckVolumeValues(&volume, randName),
					resource.TestCheckResourceAttr(resourceName, "size", "100"),
					resource.TestCheckResourceAttr(resourceName, "volume_id", randName),
					resource.TestCheckResourceAttr(resourceName, "disk_type", "High-Speed Storage A"),
					resource.TestCheckResourceAttr(resourceName, "instance_id", randName),
					resource.TestCheckResourceAttr(resourceName, "reboot", "true"),
					resource.TestCheckResourceAttr(resourceName, "accounting_type", "1"),
					resource.TestCheckResourceAttr(resourceName, "description", "memo"),
				),
			},
			{
				Config: testAccVolume(t, "testdata/volume_update.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists(resourceName, &volume),
					testAccCheckVolumeValuesUpdated(&volume, randName),
					resource.TestCheckResourceAttr(resourceName, "size", "300"),
					resource.TestCheckResourceAttr(resourceName, "volume_id", randName+"upd"),
					resource.TestCheckResourceAttr(resourceName, "disk_type", "High-Speed Storage A"),
					resource.TestCheckResourceAttr(resourceName, "instance_id", randName),
					resource.TestCheckResourceAttr(resourceName, "reboot", "true"),
					resource.TestCheckResourceAttr(resourceName, "accounting_type", "2"),
					resource.TestCheckResourceAttr(resourceName, "description", "memo"+"-upd"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"instance_id",
					"reboot",
				},
			},
		},
	})
}

func TestAcc_Volume_Unique_Id(t *testing.T) {
	var volume computing.VolumeSet

	resourceName := "nifcloud_volume.basic"
	randName := prefix + acctest.RandStringFromCharSet(7, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccVolumeResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVolume(t, "testdata/volume_unique_id.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists(resourceName, &volume),
					testAccCheckVolumeUniqueIDValues(&volume, randName),
					resource.TestCheckResourceAttr(resourceName, "size", "100"),
					resource.TestCheckResourceAttr(resourceName, "volume_id", randName),
					resource.TestCheckResourceAttr(resourceName, "disk_type", "High-Speed Storage A"),
					resource.TestCheckResourceAttrSet(resourceName, "instance_unique_id"),
					resource.TestCheckResourceAttr(resourceName, "reboot", "true"),
					resource.TestCheckResourceAttr(resourceName, "accounting_type", "1"),
					resource.TestCheckResourceAttr(resourceName, "description", "memo"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"instance_unique_id",
					"reboot",
				},
			},
		},
	})
}

func testAccVolume(t *testing.T, fileName, rName string) string {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf(string(b),
		rName,
		rName,
		rName,
		rName,
	)
}

func testAccCheckVolumeExists(n string, volume *computing.VolumeSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no volume resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no volume id is set")
		}

		svc := testAccProvider.Meta().(*client.Client).Computing
		res, err := svc.DescribeVolumesRequest(&computing.DescribeVolumesInput{
			VolumeId: []string{saved.Primary.ID},
		}).Send(context.Background())

		if err != nil {
			return err
		}
		if res == nil || len(res.VolumeSet) == 0 {
			return fmt.Errorf("volume does not found in cloud: %s", saved.Primary.ID)
		}

		foundVolume := res.VolumeSet[0]

		if nifcloud.StringValue(foundVolume.VolumeId) != saved.Primary.ID {
			return fmt.Errorf("volume does not found in cloud: %s", saved.Primary.ID)
		}

		*volume = foundVolume
		return nil
	}
}

func testAccCheckVolumeValues(volume *computing.VolumeSet, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.StringValue(volume.VolumeId) != rName {
			return fmt.Errorf("bad volume_id state, expected \"%s\", got: %#v", rName, volume.VolumeId)
		}

		if nifcloud.StringValue(volume.Size) != "100" {
			return fmt.Errorf("bad size state, expected \"100\", got: %#v", volume.Size)
		}

		if nifcloud.StringValue(volume.DiskType) != "High-Speed Storage A" {
			return fmt.Errorf("bad disk_type state,  expected \"High-Speed Storage A\", got: %#v", volume.DiskType)
		}

		if nifcloud.StringValue(volume.AttachmentSet[0].InstanceId) != rName {
			return fmt.Errorf("bad instance_id state, expected \"%s\", got: %#v", rName, volume.AttachmentSet[0].InstanceId)
		}

		if nifcloud.StringValue(volume.NextMonthAccountingType) != "1" {
			return fmt.Errorf("bad accounting_type state, expected \"1\", got: %#v", volume.NextMonthAccountingType)
		}

		if nifcloud.StringValue(volume.Description) != "memo" {
			return fmt.Errorf("bad description state,  expected \"memo\", got: %#v", volume.Description)
		}
		return nil
	}
}

func testAccCheckVolumeUniqueIDValues(volume *computing.VolumeSet, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.StringValue(volume.VolumeId) != rName {
			return fmt.Errorf("bad volume_id state, expected \"%s\", got: %#v", rName, volume.VolumeId)
		}

		if nifcloud.StringValue(volume.Size) != "100" {
			return fmt.Errorf("bad size state, expected \"100\", got: %#v", volume.Size)
		}

		if nifcloud.StringValue(volume.DiskType) != "High-Speed Storage A" {
			return fmt.Errorf("bad disk_type state,  expected \"High-Speed Storage A\", got: %#v", volume.DiskType)
		}

		if nifcloud.StringValue(volume.AttachmentSet[0].InstanceUniqueId) == "" {
			return fmt.Errorf("bad instance_unique_id state, expected not nil, got: nil")
		}

		if nifcloud.StringValue(volume.NextMonthAccountingType) != "1" {
			return fmt.Errorf("bad accounting_type state, expected \"1\", got: %#v", volume.NextMonthAccountingType)
		}

		if nifcloud.StringValue(volume.Description) != "memo" {
			return fmt.Errorf("bad description state,  expected \"memo\", got: %#v", volume.Description)
		}
		return nil
	}
}

func testAccCheckVolumeValuesUpdated(volume *computing.VolumeSet, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.StringValue(volume.VolumeId) != rName+"upd" {
			return fmt.Errorf("bad volume_id state, expected \"%s\", got: %#v", rName+"upd", volume.VolumeId)
		}

		if nifcloud.StringValue(volume.Size) != "300" {
			return fmt.Errorf("bad size state, expected \"300\", got: %#v", volume.Size)
		}

		if nifcloud.StringValue(volume.DiskType) != "High-Speed Storage A" {
			return fmt.Errorf("bad disk_type state,  expected \"High-Speed Storage A\", got: %#v", volume.DiskType)
		}

		if nifcloud.StringValue(volume.AttachmentSet[0].InstanceId) != rName {
			return fmt.Errorf("bad instance_id state, expected \"%s\", got: %#v", rName, volume.AttachmentSet[0].InstanceId)
		}

		if nifcloud.StringValue(volume.NextMonthAccountingType) != "2" {
			return fmt.Errorf("bad accounting_type state, expected \"2\", got: %#v", volume.NextMonthAccountingType)
		}

		if nifcloud.StringValue(volume.Description) != "memo-upd" {
			return fmt.Errorf("bad description state,  expected \"memo-upd\", got: %#v", volume.Description)
		}
		return nil
	}
}

func testAccVolumeResourceDestroy(s *terraform.State) error {
	svc := testAccProvider.Meta().(*client.Client).Computing

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nifcloud_volume" {
			continue
		}

		res, err := svc.DescribeVolumesRequest(&computing.DescribeVolumesInput{
			VolumeId: []string{rs.Primary.ID},
		}).Send(context.Background())

		if err != nil {
			var awsErr awserr.Error
			if errors.As(err, &awsErr) && awsErr.Code() == "Client.InvalidParameterNotFound.Volume" {
				return nil
			}
			return fmt.Errorf("failed DescribeVolumesRequest: %s", err)
		}

		if len(res.VolumeSet) > 0 {
			return fmt.Errorf("volume (%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testSweepVolume(region string) error {
	ctx := context.Background()
	svc := sharedClientForRegion(region).Computing

	res, err := svc.DescribeVolumesRequest(nil).Send(ctx)
	if err != nil {
		return err
	}

	var sweepVolumes []string
	for _, v := range res.VolumeSet {
		if strings.HasPrefix(nifcloud.StringValue(v.VolumeId), prefix) {
			sweepVolumes = append(sweepVolumes, nifcloud.StringValue(v.VolumeId))

			for _, a := range v.AttachmentSet {
				_, err = svc.DetachVolumeRequest(&computing.DetachVolumeInput{
					VolumeId:   nifcloud.String(nifcloud.StringValue(v.VolumeId)),
					InstanceId: nifcloud.String(nifcloud.StringValue(a.InstanceId)),
					Agreement:  nifcloud.Bool(true),
				}).Send(ctx)
				if err != nil {
					return err
				}

				err = svc.WaitUntilVolumeAvailable(ctx, &computing.DescribeVolumesInput{
					VolumeId: []string{nifcloud.StringValue(v.VolumeId)},
				})
			}

		}
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, n := range sweepVolumes {
		volumeID := n
		eg.Go(func() error {
			_, err = svc.DeleteVolumeRequest(&computing.DeleteVolumeInput{
				VolumeId: nifcloud.String(volumeID),
			}).Send(ctx)
			if err != nil {
				return err
			}

			err = svc.WaitUntilVolumeDeleted(ctx, &computing.DescribeVolumesInput{
				VolumeId: []string{volumeID},
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
