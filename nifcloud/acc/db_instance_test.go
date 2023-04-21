package acc

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aws/smithy-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/rdb"
	"github.com/nifcloud/nifcloud-sdk-go/service/rdb/types"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
	"golang.org/x/sync/errgroup"
)

func init() {
	resource.AddTestSweepers("nifcloud_db_instance", &resource.Sweeper{
		Name: "nifcloud_db_instance",
		F:    testSweepDBInstance,
	})
}

func TestAcc_DBInstance(t *testing.T) {
	var dbInstance, dbInstanceReplica, dbInstanceRestore types.DBInstances

	resourceName := "nifcloud_db_instance.basic"
	resourceNameReplica := "nifcloud_db_instance.replica"
	resourceNameRestore := "nifcloud_db_instance.restore"

	randName := prefix + acctest.RandString(7)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccDBInstanceResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDBInstance(t, "testdata/db_instance.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBInstanceExists(resourceName, &dbInstance),
					testAccCheckDBInstanceExists(resourceNameReplica, &dbInstanceReplica),
					testAccCheckDBInstanceExists(resourceNameRestore, &dbInstanceRestore),
					testAccCheckDBInstanceValues(&dbInstance, randName),
					resource.TestCheckResourceAttr(resourceName, "identifier", randName),
					resource.TestCheckResourceAttr(resourceName, "accounting_type", "1"),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "east-21"),
					resource.TestCheckResourceAttr(resourceName, "instance_class", "db.large"),
					resource.TestCheckResourceAttr(resourceName, "db_name", "baz"),
					resource.TestCheckResourceAttr(resourceName, "username", "for"),
					resource.TestCheckResourceAttr(resourceName, "password", "barbarbar"),
					resource.TestCheckResourceAttr(resourceName, "engine", "mysql"),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "5.7.15"),
					resource.TestCheckResourceAttr(resourceName, "allocated_storage", "50"),
					resource.TestCheckResourceAttr(resourceName, "storage_type", "0"),
					resource.TestCheckResourceAttr(resourceName, "backup_retention_period", "1"),
					resource.TestCheckResourceAttr(resourceName, "binlog_retention_period", "1"),
					resource.TestCheckResourceAttr(resourceName, "custom_binlog_retention_period", "true"),
					resource.TestCheckResourceAttr(resourceName, "backup_window", "00:00-08:00"),
					resource.TestCheckResourceAttr(resourceName, "maintenance_window", "sun:23:00-sun:23:30"),
					resource.TestCheckResourceAttr(resourceName, "multi_az", "true"),
					resource.TestCheckResourceAttr(resourceName, "multi_az_type", "0"),
					resource.TestCheckResourceAttr(resourceName, "port", "3306"),
					resource.TestCheckResourceAttr(resourceName, "publicly_accessible", "true"),
					resource.TestCheckResourceAttr(resourceName, "db_security_group_name", randName),
					resource.TestCheckResourceAttr(resourceName, "parameter_group_name", randName),
					resource.TestCheckResourceAttr(resourceNameReplica, "identifier", randName+"-replica"),
					resource.TestCheckResourceAttr(resourceNameReplica, "accounting_type", "1"),
					resource.TestCheckResourceAttr(resourceNameReplica, "instance_class", "db.large"),
					resource.TestCheckResourceAttr(resourceNameReplica, "storage_type", "0"),
					resource.TestCheckResourceAttr(resourceNameRestore, "identifier", randName+"-restore"),
					resource.TestCheckResourceAttr(resourceNameRestore, "accounting_type", "1"),
					resource.TestCheckResourceAttr(resourceNameRestore, "availability_zone", "east-21"),
					resource.TestCheckResourceAttr(resourceNameRestore, "instance_class", "db.large"),
					resource.TestCheckResourceAttr(resourceNameRestore, "storage_type", "0"),
					resource.TestCheckResourceAttr(resourceNameRestore, "multi_az", "true"),
					resource.TestCheckResourceAttr(resourceNameRestore, "multi_az_type", "0"),
					resource.TestCheckResourceAttr(resourceNameRestore, "port", "3306"),
					resource.TestCheckResourceAttr(resourceNameRestore, "publicly_accessible", "true"),
					resource.TestCheckResourceAttr(resourceNameRestore, "db_security_group_name", randName),
					resource.TestCheckResourceAttr(resourceNameRestore, "parameter_group_name", randName),
				),
			},
			{
				Config: testAccDBInstance(t, "testdata/db_instance_update.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBInstanceExists(resourceName, &dbInstance),
					testAccCheckDBInstanceExists(resourceNameReplica, &dbInstanceReplica),
					testAccCheckDBInstanceExists(resourceNameRestore, &dbInstanceRestore),
					testAccCheckDBInstanceValuesUpdated(&dbInstance, randName),
					resource.TestCheckResourceAttr(resourceName, "identifier", randName+"upd"),
					resource.TestCheckResourceAttr(resourceName, "accounting_type", "2"),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", "east-21"),
					resource.TestCheckResourceAttr(resourceName, "instance_class", "db.large8"),
					resource.TestCheckResourceAttr(resourceName, "db_name", "baz"),
					resource.TestCheckResourceAttr(resourceName, "username", "for"),
					resource.TestCheckResourceAttr(resourceName, "password", "barbarbarupd"),
					resource.TestCheckResourceAttr(resourceName, "engine", "mysql"),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "5.7.15"),
					resource.TestCheckResourceAttr(resourceName, "allocated_storage", "100"),
					resource.TestCheckResourceAttr(resourceName, "storage_type", "0"),
					resource.TestCheckResourceAttr(resourceName, "backup_retention_period", "2"),
					resource.TestCheckResourceAttr(resourceName, "binlog_retention_period", "2"),
					resource.TestCheckResourceAttr(resourceName, "custom_binlog_retention_period", "true"),
					resource.TestCheckResourceAttr(resourceName, "backup_window", "00:00-09:00"),
					resource.TestCheckResourceAttr(resourceName, "maintenance_window", "sun:22:00-sun:22:30"),
					resource.TestCheckResourceAttr(resourceName, "multi_az", "false"),
					resource.TestCheckResourceAttr(resourceName, "port", "3306"),
					resource.TestCheckResourceAttr(resourceName, "publicly_accessible", "true"),
					resource.TestCheckResourceAttr(resourceName, "db_security_group_name", randName+"upd"),
					resource.TestCheckResourceAttr(resourceName, "parameter_group_name", randName+"upd"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"apply_immediately",
					"final_snapshot_identifier",
					"password",
					"skip_final_snapshot",
					"multi_az_type",
				},
			},
		},
	})
}

func testAccDBInstance(t *testing.T, fileName, rName string) string {
	b, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf(string(b),
		rName,
		rName,
		rName,
		rName,
		rName,
		rName,
		rName,
		rName,
		rName,
	)
}

func testAccCheckDBInstanceExists(n string, dbInstance *types.DBInstances) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no db instance resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no db instance id is set")
		}

		svc := testAccProvider.Meta().(*client.Client).RDB
		res, err := svc.DescribeDBInstances(context.Background(), &rdb.DescribeDBInstancesInput{
			DBInstanceIdentifier: nifcloud.String(saved.Primary.ID),
		})

		if err != nil {
			return err
		}

		if res == nil || len(res.DBInstances) == 0 {
			return fmt.Errorf("db instance does not found in cloud: %s", saved.Primary.ID)
		}

		foundDBInstance := res.DBInstances[0]

		if nifcloud.ToString(foundDBInstance.DBInstanceIdentifier) != saved.Primary.ID {
			return fmt.Errorf("db instance does not found in cloud: %s", saved.Primary.ID)
		}

		*dbInstance = foundDBInstance
		return nil
	}
}

func testAccCheckDBInstanceValues(dbInstance *types.DBInstances, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(dbInstance.DBInstanceIdentifier) != rName {
			return fmt.Errorf("bad identifier state,  expected \"%s\", got: %#v", rName, dbInstance.DBInstanceIdentifier)
		}

		if nifcloud.ToString(dbInstance.NextMonthAccountingType) != "1" {
			return fmt.Errorf("bad accounting_type state,  expected \"1\", got: %#v", dbInstance.NextMonthAccountingType)
		}

		if nifcloud.ToString(dbInstance.AvailabilityZone) != "east-21" {
			return fmt.Errorf("bad availability_zone state,  expected \"east-21\", got: %#v", dbInstance.AvailabilityZone)
		}

		if nifcloud.ToString(dbInstance.DBInstanceClass) != "db.large" {
			return fmt.Errorf("bad instance_class state,  expected \"db.large\", got: %#v", dbInstance.DBInstanceClass)
		}

		if nifcloud.ToString(dbInstance.DBName) != "baz" {
			return fmt.Errorf("bad db_name state,  expected \"baz\", got: %#v", dbInstance.DBName)
		}

		if nifcloud.ToString(dbInstance.MasterUsername) != "for" {
			return fmt.Errorf("bad username state,  expected \"for\", got: %#v", dbInstance.MasterUsername)
		}

		if nifcloud.ToString(dbInstance.Engine) != "mysql" {
			return fmt.Errorf("bad engine state,  expected \"mysql\", got: %#v", dbInstance.Engine)
		}

		if nifcloud.ToString(dbInstance.EngineVersion) != "5.7.15" {
			return fmt.Errorf("bad engine_version state,  expected \"5.7.15\", got: %#v", dbInstance.EngineVersion)
		}

		if nifcloud.ToInt32(dbInstance.AllocatedStorage) != 50 {
			return fmt.Errorf("bad allocated_storage state,  expected \"50\", got: %#v", dbInstance.AllocatedStorage)
		}

		if nifcloud.ToInt32(dbInstance.BackupRetentionPeriod) != 1 {
			return fmt.Errorf("bad backup_retention_period state,  expected \"1\", got: %#v", dbInstance.BackupRetentionPeriod)
		}

		if nifcloud.ToInt32(dbInstance.BinlogRetentionPeriod) != 1 {
			return fmt.Errorf("bad binlog_retention_period state,  expected \"1\", got: %#v", dbInstance.BinlogRetentionPeriod)
		}

		if nifcloud.ToString(dbInstance.PreferredBackupWindow) != "00:00-08:00" {
			return fmt.Errorf("bad backup_window state,  expected \"00:00-08:00\", got: %#v", dbInstance.PreferredBackupWindow)
		}

		if nifcloud.ToString(dbInstance.PreferredMaintenanceWindow) != "sun:23:00-sun:23:30" {
			return fmt.Errorf("bad maintenance_window state,  expected \"sun:23:00-sun:23:30\", got: %#v", dbInstance.PreferredMaintenanceWindow)
		}

		if nifcloud.ToBool(dbInstance.MultiAZ) != true {
			return fmt.Errorf("bad multi_az state,  expected \"true\", got: %#v", dbInstance.MultiAZ)
		}

		if nifcloud.ToString(dbInstance.NiftyMultiAZType) != "0" {
			return fmt.Errorf("bad multi_az_type state,  expected \"0\", got: %#v", dbInstance.NiftyMultiAZType)
		}

		if nifcloud.ToInt32(dbInstance.Endpoint.Port) != 3306 {
			return fmt.Errorf("bad port state,  expected \"3306\", got: %#v", dbInstance.Endpoint.Port)
		}

		if nifcloud.ToBool(dbInstance.PubliclyAccessible) != true {
			return fmt.Errorf("bad publicly_accessible state,  expected \"true\", got: %#v", dbInstance.PubliclyAccessible)
		}

		if nifcloud.ToString(dbInstance.DBSecurityGroups[0].DBSecurityGroupName) != rName {
			return fmt.Errorf("bad db_security_group_name state,  expected \"%s\", got: %#v", rName, dbInstance.DBSecurityGroups[0].DBSecurityGroupName)
		}

		if nifcloud.ToString(dbInstance.DBParameterGroups[0].DBParameterGroupName) != rName {
			return fmt.Errorf("bad db_parameter_group_name state,  expected \"%s\", got: %#v", rName, dbInstance.DBParameterGroups[0].DBParameterGroupName)
		}
		return nil
	}
}

func testAccCheckDBInstanceValuesUpdated(dbInstance *types.DBInstances, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(dbInstance.DBInstanceIdentifier) != rName+"upd" {
			return fmt.Errorf("bad identifier state,  expected \"%s\", got: %#v", rName+"upd", dbInstance.DBInstanceIdentifier)
		}

		if nifcloud.ToString(dbInstance.NextMonthAccountingType) != "2" {
			return fmt.Errorf("bad accounting_type state,  expected \"2\", got: %#v", dbInstance.NextMonthAccountingType)
		}

		if nifcloud.ToString(dbInstance.AvailabilityZone) != "east-21" {
			return fmt.Errorf("bad availability_zone state,  expected \"east-21\", got: %#v", dbInstance.AvailabilityZone)
		}

		if nifcloud.ToString(dbInstance.DBInstanceClass) != "db.large8" {
			return fmt.Errorf("bad instance_class state,  expected \"db.large8\", got: %#v", dbInstance.DBInstanceClass)
		}

		if nifcloud.ToString(dbInstance.DBName) != "baz" {
			return fmt.Errorf("bad db_name state,  expected \"baz\", got: %#v", dbInstance.DBName)
		}

		if nifcloud.ToString(dbInstance.MasterUsername) != "for" {
			return fmt.Errorf("bad username state,  expected \"for\", got: %#v", dbInstance.MasterUsername)
		}

		if nifcloud.ToString(dbInstance.Engine) != "mysql" {
			return fmt.Errorf("bad engine state,  expected \"mysql\", got: %#v", dbInstance.Engine)
		}

		if nifcloud.ToString(dbInstance.EngineVersion) != "5.7.15" {
			return fmt.Errorf("bad engine_version state,  expected \"5.7.15\", got: %#v", dbInstance.EngineVersion)
		}

		if nifcloud.ToInt32(dbInstance.AllocatedStorage) != 100 {
			return fmt.Errorf("bad allocated_storage state,  expected \"100\", got: %#v", dbInstance.AllocatedStorage)
		}

		if nifcloud.ToInt32(dbInstance.BackupRetentionPeriod) != 2 {
			return fmt.Errorf("bad backup_retention_period state,  expected \"2\", got: %#v", dbInstance.BackupRetentionPeriod)
		}

		if nifcloud.ToInt32(dbInstance.BinlogRetentionPeriod) != 2 {
			return fmt.Errorf("bad binlog_retention_period state,  expected \"2\", got: %#v", dbInstance.BinlogRetentionPeriod)
		}

		if nifcloud.ToString(dbInstance.PreferredBackupWindow) != "00:00-09:00" {
			return fmt.Errorf("bad backup_window state,  expected \"00:00-09:00\", got: %#v", dbInstance.PreferredBackupWindow)
		}

		if nifcloud.ToString(dbInstance.PreferredMaintenanceWindow) != "sun:22:00-sun:22:30" {
			return fmt.Errorf("bad maintenance_window state,  expected \"sun:22:00-sun:22:30\", got: %#v", dbInstance.PreferredMaintenanceWindow)
		}

		if nifcloud.ToBool(dbInstance.MultiAZ) != false {
			return fmt.Errorf("bad multi_az state,  expected \"true\", got: %#v", dbInstance.MultiAZ)
		}

		if nifcloud.ToInt32(dbInstance.Endpoint.Port) != 3306 {
			return fmt.Errorf("bad port state,  expected \"3306\", got: %#v", dbInstance.Endpoint.Port)
		}

		if nifcloud.ToBool(dbInstance.PubliclyAccessible) != true {
			return fmt.Errorf("bad publicly_accessible state,  expected \"true\", got: %#v", dbInstance.PubliclyAccessible)
		}

		if nifcloud.ToString(dbInstance.DBSecurityGroups[0].DBSecurityGroupName) != rName+"upd" {
			return fmt.Errorf("bad db_security_group_name state,  expected \"%s\", got: %#v", rName+"upd", dbInstance.DBSecurityGroups[0].DBSecurityGroupName)
		}

		if nifcloud.ToString(dbInstance.DBParameterGroups[0].DBParameterGroupName) != rName+"upd" {
			return fmt.Errorf("bad db_parameter_group_name state,  expected \"%s\", got: %#v", rName+"upd", dbInstance.DBParameterGroups[0].DBParameterGroupName)
		}
		return nil
	}
}

func testAccDBInstanceResourceDestroy(s *terraform.State) error {
	svc := testAccProvider.Meta().(*client.Client).RDB

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nifcloud_db_instance" {
			continue
		}

		res, err := svc.DescribeDBInstances(context.Background(), &rdb.DescribeDBInstancesInput{
			DBInstanceIdentifier: nifcloud.String(rs.Primary.ID),
		})

		if err != nil {
			var awsErr smithy.APIError
			if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameterNotFound.DBInstance" {
				return nil
			}
			return fmt.Errorf("failed DescribeDBInstancesRequest: %s", err)
		}

		if len(res.DBInstances) > 0 {
			return fmt.Errorf("db instance (%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testSweepDBInstance(region string) error {
	ctx := context.Background()
	svc := sharedClientForRegion(region).RDB

	res, err := svc.DescribeDBInstances(ctx, nil)
	if err != nil {
		return err
	}

	var sweepDBInstances []string
	for _, i := range res.DBInstances {
		if strings.HasPrefix(nifcloud.ToString(i.DBInstanceIdentifier), prefix) {
			sweepDBInstances = append(sweepDBInstances, nifcloud.ToString(i.DBInstanceIdentifier))
		}
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, n := range sweepDBInstances {
		identifier := n
		eg.Go(func() error {
			_, err = svc.DeleteDBInstance(ctx, &rdb.DeleteDBInstanceInput{
				DBInstanceIdentifier: nifcloud.String(identifier),
				SkipFinalSnapshot:    nifcloud.Bool(true),
			})
			if err != nil {
				return err
			}

			err = rdb.NewDBInstanceDeletedWaiter(svc).Wait(ctx, &rdb.DescribeDBInstancesInput{
				DBInstanceIdentifier: nifcloud.String(identifier),
			}, 600*time.Second)
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
