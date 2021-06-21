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
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/acc/helper"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
	"golang.org/x/sync/errgroup"
)

func init() {
	resource.AddTestSweepers("nifcloud_ssl_certificate", &resource.Sweeper{
		Name: "nifcloud_ssl_certificate",
		F:    testSweepSSLCertificate,
	})
}

func TestAcc_SSLCertificate(t *testing.T) {
	var sslCertificate computing.CertsSet

	resourceName := "nifcloud_ssl_certificate.basic"
	randName := prefix + acctest.RandString(10)

	caKey := helper.GeneratePrivateKey(t, 4096)
	caCert := helper.GenerateSelfSignedCertificateAuthority(t, caKey)
	key := helper.GeneratePrivateKey(t, 4096)
	cert := helper.GenerateCertificate(t, caKey, caCert, key, randName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccSSLCertificateResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSSLCertificate(t, "testdata/ssl_certificate.tf", cert, key, caCert),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSSLCertificateExists(resourceName, &sslCertificate),
					testAccCheckSSLCertificateValues(resourceName, &sslCertificate, randName, cert, key, caCert),
					resource.TestCheckResourceAttrSet(resourceName, "fqdn_id"),
					resource.TestCheckResourceAttr(resourceName, "fqdn", randName),
					resource.TestCheckResourceAttr(resourceName, "description", "memo"),
					resource.TestCheckResourceAttr(resourceName, "certificate", cert),
					resource.TestCheckResourceAttr(resourceName, "key", key),
					resource.TestCheckResourceAttr(resourceName, "ca", caCert),
				),
			},
			{
				Config: testAccSSLCertificate(t, "testdata/ssl_certificate_update.tf", cert, key, caCert),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSSLCertificateExists(resourceName, &sslCertificate),
					testAccCheckSSLCertificateValuesUpdated(resourceName, &sslCertificate, randName, cert, key, caCert),
					resource.TestCheckResourceAttrSet(resourceName, "fqdn_id"),
					resource.TestCheckResourceAttr(resourceName, "fqdn", randName),
					resource.TestCheckResourceAttr(resourceName, "description", "memo-upd"),
					resource.TestCheckResourceAttr(resourceName, "certificate", cert),
					resource.TestCheckResourceAttr(resourceName, "key", key),
					resource.TestCheckResourceAttr(resourceName, "ca", caCert),
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

func testAccSSLCertificate(t *testing.T, fileName, certificate, key, ca string) string {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf(string(b), certificate, key, ca)
}

func testAccCheckSSLCertificateExists(n string, cert *computing.CertsSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no SSLCertificate resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no SSLCertificate id is set")
		}

		svc := testAccProvider.Meta().(*client.Client).Computing
		res, err := svc.DescribeSslCertificatesRequest(&computing.DescribeSslCertificatesInput{
			FqdnId: []string{saved.Primary.ID},
		}).Send(context.Background())

		if err != nil {
			return err
		}

		if len(res.CertsSet) == 0 {
			return fmt.Errorf("SSLCertificate does not found in cloud: %s", saved.Primary.ID)
		}

		foundCert := res.CertsSet[0]

		if nifcloud.StringValue(foundCert.FqdnId) != saved.Primary.ID {
			return fmt.Errorf("SSLCertificate does not found in cloud: %s", saved.Primary.ID)
		}

		*cert = foundCert

		return nil
	}
}

func testAccCheckSSLCertificateValues(
	n string, certSet *computing.CertsSet, fqdn, cert, key, caCert string,
) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no SSLCertificate resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no SSLCertificate id is set")
		}
		fqdnID := saved.Primary.ID

		if nifcloud.StringValue(certSet.FqdnId) != fqdnID {
			return fmt.Errorf("bad fqdn_id state, expected \"%s\", got: %#v", fqdnID, nifcloud.StringValue(certSet.FqdnId))
		}

		if nifcloud.StringValue(certSet.Fqdn) != fqdn {
			return fmt.Errorf("bad fqdn state, expected \"%s\", got: %#v", fqdn, nifcloud.StringValue(certSet.Fqdn))
		}

		if nifcloud.StringValue(certSet.Description) != "memo" {
			return fmt.Errorf("bad description state, expected \"memo\", got: %#v", nifcloud.StringValue(certSet.Description))
		}

		return testAccCheckSSLCertificateFileData(fqdnID, cert, key, caCert)
	}
}

func testAccCheckSSLCertificateValuesUpdated(
	n string, certsSet *computing.CertsSet, fqdn, cert, key, caCert string,
) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		saved, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("no SSLCertificate resource: %s", n)
		}

		if saved.Primary.ID == "" {
			return fmt.Errorf("no SSLCertificate id is set")
		}
		fqdnID := saved.Primary.ID

		if nifcloud.StringValue(certsSet.FqdnId) != fqdnID {
			return fmt.Errorf("bad fqdn_id state, expected \"%s\", got: %#v", fqdnID, nifcloud.StringValue(certsSet.FqdnId))
		}

		if nifcloud.StringValue(certsSet.Fqdn) != fqdn {
			return fmt.Errorf("bad fqdn state, expected \"%s\", got: %#v", fqdn, nifcloud.StringValue(certsSet.Fqdn))
		}

		if nifcloud.StringValue(certsSet.Description) != "memo-upd" {
			return fmt.Errorf("bad description state, expected \"memo-upd\", got: %#v", nifcloud.StringValue(certsSet.Description))
		}

		return testAccCheckSSLCertificateFileData(fqdnID, cert, key, caCert)
	}
}

func testAccCheckSSLCertificateFileData(fqdnID, wantCert, wantKey, wantCACert string) error {
	svc := testAccProvider.Meta().(*client.Client).Computing

	ctx := context.Background()
	res, err := svc.DownloadSslCertificateRequest(&computing.DownloadSslCertificateInput{
		FqdnId:   nifcloud.String(fqdnID),
		FileType: computing.FileTypeOfDownloadSslCertificateRequest1,
	}).Send(ctx)
	if err != nil {
		return err
	}

	gotKey := nifcloud.StringValue(res.DownloadSslCertificateOutput.FileData)
	if gotKey != wantKey {
		return fmt.Errorf("bad private key, expected \"%s\", got \"%s\"", wantKey, gotKey)
	}

	res, err = svc.DownloadSslCertificateRequest(&computing.DownloadSslCertificateInput{
		FqdnId:   nifcloud.String(fqdnID),
		FileType: computing.FileTypeOfDownloadSslCertificateRequest2,
	}).Send(ctx)
	if err != nil {
		return err
	}

	gotCACert := nifcloud.StringValue(res.DownloadSslCertificateOutput.FileData)
	if gotCACert != wantCACert {
		return fmt.Errorf("bad ca, expected \"%s\", got %#v", wantCACert, gotCACert)
	}

	res, err = svc.DownloadSslCertificateRequest(&computing.DownloadSslCertificateInput{
		FqdnId:   nifcloud.String(fqdnID),
		FileType: computing.FileTypeOfDownloadSslCertificateRequest3,
	}).Send(ctx)
	if err != nil {
		return err
	}

	gotCert := nifcloud.StringValue(res.DownloadSslCertificateOutput.FileData)
	if gotCert != wantCert {
		return fmt.Errorf("bad cert, expected \"%s\", got \"%s\"", wantCert, gotCert)
	}

	return nil
}

func testAccSSLCertificateResourceDestroy(s *terraform.State) error {
	svc := testAccProvider.Meta().(*client.Client).Computing

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nifcloud_ssl_certificate" {
			continue
		}

		res, err := svc.DescribeSslCertificatesRequest(&computing.DescribeSslCertificatesInput{
			FqdnId: []string{rs.Primary.ID},
		}).Send(context.Background())

		if err != nil {
			var awsErr awserr.Error
			if errors.As(err, &awsErr) && awsErr.Code() == "Client.InvalidParameterNotFound.SslCertificate" {
				return nil
			}
			return fmt.Errorf("failed DesribeSslCertificates: %s", err)
		}

		if len(res.CertsSet) > 0 {
			return fmt.Errorf("SSLCertificate (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testSweepSSLCertificate(region string) error {
	ctx := context.Background()
	svc := sharedClientForRegion(region).Computing

	res, err := svc.DescribeSslCertificatesRequest(nil).Send(ctx)
	if err != nil {
		return err
	}

	var sweepSSLCertificates []string
	for _, k := range res.CertsSet {
		if strings.HasPrefix(nifcloud.StringValue(k.FqdnId), prefix) {
			sweepSSLCertificates = append(sweepSSLCertificates, nifcloud.StringValue(k.FqdnId))
		}
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, id := range sweepSSLCertificates {
		fqdnID := id
		eg.Go(func() error {
			_, err := svc.DeleteSslCertificateRequest(&computing.DeleteSslCertificateInput{
				FqdnId: nifcloud.String(fqdnID),
			}).Send(ctx)
			return err
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}
