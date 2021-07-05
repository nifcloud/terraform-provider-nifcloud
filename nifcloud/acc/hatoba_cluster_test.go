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
	"github.com/nifcloud/nifcloud-sdk-go/service/hatoba"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func init() {
	resource.AddTestSweepers("nifcloud_hatoba_cluster", &resource.Sweeper{
		Name:         "nifcloud_hatoba_cluster",
		F:            testSweepHatobaCluster,
		Dependencies: []string{},
	})
}

func TestAcc_HatobaCluster(t *testing.T) {
	var cluster hatoba.Cluster

	resourceName := "nifcloud_hatoba_cluster.basic"
	randName := prefix + acctest.RandString(7)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccHatobaClusterResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccHatobaCluster(t, "testdata/hatoba_cluster.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHatobaClusterExists(resourceName, &cluster),
					testAccCheckHatobaClusterValues(&cluster, randName),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "description", "memo"),
					resource.TestCheckResourceAttr(resourceName, "firewall_group", randName),
					resource.TestCheckResourceAttr(resourceName, "locations.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "locations.0", "east-21"),
					resource.TestCheckResourceAttr(resourceName, "addons_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "addons_config.0.http_load_balancing.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "addons_config.0.http_load_balancing.0.disabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "network_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "network_config.0.network_id", "net-COMMON_PRIVATE"),
					resource.TestCheckResourceAttr(resourceName, "node_pools.#", "2"),
				),
			},
			{
				Config: testAccHatobaCluster(t, "testdata/hatoba_cluster_update.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHatobaClusterExists(resourceName, &cluster),
					testAccCheckHatobaClusterValuesUpdated(&cluster, randName),
					resource.TestCheckResourceAttr(resourceName, "name", randName+"upd"),
					resource.TestCheckResourceAttr(resourceName, "description", "memo-upd"),
					resource.TestCheckResourceAttr(resourceName, "firewall_group", randName),
					resource.TestCheckResourceAttr(resourceName, "locations.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "locations.0", "east-21"),
					resource.TestCheckResourceAttr(resourceName, "locations.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "locations.0", "east-21"),
					resource.TestCheckResourceAttr(resourceName, "addons_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "addons_config.0.http_load_balancing.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "addons_config.0.http_load_balancing.0.disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "network_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "network_config.0.network_id", "net-COMMON_PRIVATE"),
					resource.TestCheckResourceAttr(resourceName, "node_pools.#", "2"),
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

func fetchDefaultKubernetesVersion(region string) (string, error) {
	svc := sharedClientForRegion(region).Hatoba
	res, err := svc.GetServerConfigRequest(nil).Send(context.Background())
	if err != nil {
		return "", err
	}

	return nifcloud.StringValue(res.ServerConfig.DefaultKubernetesVersion), nil
}

func testAccHatobaCluster(t *testing.T, fileName, rName string) string {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf(string(b),
		rName,
		rName,
	)
}

func testAccCheckHatobaClusterExists(n string, cluster *hatoba.Cluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no Hatoba cluster resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no Hatoba cluster id is set")
		}

		svc := testAccProvider.Meta().(*client.Client).Hatoba
		res, err := svc.GetClusterRequest(&hatoba.GetClusterInput{
			ClusterName: nifcloud.String(saved.Primary.ID),
		}).Send(context.Background())
		if err != nil {
			return err
		}

		foundCluster := res.Cluster

		if nifcloud.StringValue(foundCluster.Name) != saved.Primary.ID {
			return fmt.Errorf("Hatoba cluster does not found in cloud: %s", saved.Primary.ID)
		}

		*cluster = *foundCluster

		return nil
	}
}

func testAccCheckHatobaClusterValues(cluster *hatoba.Cluster, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		defaultKubernetesVersion, err := fetchDefaultKubernetesVersion("jp-east-2")
		if err != nil {
			return err
		}

		if nifcloud.StringValue(cluster.Name) != name {
			return fmt.Errorf("bad name state, expected \"%s\", got: %#v", name, nifcloud.StringValue(cluster.Name))
		}

		if nifcloud.StringValue(cluster.Description) != "memo" {
			return fmt.Errorf("bad description state, expected \"memo\", got: %#v", nifcloud.StringValue(cluster.Description))
		}

		if nifcloud.StringValue(cluster.KubernetesVersion) != defaultKubernetesVersion {
			return fmt.Errorf("bad kubernetes_version state, expected %q, got: %#v", defaultKubernetesVersion, nifcloud.StringValue(cluster.KubernetesVersion))
		}

		if len(cluster.Locations) != 1 {
			return fmt.Errorf("bad locations state, expected length is 1, got: %d", len(cluster.Locations))
		}

		if cluster.Locations[0] != "east-21" {
			return fmt.Errorf("bad locations state, expected \"east-21\", got: %#v", cluster.Locations[0])
		}

		if cluster.AddonsConfig == nil || cluster.AddonsConfig.HttpLoadBalancing == nil {
			return fmt.Errorf("bad addons_config state, response is nil")
		}

		if nifcloud.BoolValue(cluster.AddonsConfig.HttpLoadBalancing.Disabled) != true {
			return fmt.Errorf("bad http_load_balancing state, expected true, got: false")
		}

		if cluster.NetworkConfig == nil {
			return fmt.Errorf("bad network_config state, response is nil")
		}

		if nifcloud.StringValue(cluster.NetworkConfig.NetworkId) != "net-COMMON_PRIVATE" {
			return fmt.Errorf("bad network_id state, expected \"net-COMMON_PRIVATE\", got: %#v", nifcloud.StringValue(cluster.NetworkConfig.NetworkId))
		}

		wantNodePools := []hatoba.NodePool{
			{
				Name:         nifcloud.String("default"),
				InstanceType: nifcloud.String("medium"),
				NodeCount:    nifcloud.Int64(1),
			},
			{
				Name:         nifcloud.String("lowspec"),
				InstanceType: nifcloud.String("e-medium"),
				NodeCount:    nifcloud.Int64(1),
			},
		}

		if len(cluster.NodePools) != len(wantNodePools) {
			return fmt.Errorf("bad node_pools state, expected length %d, got %d", len(wantNodePools), len(cluster.NodePools))
		}

		foundNodePoolNames := []string{}

		for _, g := range cluster.NodePools {
			for _, w := range wantNodePools {
				if nifcloud.StringValue(g.Name) == nifcloud.StringValue(w.Name) &&
					nifcloud.StringValue(g.InstanceType) == nifcloud.StringValue(w.InstanceType) &&
					nifcloud.Int64Value(g.NodeCount) == nifcloud.Int64Value(w.NodeCount) {
					foundNodePoolNames = append(foundNodePoolNames, nifcloud.StringValue(g.Name))
					break
				}
			}
		}

		if len(foundNodePoolNames) != len(wantNodePools) {
			return fmt.Errorf("bad node_pools state, expected node pool not found in cloud. found: %#v", foundNodePoolNames)
		}

		return nil
	}
}

func testAccCheckHatobaClusterValuesUpdated(cluster *hatoba.Cluster, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		defaultKubernetesVersion, err := fetchDefaultKubernetesVersion("jp-east-2")
		if err != nil {
			return err
		}

		if nifcloud.StringValue(cluster.Name) != name+"upd" {
			return fmt.Errorf("bad name state, expected \"%s\", got: %#v", name+"upd", nifcloud.StringValue(cluster.Name))
		}

		if nifcloud.StringValue(cluster.Description) != "memo-upd" {
			return fmt.Errorf("bad description state, expected \"memo-upd\", got: %#v", nifcloud.StringValue(cluster.Description))
		}

		if nifcloud.StringValue(cluster.KubernetesVersion) != defaultKubernetesVersion {
			return fmt.Errorf("bad kubernetes_version state, expected %q, got: %#v", defaultKubernetesVersion, nifcloud.StringValue(cluster.KubernetesVersion))
		}

		if len(cluster.Locations) != 1 {
			return fmt.Errorf("bad locations state, expected length is 1, got: %d", len(cluster.Locations))
		}

		if cluster.Locations[0] != "east-21" {
			return fmt.Errorf("bad locations state, expected \"east-21\", got: %#v", cluster.Locations[0])
		}

		if cluster.AddonsConfig == nil || cluster.AddonsConfig.HttpLoadBalancing == nil {
			return fmt.Errorf("bad addons_config state, response is nil")
		}

		if nifcloud.BoolValue(cluster.AddonsConfig.HttpLoadBalancing.Disabled) != false {
			return fmt.Errorf("bad http_load_balancing state, expected false, got: true")
		}

		if cluster.NetworkConfig == nil {
			return fmt.Errorf("bad network_config state, response is nil")
		}

		if nifcloud.StringValue(cluster.NetworkConfig.NetworkId) != "net-COMMON_PRIVATE" {
			return fmt.Errorf("bad network_id state, expected \"net-COMMON_PRIVATE\", got: %#v", nifcloud.StringValue(cluster.NetworkConfig.NetworkId))
		}

		wantNodePools := []hatoba.NodePool{
			{
				Name:         nifcloud.String("default"),
				InstanceType: nifcloud.String("medium"),
				NodeCount:    nifcloud.Int64(3),
			},
			{
				Name:         nifcloud.String("highspec"),
				InstanceType: nifcloud.String("large"),
				NodeCount:    nifcloud.Int64(1),
			},
		}

		if len(cluster.NodePools) != len(wantNodePools) {
			return fmt.Errorf("bad node_pools state, expected length %d, got %d", len(wantNodePools), len(cluster.NodePools))
		}

		foundNodePoolNames := []string{}

		for _, g := range cluster.NodePools {
			for _, w := range wantNodePools {
				if nifcloud.StringValue(g.Name) == nifcloud.StringValue(w.Name) &&
					nifcloud.StringValue(g.InstanceType) == nifcloud.StringValue(w.InstanceType) &&
					nifcloud.Int64Value(g.NodeCount) == nifcloud.Int64Value(w.NodeCount) {
					foundNodePoolNames = append(foundNodePoolNames, nifcloud.StringValue(g.Name))
					break
				}
			}
		}

		if len(foundNodePoolNames) != len(wantNodePools) {
			return fmt.Errorf("bad node_pools state, expected node pool not found in cloud. found: %#v", foundNodePoolNames)
		}

		return nil
	}
}

func testAccHatobaClusterResourceDestroy(s *terraform.State) error {
	svc := testAccProvider.Meta().(*client.Client).Hatoba

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nifcloud_hatoba_cluster" {
			continue
		}

		_, err := svc.GetClusterRequest(&hatoba.GetClusterInput{
			ClusterName: nifcloud.String(rs.Primary.ID),
		}).Send(context.Background())

		if err != nil {
			var awsErr awserr.Error
			if errors.As(err, &awsErr) && awsErr.Code() != "Client.InvalidParameterNotFound.Cluster" {
				return fmt.Errorf("failed GetClusterRequest: %s", err)
			}
		}
	}
	return nil
}

func testSweepHatobaCluster(region string) error {
	ctx := context.Background()
	svc := sharedClientForRegion(region).Hatoba

	res, err := svc.ListClustersRequest(nil).Send(ctx)
	if err != nil {
		return err
	}

	var sweepHatobaClusters []string
	for _, k := range res.Clusters {
		if strings.HasPrefix(nifcloud.StringValue(k.Name), prefix) {
			sweepHatobaClusters = append(sweepHatobaClusters, nifcloud.StringValue(k.Name))
		}
	}

	if _, err := svc.DeleteClustersRequest(&hatoba.DeleteClustersInput{
		Names: nifcloud.String(strings.Join(sweepHatobaClusters, ",")),
	}).Send(ctx); err != nil {
		return err
	}

	return nil
}
