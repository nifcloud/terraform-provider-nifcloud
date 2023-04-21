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
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing/types"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
	"golang.org/x/sync/errgroup"
)

func init() {
	resource.AddTestSweepers("nifcloud_web_proxy", &resource.Sweeper{
		Name: "nifcloud_web_proxy",
		F:    testSweepWebProxy,
	})
}

func TestAcc_WebProxy(t *testing.T) {
	var webProxy types.WebProxyOfNiftyDescribeWebProxies

	resourceName := "nifcloud_web_proxy.basic"
	randName := prefix + acctest.RandString(7)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccWebProxyResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWebProxy(t, "testdata/web_proxy.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWebProxyExists(resourceName, &webProxy),
					testAccCheckWebProxyValues(&webProxy, randName),
					resource.TestCheckResourceAttr(resourceName, "description", "memo"),
					resource.TestCheckResourceAttr(resourceName, "name_server", "1.1.1.1"),
					resource.TestCheckResourceAttr(resourceName, "listen_port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "listen_interface_network_id", "net-COMMON_GLOBAL"),
					resource.TestCheckResourceAttr(resourceName, "bypass_interface_network_name", randName),
					resource.TestCheckResourceAttr(resourceName, "router_name", randName),
					resource.TestCheckResourceAttrSet(resourceName, "router_id"),
					resource.TestCheckResourceAttrSet(resourceName, "bypass_interface_network_id"),
				),
			},
			{
				Config: testAccWebProxy(t, "testdata/web_proxy_update.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWebProxyExists(resourceName, &webProxy),
					testAccCheckWebProxyValuesUpdated(&webProxy, randName),
					resource.TestCheckResourceAttr(resourceName, "description", "memo-upd"),
					resource.TestCheckResourceAttr(resourceName, "name_server", "8.8.8.8"),
					resource.TestCheckResourceAttr(resourceName, "listen_port", "3000"),
					resource.TestCheckResourceAttr(resourceName, "bypass_interface_network_id", "net-COMMON_GLOBAL"),
					resource.TestCheckResourceAttr(resourceName, "listen_interface_network_name", randName),
					resource.TestCheckResourceAttr(resourceName, "router_name", randName),
					resource.TestCheckResourceAttrSet(resourceName, "router_id"),
					resource.TestCheckResourceAttrSet(resourceName, "listen_interface_network_id"),
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

func testAccWebProxy(t *testing.T, fileName, rName string) string {
	b, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf(string(b),
		rName,
		rName,
	)
}

func testAccCheckWebProxyExists(n string, webProxy *types.WebProxyOfNiftyDescribeWebProxies) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no web proxy resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no web proxy id is set")
		}

		svc := testAccProvider.Meta().(*client.Client).Computing
		res, err := svc.NiftyDescribeWebProxies(context.Background(), &computing.NiftyDescribeWebProxiesInput{
			RouterId: []string{saved.Primary.ID},
		})
		if err != nil {
			return err
		}

		if res == nil || len(res.WebProxy) == 0 {
			return fmt.Errorf("web proxy does not found in cloud: %s", saved.Primary.ID)
		}

		foundWebProxy := res.WebProxy[0]

		if nifcloud.ToString(foundWebProxy.RouterId) != saved.Primary.ID {
			return fmt.Errorf("web proxy does not found in cloud: %s", saved.Primary.ID)
		}

		*webProxy = foundWebProxy

		return nil
	}
}

func testAccCheckWebProxyValues(webProxy *types.WebProxyOfNiftyDescribeWebProxies, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(webProxy.RouterName) != rName {
			return fmt.Errorf("bad router_name state,  expected \"%s\", got: %#v", rName, webProxy.RouterName)
		}

		if nifcloud.ToString(webProxy.BypassInterface.NetworkName) != rName {
			return fmt.Errorf("bad bypass_interface_network_name state,  expected \"%s\", got: %#v", rName, webProxy.BypassInterface.NetworkName)
		}

		if nifcloud.ToString(webProxy.BypassInterface.NetworkId) == "" {
			return fmt.Errorf("bad bypass_interface_network_id state ,  expected not empty string, got: %#v", webProxy.BypassInterface.NetworkId)
		}

		if nifcloud.ToString(webProxy.Description) != "memo" {
			return fmt.Errorf("bad description state,  expected \"memo\", got: %#v", webProxy.Description)
		}

		if nifcloud.ToString(webProxy.ListenPort) != "8080" {
			return fmt.Errorf("bad listen_port state,  expected \"8080\", got: %#v", webProxy.ListenPort)
		}

		if nifcloud.ToString(webProxy.Option.NameServer) != "1.1.1.1" {
			return fmt.Errorf("bad name_server state,  expected \"1.1.1.1\", got: %#v", webProxy.Option.NameServer)
		}

		if nifcloud.ToString(webProxy.ListenInterface.NetworkId) != "net-COMMON_GLOBAL" {
			return fmt.Errorf("bad listen_interface_network_id state,  expected \"net-COMMON_GLOBAL\", got: %#v", webProxy.ListenInterface.NetworkId)
		}

		if nifcloud.ToString(webProxy.RouterId) == "" {
			return fmt.Errorf("bad router_id,  expected not empty string, got: %#v", webProxy.RouterId)
		}
		return nil
	}
}

func testAccCheckWebProxyValuesUpdated(webProxy *types.WebProxyOfNiftyDescribeWebProxies, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(webProxy.RouterName) != rName {
			return fmt.Errorf("bad router_name state,  expected \"%s\", got: %#v", rName, webProxy.RouterName)
		}

		if nifcloud.ToString(webProxy.ListenInterface.NetworkName) != rName {
			return fmt.Errorf("bad listen_interface_network_name state,  expected \"%s\", got: %#v", rName, webProxy.ListenInterface.NetworkName)
		}

		if nifcloud.ToString(webProxy.ListenInterface.NetworkId) == "" {
			return fmt.Errorf("bad listen_interface_network_id state ,  expected not empty string, got: %#v", webProxy.ListenInterface.NetworkId)
		}

		if nifcloud.ToString(webProxy.Description) != "memo-upd" {
			return fmt.Errorf("bad description state,  expected \"memo-upd\", got: %#v", webProxy.Description)
		}

		if nifcloud.ToString(webProxy.ListenPort) != "3000" {
			return fmt.Errorf("bad listen_port state,  expected \"3000\", got: %#v", webProxy.ListenPort)
		}

		if nifcloud.ToString(webProxy.Option.NameServer) != "8.8.8.8" {
			return fmt.Errorf("bad name_server state,  expected \"8.8.8.8\", got: %#v", webProxy.Option.NameServer)
		}

		if nifcloud.ToString(webProxy.BypassInterface.NetworkId) != "net-COMMON_GLOBAL" {
			return fmt.Errorf("bad bypass_interface_network_id state,  expected \"net-COMMON_GLOBAL\", got: %#v", webProxy.BypassInterface.NetworkId)
		}

		if nifcloud.ToString(webProxy.RouterId) == "" {
			return fmt.Errorf("bad router_id,  expected not empty string, got: %#v", webProxy.RouterId)
		}

		return nil
	}
}

func testAccWebProxyResourceDestroy(s *terraform.State) error {
	svc := testAccProvider.Meta().(*client.Client).Computing

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nifcloud_web_proxy" {
			continue
		}

		res, err := svc.NiftyDescribeWebProxies(context.Background(), &computing.NiftyDescribeWebProxiesInput{
			RouterId: []string{rs.Primary.ID},
		})

		if err != nil {
			var awsErr smithy.APIError
			if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameterNotFound.RouterId" {
				return nil
			}
			return fmt.Errorf("failed listing web proxy: %s", err)
		}

		if len(res.WebProxy) > 0 {
			return fmt.Errorf("web proxy (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testSweepWebProxy(region string) error {
	ctx := context.Background()
	svc := sharedClientForRegion(region).Computing

	res, err := svc.NiftyDescribeWebProxies(ctx, nil)
	if err != nil {
		return err
	}

	var sweepWebProxies []string
	for _, w := range res.WebProxy {
		if strings.HasPrefix(nifcloud.ToString(w.RouterName), prefix) {
			sweepWebProxies = append(sweepWebProxies, nifcloud.ToString(w.RouterId))
		}
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, n := range sweepWebProxies {
		routerID := n
		eg.Go(func() error {
			_, err = svc.NiftyDeleteWebProxy(ctx, &computing.NiftyDeleteWebProxyInput{
				RouterId: nifcloud.String(routerID),
			})
			if err != nil {
				return err
			}

			err = computing.NewRouterAvailableWaiter(svc).Wait(ctx, &computing.NiftyDescribeRoutersInput{
				RouterId: []string{routerID},
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
