package acc

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDatasourceImage_basic(t *testing.T) {
	datasourceName := "data.nifcloud_image.basic"

	//lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccImageDataSource(t, "testdata/data_image.tf"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImageDataSourceID(datasourceName),
					resource.TestCheckResourceAttr(datasourceName, "image_id", "283"),
					resource.TestCheckResourceAttr(datasourceName, "image_name", "Ubuntu Server 22.04 LTS"),
					resource.TestCheckResourceAttr(datasourceName, "owner", "niftycloud"),
				),
			},
		},
	})
}

func testAccImageDataSource(t *testing.T, fileName string) string {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func testAccCheckImageDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find image data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("image data source ID not set")
		}
		return nil
	}
}
