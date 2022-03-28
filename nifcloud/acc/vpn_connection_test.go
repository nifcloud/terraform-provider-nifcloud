package acc

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
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
	resource.AddTestSweepers("nifcloud_vpn_connection", &resource.Sweeper{
		Name: "nifcloud_vpn_connection",
		F:    testSweepVpnConnection,
	})
}

func TestAcc_VpnConnection(t *testing.T) {
	var VpnConnection types.VpnConnectionSet

	resourceName := "nifcloud_vpn_connection.basic"
	randName := prefix + acctest.RandString(7)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccVpnConnectionResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpnConnection(t, "testdata/vpn_connection.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnConnectionExists(resourceName, &VpnConnection),
					testAccCheckVpnConnectionValues(&VpnConnection, randName),
					resource.TestCheckResourceAttr(resourceName, "type", "L2TPv3 / IPsec"),
					resource.TestCheckResourceAttr(resourceName, "vpn_gateway_name", randName),
					resource.TestCheckResourceAttr(resourceName, "customer_gateway_name", randName),
					resource.TestCheckResourceAttr(resourceName, "tunnel_type", "L2TPv3"),
					resource.TestCheckResourceAttr(resourceName, "tunnel_mode", "Unmanaged"),
					resource.TestCheckResourceAttr(resourceName, "tunnel_encapsulation", "UDP"),
					resource.TestCheckResourceAttr(resourceName, "tunnel_id", "1"),
					resource.TestCheckResourceAttr(resourceName, "tunnel_peer_id", "2"),
					resource.TestCheckResourceAttr(resourceName, "tunnel_session_id", "1"),
					resource.TestCheckResourceAttr(resourceName, "tunnel_peer_session_id", "2"),
					resource.TestCheckResourceAttr(resourceName, "tunnel_source_port", "7777"),
					resource.TestCheckResourceAttr(resourceName, "tunnel_destination_port", "7778"),
					resource.TestCheckResourceAttr(resourceName, "mtu", "1000"),
					resource.TestCheckResourceAttr(resourceName, "ipsec_config_encryption_algorithm", "AES256"),
					resource.TestCheckResourceAttr(resourceName, "ipsec_config_hash_algorithm", "SHA256"),
					resource.TestCheckResourceAttr(resourceName, "ipsec_config_pre_shared_key", "test"),
					resource.TestCheckResourceAttr(resourceName, "ipsec_config_internet_key_exchange", "IKEv2"),
					resource.TestCheckResourceAttr(resourceName, "ipsec_config_internet_key_exchange_lifetime", "300"),
					resource.TestCheckResourceAttr(resourceName, "ipsec_config_encapsulating_security_payload_lifetime", "301"),
					resource.TestCheckResourceAttr(resourceName, "ipsec_config_diffie_hellman_group", "5"),
					resource.TestCheckResourceAttr(resourceName, "description", "tfacc-memo"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"customer_gateway_name",
					"vpn_gateway_name",
					"customer_gateway_id",
					"vpn_gateway_id",
				},
			},
		},
	})
}

func TestAcc_VpnConnection_Id_No_Tunnel(t *testing.T) {
	var VpnConnection types.VpnConnectionSet

	resourceName := "nifcloud_vpn_connection.basic"
	randName := prefix + acctest.RandString(7)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccVpnConnectionResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpnConnection(t, "testdata/vpn_connection_id_no_tunnel.tf", randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnConnectionExists(resourceName, &VpnConnection),
					testAccCheckVpnConnectionValuesIDNoTunnel(&VpnConnection, randName),
					resource.TestCheckResourceAttr(resourceName, "type", "IPsec"),
					resource.TestCheckResourceAttrSet(resourceName, "vpn_gateway_id"),
					resource.TestCheckResourceAttrSet(resourceName, "customer_gateway_id"),
					resource.TestCheckResourceAttr(resourceName, "ipsec_config_encryption_algorithm", "AES128"),
					resource.TestCheckResourceAttr(resourceName, "ipsec_config_hash_algorithm", "SHA1"),
					resource.TestCheckResourceAttr(resourceName, "ipsec_config_pre_shared_key", "test"),
					resource.TestCheckResourceAttr(resourceName, "ipsec_config_internet_key_exchange", "IKEv2"),
					resource.TestCheckResourceAttr(resourceName, "ipsec_config_internet_key_exchange_lifetime", "300"),
					resource.TestCheckResourceAttr(resourceName, "ipsec_config_encapsulating_security_payload_lifetime", "301"),
					resource.TestCheckResourceAttr(resourceName, "ipsec_config_diffie_hellman_group", "2"),
					resource.TestCheckResourceAttr(resourceName, "description", "tfacc-memo"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"mtu",
				},
			},
		},
	})
}

func testAccVpnConnection(t *testing.T, fileName, rName string) string {
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

func testAccCheckVpnConnectionExists(n string, VpnConnection *types.VpnConnectionSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no vpn connection resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no vpn connection id is set")
		}

		svc := testAccProvider.Meta().(*client.Client).Computing
		res, err := svc.DescribeVpnConnections(context.Background(), &computing.DescribeVpnConnectionsInput{
			VpnConnectionId: []string{saved.Primary.ID},
		})

		if err != nil {
			return err
		}

		if len(res.VpnConnectionSet) == 0 {
			return fmt.Errorf("vpn connection does not found in cloud: %s", saved.Primary.ID)
		}

		foundVpnConnection := res.VpnConnectionSet[0]

		if nifcloud.ToString(foundVpnConnection.VpnConnectionId) != saved.Primary.ID {
			return fmt.Errorf("vpn connection does not found in cloud: %s", saved.Primary.ID)
		}

		*VpnConnection = foundVpnConnection
		return nil
	}
}

func testAccCheckVpnConnectionValues(VpnConnection *types.VpnConnectionSet, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(VpnConnection.Type) != "L2TPv3 / IPsec" {
			return fmt.Errorf("bad type state, expected \"L2TPv3 / IPsec\", got: %#v", VpnConnection.Type)
		}

		if nifcloud.ToString(VpnConnection.NiftyVpnGatewayName) != rName {
			return fmt.Errorf("bad vpn gateway name state, expected \"%s\", got: %#v", rName, VpnConnection.NiftyVpnGatewayName)
		}

		if nifcloud.ToString(VpnConnection.NiftyCustomerGatewayName) != rName {
			return fmt.Errorf("bad customer gateway name state, expected \"%s\", got: %#v", rName, VpnConnection.NiftyCustomerGatewayName)
		}

		if nifcloud.ToString(VpnConnection.NiftyTunnel.Type) != "L2TPv3" {
			return fmt.Errorf("bad tunnel type state, expected \"L2TPv3\", got: %#v", VpnConnection.NiftyTunnel.Type)
		}

		if nifcloud.ToString(VpnConnection.NiftyTunnel.Mode) != "Unmanaged" {
			return fmt.Errorf("bad tunnel mode state, expected \"L2TPv3\", got: %#v", VpnConnection.NiftyTunnel.Mode)
		}

		if nifcloud.ToString(VpnConnection.NiftyTunnel.Encapsulation) != "UDP" {
			return fmt.Errorf("bad tunnel encapsulation state, expected \"UDP\", got: %#v", VpnConnection.NiftyTunnel.Encapsulation)
		}

		if nifcloud.ToString(VpnConnection.NiftyTunnel.TunnelId) != "1" {
			return fmt.Errorf("bad tunnel id state, expected \"1\", got: %#v", VpnConnection.NiftyTunnel.TunnelId)
		}

		if nifcloud.ToString(VpnConnection.NiftyTunnel.PeerTunnelId) != "2" {
			return fmt.Errorf("bad tunnel peer id state, expected \"2\", got: %#v", VpnConnection.NiftyTunnel.PeerTunnelId)
		}

		if nifcloud.ToString(VpnConnection.NiftyTunnel.SessionId) != "1" {
			return fmt.Errorf("bad tunnel session id state, expected \"1\", got: %#v", VpnConnection.NiftyTunnel.SessionId)
		}

		if nifcloud.ToString(VpnConnection.NiftyTunnel.PeerSessionId) != "2" {
			return fmt.Errorf("bad tunnel peer session id state, expected \"2\", got: %#v", VpnConnection.NiftyTunnel.PeerSessionId)
		}

		if nifcloud.ToString(VpnConnection.NiftyTunnel.SourcePort) != "7777" {
			return fmt.Errorf("bad tunnel source port state, expected \"7777\", got: %#v", VpnConnection.NiftyTunnel.SourcePort)
		}

		if nifcloud.ToString(VpnConnection.NiftyTunnel.DestinationPort) != "7778" {
			return fmt.Errorf("bad tunnel destination port state, expected \"7778\", got: %#v", VpnConnection.NiftyTunnel.DestinationPort)
		}

		if nifcloud.ToString(VpnConnection.NiftyIpsecConfiguration.Mtu) != "1000" {
			return fmt.Errorf("bad mtu state, expected \"1000\", got: %#v", VpnConnection.NiftyIpsecConfiguration.Mtu)
		}

		if nifcloud.ToString(VpnConnection.NiftyIpsecConfiguration.EncryptionAlgorithm) != "AES256" {
			return fmt.Errorf("bad ipsec config encryption algorithm state, expected \"AES256\", got: %#v", VpnConnection.NiftyIpsecConfiguration.EncryptionAlgorithm)
		}

		if nifcloud.ToString(VpnConnection.NiftyIpsecConfiguration.HashingAlgorithm) != "SHA256" {
			return fmt.Errorf("bad ipsec config hash algorithm state, expected \"SHA256\", got: %#v", VpnConnection.NiftyIpsecConfiguration.HashingAlgorithm)
		}

		if nifcloud.ToString(VpnConnection.NiftyIpsecConfiguration.PreSharedKey) != "test" {
			return fmt.Errorf("bad ipsec config pre shared key state, expected \"test\", got: %#v", VpnConnection.NiftyIpsecConfiguration.PreSharedKey)
		}

		if nifcloud.ToString(VpnConnection.NiftyIpsecConfiguration.InternetKeyExchange) != "IKEv2" {
			return fmt.Errorf("bad ipsec config internet key exchange state, expected \"IKEv2\", got: %#v", VpnConnection.NiftyIpsecConfiguration.InternetKeyExchange)
		}

		if nifcloud.ToInt32(VpnConnection.NiftyIpsecConfiguration.InternetKeyExchangeLifetime) != 300 {
			return fmt.Errorf("bad ipsec config internet key exchange lifetime state, expected \"300\", got: %#v", VpnConnection.NiftyIpsecConfiguration.InternetKeyExchangeLifetime)
		}

		if nifcloud.ToInt32(VpnConnection.NiftyIpsecConfiguration.EncapsulatingSecurityPayloadLifetime) != 301 {
			return fmt.Errorf("bad ipsec config encapsulating security payload lifetime state, expected \"301\", got: %#v", VpnConnection.NiftyIpsecConfiguration.EncapsulatingSecurityPayloadLifetime)
		}

		if nifcloud.ToInt32(VpnConnection.NiftyIpsecConfiguration.DiffieHellmanGroup) != 5 {
			return fmt.Errorf("bad ipsec config diffie hellman group state, expected \"5\", got: %#v", VpnConnection.NiftyIpsecConfiguration.DiffieHellmanGroup)
		}

		if nifcloud.ToString(VpnConnection.NiftyVpnConnectionDescription) != "tfacc-memo" {
			return fmt.Errorf("bad description state, expected \"tfacc-memo\", got: %#v", VpnConnection.NiftyVpnConnectionDescription)
		}
		return nil
	}
}

func testAccCheckVpnConnectionValuesIDNoTunnel(VpnConnection *types.VpnConnectionSet, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.ToString(VpnConnection.Type) != "IPsec" {
			return fmt.Errorf("bad type state, expected \"IPsec\", got: %#v", VpnConnection.Type)
		}

		if nifcloud.ToString(VpnConnection.VpnGatewayId) == "" {
			return fmt.Errorf("bad vpn gateway id state, expected \"%s\", got: %#v", rName, VpnConnection.VpnGatewayId)
		}

		if nifcloud.ToString(VpnConnection.CustomerGatewayId) == "" {
			return fmt.Errorf("bad customer gateway id state, expected \"%s\", got: %#v", rName, VpnConnection.CustomerGatewayId)
		}

		if nifcloud.ToString(VpnConnection.NiftyIpsecConfiguration.EncryptionAlgorithm) != "AES128" {
			return fmt.Errorf("bad ipsec config encryption algorithm state, expected \"AES128\", got: %#v", VpnConnection.NiftyIpsecConfiguration.EncryptionAlgorithm)
		}

		if nifcloud.ToString(VpnConnection.NiftyIpsecConfiguration.HashingAlgorithm) != "SHA1" {
			return fmt.Errorf("bad ipsec config hash algorithm state, expected \"SHA1\", got: %#v", VpnConnection.NiftyIpsecConfiguration.HashingAlgorithm)
		}

		if nifcloud.ToString(VpnConnection.NiftyIpsecConfiguration.PreSharedKey) != "test" {
			return fmt.Errorf("bad ipsec config pre shared key state, expected \"test\", got: %#v", VpnConnection.NiftyIpsecConfiguration.PreSharedKey)
		}

		if nifcloud.ToString(VpnConnection.NiftyIpsecConfiguration.InternetKeyExchange) != "IKEv2" {
			return fmt.Errorf("bad ipsec config internet key exchange state, expected \"IKEv2\", got: %#v", VpnConnection.NiftyIpsecConfiguration.InternetKeyExchange)
		}

		if nifcloud.ToInt32(VpnConnection.NiftyIpsecConfiguration.InternetKeyExchangeLifetime) != 300 {
			return fmt.Errorf("bad ipsec config internet key exchange lifetime state, expected \"300\", got: %#v", VpnConnection.NiftyIpsecConfiguration.InternetKeyExchangeLifetime)
		}

		if nifcloud.ToInt32(VpnConnection.NiftyIpsecConfiguration.EncapsulatingSecurityPayloadLifetime) != 301 {
			return fmt.Errorf("bad ipsec config encapsulating security payload lifetime state, expected \"301\", got: %#v", VpnConnection.NiftyIpsecConfiguration.EncapsulatingSecurityPayloadLifetime)
		}

		if nifcloud.ToInt32(VpnConnection.NiftyIpsecConfiguration.DiffieHellmanGroup) != 2 {
			return fmt.Errorf("bad ipsec config diffie hellman group state, expected \"2\", got: %#v", VpnConnection.NiftyIpsecConfiguration.DiffieHellmanGroup)
		}

		if nifcloud.ToString(VpnConnection.NiftyVpnConnectionDescription) != "tfacc-memo" {
			return fmt.Errorf("bad description state, expected \"tfacc-memo\", got: %#v", VpnConnection.NiftyVpnConnectionDescription)
		}
		return nil
	}

}

func testAccVpnConnectionResourceDestroy(s *terraform.State) error {
	svc := testAccProvider.Meta().(*client.Client).Computing

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nifcloud_vpn_connection" {
			continue
		}

		res, err := svc.DescribeVpnConnections(context.Background(), &computing.DescribeVpnConnectionsInput{
			VpnConnectionId: []string{rs.Primary.ID},
		})

		if err != nil {
			var awsErr smithy.APIError
			if errors.As(err, &awsErr) && awsErr.ErrorCode() == "Client.InvalidParameterNotFound.VpnConnectionId" {
				return nil
			}
			return fmt.Errorf("failed DescribeVpnConnectionsRequest: %s", err)
		}

		if len(res.VpnConnectionSet) > 0 {
			return fmt.Errorf("vpn connection (%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testSweepVpnConnection(region string) error {
	ctx := context.Background()
	svc := sharedClientForRegion(region).Computing

	res, err := svc.DescribeVpnConnections(ctx, nil)
	if err != nil {
		return err
	}

	var sweepVpnConnections []string
	for _, k := range res.VpnConnectionSet {
		if strings.HasPrefix(nifcloud.ToString(k.NiftyVpnConnectionDescription), prefix) {
			sweepVpnConnections = append(sweepVpnConnections, nifcloud.ToString(k.VpnConnectionId))
		}
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, n := range sweepVpnConnections {
		vpnConnectionID := n
		eg.Go(func() error {
			_, err := svc.DeleteVpnConnection(ctx, &computing.DeleteVpnConnectionInput{
				VpnConnectionId: nifcloud.String(vpnConnectionID),
				Agreement:       nifcloud.Bool(false),
			})
			if err != nil {
				return err
			}

			err = computing.NewVpnConnectionDeletedWaiter(svc).Wait(ctx, &computing.DescribeVpnConnectionsInput{
				VpnConnectionId: []string{vpnConnectionID},
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
