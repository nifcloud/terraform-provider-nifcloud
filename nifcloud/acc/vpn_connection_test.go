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
	resource.AddTestSweepers("nifcloud_vpn_connection", &resource.Sweeper{
		Name: "nifcloud_vpn_connection",
		F:    testSweepVpnConnection,
	})
}

func TestAcc_VpnConnection(t *testing.T) {
	var VpnConnection computing.VpnConnectionSet

	resourceName := "nifcloud_vpn_connection.basic"
	randName := prefix + acctest.RandStringFromCharSet(7, acctest.CharSetAlphaNum)

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
			},
		},
	})
}

func TestAcc_VpnConnection_Id_No_Tunnel(t *testing.T) {
	var VpnConnection computing.VpnConnectionSet

	resourceName := "nifcloud_vpn_connection.basic"
	randName := prefix + acctest.RandStringFromCharSet(7, acctest.CharSetAlphaNum)

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

func testAccCheckVpnConnectionExists(n string, VpnConnection *computing.VpnConnectionSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no vpn connection resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no vpn connection id is set")
		}

		svc := testAccProvider.Meta().(*client.Client).Computing
		res, err := svc.DescribeVpnConnectionsRequest(&computing.DescribeVpnConnectionsInput{
			VpnConnectionId: []string{saved.Primary.ID},
		}).Send(context.Background())

		if err != nil {
			return err
		}

		if len(res.VpnConnectionSet) == 0 {
			return fmt.Errorf("vpn connection does not found in cloud: %s", saved.Primary.ID)
		}

		foundVpnConnection := res.VpnConnectionSet[0]

		if nifcloud.StringValue(foundVpnConnection.VpnConnectionId) != saved.Primary.ID {
			return fmt.Errorf("vpn connection does not found in cloud: %s", saved.Primary.ID)
		}

		*VpnConnection = foundVpnConnection
		return nil
	}
}

func testAccCheckVpnConnectionValues(VpnConnection *computing.VpnConnectionSet, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.StringValue(VpnConnection.Type) != "L2TPv3 / IPsec" {
			return fmt.Errorf("bad type state, expected \"L2TPv3 / IPsec\", got: %#v", VpnConnection.Type)
		}

		if nifcloud.StringValue(VpnConnection.NiftyVpnGatewayName) != rName {
			return fmt.Errorf("bad vpn gateway name state, expected \"%s\", got: %#v", rName, VpnConnection.NiftyVpnGatewayName)
		}

		if nifcloud.StringValue(VpnConnection.NiftyCustomerGatewayName) != rName {
			return fmt.Errorf("bad customer gateway name state, expected \"%s\", got: %#v", rName, VpnConnection.NiftyCustomerGatewayName)
		}

		if nifcloud.StringValue(VpnConnection.NiftyTunnel.Type) != "L2TPv3" {
			return fmt.Errorf("bad tunnel type state, expected \"L2TPv3\", got: %#v", VpnConnection.NiftyTunnel.Type)
		}

		if nifcloud.StringValue(VpnConnection.NiftyTunnel.Mode) != "Unmanaged" {
			return fmt.Errorf("bad tunnel mode state, expected \"L2TPv3\", got: %#v", VpnConnection.NiftyTunnel.Mode)
		}

		if nifcloud.StringValue(VpnConnection.NiftyTunnel.Encapsulation) != "UDP" {
			return fmt.Errorf("bad tunnel encapsulation state, expected \"UDP\", got: %#v", VpnConnection.NiftyTunnel.Encapsulation)
		}

		if nifcloud.StringValue(VpnConnection.NiftyTunnel.TunnelId) != "1" {
			return fmt.Errorf("bad tunnel id state, expected \"1\", got: %#v", VpnConnection.NiftyTunnel.TunnelId)
		}

		if nifcloud.StringValue(VpnConnection.NiftyTunnel.PeerTunnelId) != "2" {
			return fmt.Errorf("bad tunnel peer id state, expected \"2\", got: %#v", VpnConnection.NiftyTunnel.PeerTunnelId)
		}

		if nifcloud.StringValue(VpnConnection.NiftyTunnel.SessionId) != "1" {
			return fmt.Errorf("bad tunnel session id state, expected \"1\", got: %#v", VpnConnection.NiftyTunnel.SessionId)
		}

		if nifcloud.StringValue(VpnConnection.NiftyTunnel.PeerSessionId) != "2" {
			return fmt.Errorf("bad tunnel peer session id state, expected \"2\", got: %#v", VpnConnection.NiftyTunnel.PeerSessionId)
		}

		if nifcloud.StringValue(VpnConnection.NiftyTunnel.SourcePort) != "7777" {
			return fmt.Errorf("bad tunnel source port state, expected \"7777\", got: %#v", VpnConnection.NiftyTunnel.SourcePort)
		}

		if nifcloud.StringValue(VpnConnection.NiftyTunnel.DestinationPort) != "7778" {
			return fmt.Errorf("bad tunnel destination port state, expected \"7778\", got: %#v", VpnConnection.NiftyTunnel.DestinationPort)
		}

		if nifcloud.StringValue(VpnConnection.NiftyIpsecConfiguration.Mtu) != "1000" {
			return fmt.Errorf("bad mtu state, expected \"1000\", got: %#v", VpnConnection.NiftyIpsecConfiguration.Mtu)
		}

		if nifcloud.StringValue(VpnConnection.NiftyIpsecConfiguration.EncryptionAlgorithm) != "AES256" {
			return fmt.Errorf("bad ipsec config encryption algorithm state, expected \"AES256\", got: %#v", VpnConnection.NiftyIpsecConfiguration.EncryptionAlgorithm)
		}

		if nifcloud.StringValue(VpnConnection.NiftyIpsecConfiguration.HashingAlgorithm) != "SHA256" {
			return fmt.Errorf("bad ipsec config hash algorithm state, expected \"SHA256\", got: %#v", VpnConnection.NiftyIpsecConfiguration.HashingAlgorithm)
		}

		if nifcloud.StringValue(VpnConnection.NiftyIpsecConfiguration.PreSharedKey) != "test" {
			return fmt.Errorf("bad ipsec config pre shared key state, expected \"test\", got: %#v", VpnConnection.NiftyIpsecConfiguration.PreSharedKey)
		}

		if nifcloud.StringValue(VpnConnection.NiftyIpsecConfiguration.InternetKeyExchange) != "IKEv2" {
			return fmt.Errorf("bad ipsec config internet key exchange state, expected \"IKEv2\", got: %#v", VpnConnection.NiftyIpsecConfiguration.InternetKeyExchange)
		}

		if nifcloud.Int64Value(VpnConnection.NiftyIpsecConfiguration.InternetKeyExchangeLifetime) != 300 {
			return fmt.Errorf("bad ipsec config internet key exchange lifetime state, expected \"300\", got: %#v", VpnConnection.NiftyIpsecConfiguration.InternetKeyExchangeLifetime)
		}

		if nifcloud.Int64Value(VpnConnection.NiftyIpsecConfiguration.EncapsulatingSecurityPayloadLifetime) != 301 {
			return fmt.Errorf("bad ipsec config encapsulating security payload lifetime state, expected \"301\", got: %#v", VpnConnection.NiftyIpsecConfiguration.EncapsulatingSecurityPayloadLifetime)
		}

		if nifcloud.Int64Value(VpnConnection.NiftyIpsecConfiguration.DiffieHellmanGroup) != 5 {
			return fmt.Errorf("bad ipsec config diffie hellman group state, expected \"5\", got: %#v", VpnConnection.NiftyIpsecConfiguration.DiffieHellmanGroup)
		}

		if nifcloud.StringValue(VpnConnection.NiftyVpnConnectionDescription) != "tfacc-memo" {
			return fmt.Errorf("bad description state, expected \"tfacc-memo\", got: %#v", VpnConnection.NiftyVpnConnectionDescription)
		}
		return nil
	}
}

func testAccCheckVpnConnectionValuesIDNoTunnel(VpnConnection *computing.VpnConnectionSet, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if nifcloud.StringValue(VpnConnection.Type) != "IPsec" {
			return fmt.Errorf("bad type state, expected \"IPsec\", got: %#v", VpnConnection.Type)
		}

		if nifcloud.StringValue(VpnConnection.VpnGatewayId) == "" {
			return fmt.Errorf("bad vpn gateway id state, expected \"%s\", got: %#v", rName, VpnConnection.VpnGatewayId)
		}

		if nifcloud.StringValue(VpnConnection.CustomerGatewayId) == "" {
			return fmt.Errorf("bad customer gateway id state, expected \"%s\", got: %#v", rName, VpnConnection.CustomerGatewayId)
		}

		if nifcloud.StringValue(VpnConnection.NiftyIpsecConfiguration.EncryptionAlgorithm) != "AES128" {
			return fmt.Errorf("bad ipsec config encryption algorithm state, expected \"AES128\", got: %#v", VpnConnection.NiftyIpsecConfiguration.EncryptionAlgorithm)
		}

		if nifcloud.StringValue(VpnConnection.NiftyIpsecConfiguration.HashingAlgorithm) != "SHA1" {
			return fmt.Errorf("bad ipsec config hash algorithm state, expected \"SHA1\", got: %#v", VpnConnection.NiftyIpsecConfiguration.HashingAlgorithm)
		}

		if nifcloud.StringValue(VpnConnection.NiftyIpsecConfiguration.PreSharedKey) != "test" {
			return fmt.Errorf("bad ipsec config pre shared key state, expected \"test\", got: %#v", VpnConnection.NiftyIpsecConfiguration.PreSharedKey)
		}

		if nifcloud.StringValue(VpnConnection.NiftyIpsecConfiguration.InternetKeyExchange) != "IKEv2" {
			return fmt.Errorf("bad ipsec config internet key exchange state, expected \"IKEv2\", got: %#v", VpnConnection.NiftyIpsecConfiguration.InternetKeyExchange)
		}

		if nifcloud.Int64Value(VpnConnection.NiftyIpsecConfiguration.InternetKeyExchangeLifetime) != 300 {
			return fmt.Errorf("bad ipsec config internet key exchange lifetime state, expected \"300\", got: %#v", VpnConnection.NiftyIpsecConfiguration.InternetKeyExchangeLifetime)
		}

		if nifcloud.Int64Value(VpnConnection.NiftyIpsecConfiguration.EncapsulatingSecurityPayloadLifetime) != 301 {
			return fmt.Errorf("bad ipsec config encapsulating security payload lifetime state, expected \"301\", got: %#v", VpnConnection.NiftyIpsecConfiguration.EncapsulatingSecurityPayloadLifetime)
		}

		if nifcloud.Int64Value(VpnConnection.NiftyIpsecConfiguration.DiffieHellmanGroup) != 2 {
			return fmt.Errorf("bad ipsec config diffie hellman group state, expected \"2\", got: %#v", VpnConnection.NiftyIpsecConfiguration.DiffieHellmanGroup)
		}

		if nifcloud.StringValue(VpnConnection.NiftyVpnConnectionDescription) != "tfacc-memo" {
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

		res, err := svc.DescribeVpnConnectionsRequest(&computing.DescribeVpnConnectionsInput{
			VpnConnectionId: []string{rs.Primary.ID},
		}).Send(context.Background())

		if err != nil {
			var awsErr awserr.Error
			if errors.As(err, &awsErr) && awsErr.Code() != "Client.InvalidParameterNotFound.VpnConnectionId" {
				return fmt.Errorf("failed DescribeVpnConnectionsRequest: %s", err)
			}
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

	res, err := svc.DescribeVpnConnectionsRequest(nil).Send(ctx)
	if err != nil {
		return err
	}

	var sweepVpnConnections []string
	for _, k := range res.VpnConnectionSet {
		if strings.HasPrefix(nifcloud.StringValue(k.NiftyVpnConnectionDescription), prefix) {
			sweepVpnConnections = append(sweepVpnConnections, nifcloud.StringValue(k.VpnConnectionId))
		}
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, n := range sweepVpnConnections {
		vpnConnectionID := n
		eg.Go(func() error {
			_, err := svc.DeleteVpnConnectionRequest(&computing.DeleteVpnConnectionInput{
				VpnConnectionId: nifcloud.String(vpnConnectionID),
				Agreement:       nifcloud.Bool(false),
			}).Send(ctx)
			if err != nil {
				return err
			}

			err = svc.WaitUntilVpnConnectionDeleted(ctx, &computing.DescribeVpnConnectionsInput{
				VpnConnectionId: []string{vpnConnectionID},
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
