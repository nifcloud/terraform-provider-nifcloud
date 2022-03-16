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
	"github.com/nifcloud/nifcloud-sdk-go/service/hatoba"
	"github.com/nifcloud/nifcloud-sdk-go/service/hatoba/types"
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
	var cluster types.Cluster

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
					resource.TestCheckResourceAttrSet(resourceName, "kube_config_raw"),
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
					resource.TestCheckResourceAttrSet(resourceName, "kube_config_raw"),
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
	res, err := svc.GetServerConfig(context.Background(), nil)
	if err != nil {
		return "", err
	}

	return nifcloud.ToString(res.ServerConfig.DefaultKubernetesVersion), nil
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

func testAccCheckHatobaClusterExists(n string, cluster *types.Cluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no types cluster resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no types cluster id is set")
		}

		svc := testAccProvider.Meta().(*client.Client).Hatoba
		res, err := svc.GetCluster(context.Background(), &hatoba.GetClusterInput{
			ClusterName: nifcloud.String(saved.Primary.ID),
		})
		if err != nil {
			return err
		}

		foundCluster := res.Cluster

		if nifcloud.ToString(foundCluster.Name) != saved.Primary.ID {
			return fmt.Errorf("Hatoba cluster does not found in cloud: %s", saved.Primary.ID)
		}

		*cluster = *foundCluster

		return nil
	}
}

func testAccCheckHatobaClusterValues(cluster *types.Cluster, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		defaultKubernetesVersion, err := fetchDefaultKubernetesVersion("jp-east-2")
		if err != nil {
			return err
		}

		if nifcloud.ToString(cluster.Name) != name {
			return fmt.Errorf("bad name state, expected \"%s\", got: %#v", name, nifcloud.ToString(cluster.Name))
		}

		if nifcloud.ToString(cluster.Description) != "memo" {
			return fmt.Errorf("bad description state, expected \"memo\", got: %#v", nifcloud.ToString(cluster.Description))
		}

		if nifcloud.ToString(cluster.KubernetesVersion) != defaultKubernetesVersion {
			return fmt.Errorf("bad kubernetes_version state, expected %q, got: %#v", defaultKubernetesVersion, nifcloud.ToString(cluster.KubernetesVersion))
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

		if nifcloud.ToBool(cluster.AddonsConfig.HttpLoadBalancing.Disabled) != true {
			return fmt.Errorf("bad http_load_balancing state, expected true, got: false")
		}

		if cluster.NetworkConfig == nil {
			return fmt.Errorf("bad network_config state, response is nil")
		}

		if nifcloud.ToString(cluster.NetworkConfig.NetworkId) != "net-COMMON_PRIVATE" {
			return fmt.Errorf("bad network_id state, expected \"net-COMMON_PRIVATE\", got: %#v", nifcloud.ToString(cluster.NetworkConfig.NetworkId))
		}

		wantNodePools := []types.NodePool{
			{
				Name:         nifcloud.String("default"),
				InstanceType: nifcloud.String("medium"),
				NodeCount:    nifcloud.Int32(1),
			},
			{
				Name:         nifcloud.String("lowspec"),
				InstanceType: nifcloud.String("e-medium"),
				NodeCount:    nifcloud.Int32(1),
			},
		}

		if len(cluster.NodePools) != len(wantNodePools) {
			return fmt.Errorf("bad node_pools state, expected length %d, got %d", len(wantNodePools), len(cluster.NodePools))
		}

		foundNodePoolNames := []string{}

		for _, g := range cluster.NodePools {
			for _, w := range wantNodePools {
				if nifcloud.ToString(g.Name) == nifcloud.ToString(w.Name) &&
					nifcloud.ToString(g.InstanceType) == nifcloud.ToString(w.InstanceType) &&
					nifcloud.ToInt32(g.NodeCount) == nifcloud.ToInt32(w.NodeCount) {
					foundNodePoolNames = append(foundNodePoolNames, nifcloud.ToString(g.Name))
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

func testAccCheckHatobaClusterValuesUpdated(cluster *types.Cluster, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		defaultKubernetesVersion, err := fetchDefaultKubernetesVersion("jp-east-2")
		if err != nil {
			return err
		}

		if nifcloud.ToString(cluster.Name) != name+"upd" {
			return fmt.Errorf("bad name state, expected \"%s\", got: %#v", name+"upd", nifcloud.ToString(cluster.Name))
		}

		if nifcloud.ToString(cluster.Description) != "memo-upd" {
			return fmt.Errorf("bad description state, expected \"memo-upd\", got: %#v", nifcloud.ToString(cluster.Description))
		}

		if nifcloud.ToString(cluster.KubernetesVersion) != defaultKubernetesVersion {
			return fmt.Errorf("bad kubernetes_version state, expected %q, got: %#v", defaultKubernetesVersion, nifcloud.ToString(cluster.KubernetesVersion))
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

		if nifcloud.ToBool(cluster.AddonsConfig.HttpLoadBalancing.Disabled) != false {
			return fmt.Errorf("bad http_load_balancing state, expected false, got: true")
		}

		if cluster.NetworkConfig == nil {
			return fmt.Errorf("bad network_config state, response is nil")
		}

		if nifcloud.ToString(cluster.NetworkConfig.NetworkId) != "net-COMMON_PRIVATE" {
			return fmt.Errorf("bad network_id state, expected \"net-COMMON_PRIVATE\", got: %#v", nifcloud.ToString(cluster.NetworkConfig.NetworkId))
		}

		wantNodePools := []types.NodePool{
			{
				Name:         nifcloud.String("default"),
				InstanceType: nifcloud.String("medium"),
				NodeCount:    nifcloud.Int32(3),
			},
			{
				Name:         nifcloud.String("highspec"),
				InstanceType: nifcloud.String("large"),
				NodeCount:    nifcloud.Int32(1),
			},
		}

		if len(cluster.NodePools) != len(wantNodePools) {
			return fmt.Errorf("bad node_pools state, expected length %d, got %d", len(wantNodePools), len(cluster.NodePools))
		}

		foundNodePoolNames := []string{}

		for _, g := range cluster.NodePools {
			for _, w := range wantNodePools {
				if nifcloud.ToString(g.Name) == nifcloud.ToString(w.Name) &&
					nifcloud.ToString(g.InstanceType) == nifcloud.ToString(w.InstanceType) &&
					nifcloud.ToInt32(g.NodeCount) == nifcloud.ToInt32(w.NodeCount) {
					foundNodePoolNames = append(foundNodePoolNames, nifcloud.ToString(g.Name))
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

		_, err := svc.GetCluster(context.Background(), &hatoba.GetClusterInput{
			ClusterName: nifcloud.String(rs.Primary.ID),
		})

		if err != nil {
			var awsErr smithy.APIError
			if errors.As(err, &awsErr) && awsErr.ErrorCode() != "Client.InvalidParameterNotFound.Cluster" {
				return fmt.Errorf("failed GetClusterRequest: %s", err)
			}
		}
	}
	return nil
}

func testSweepHatobaCluster(region string) error {
	ctx := context.Background()
	svc := sharedClientForRegion(region).Hatoba

	res, err := svc.ListClusters(ctx, nil)
	if err != nil {
		return err
	}

	var sweepHatobaClusters []string
	for _, k := range res.Clusters {
		if strings.HasPrefix(nifcloud.ToString(k.Name), prefix) {
			sweepHatobaClusters = append(sweepHatobaClusters, nifcloud.ToString(k.Name))
		}
	}

	if len(sweepHatobaClusters) > 0 {
		if _, err := svc.DeleteClusters(ctx, &hatoba.DeleteClustersInput{
			Names: nifcloud.String(strings.Join(sweepHatobaClusters, ",")),
		}); err != nil {
			return err
		}
	}

	return nil
}
