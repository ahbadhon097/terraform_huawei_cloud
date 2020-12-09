package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccHuaweiCloudRdsFlavorV1DataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckDeprecated(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHuaweiCloudRdsFlavorV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsFlavorV1DataSourceID("data.huaweicloud_rds_flavors_v1.flavor"),
					resource.TestCheckResourceAttrSet(
						"data.huaweicloud_rds_flavors_v1.flavor", "name"),
					resource.TestCheckResourceAttrSet(
						"data.huaweicloud_rds_flavors_v1.flavor", "id"),
					resource.TestCheckResourceAttrSet(
						"data.huaweicloud_rds_flavors_v1.flavor", "speccode"),
				),
			},
		},
	})
}

func TestAccHuaweiCloudRdsFlavorV1DataSource_speccode(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckDeprecated(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHuaweiCloudRdsFlavorV1DataSource_speccode,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkV2DataSourceID("data.huaweicloud_rds_flavors_v1.flavor"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_rds_flavors_v1.flavor", "name", "OTC_PGCM_XLARGE"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_rds_flavors_v1.flavor", "speccode", "rds.pg.s1.xlarge"),
				),
			},
		},
	})
}

func testAccCheckRdsFlavorV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find rds data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Rds data source ID not set ")
		}

		return nil
	}
}

var testAccHuaweiCloudRdsFlavorV1DataSource_basic = fmt.Sprintf(`

data "huaweicloud_rds_flavors_v1" "flavor" {
    region = "%s"
	datastore_name = "PostgreSQL"
    datastore_version = "9.5.5"
}
`, HW_REGION_NAME)

var testAccHuaweiCloudRdsFlavorV1DataSource_speccode = fmt.Sprintf(`

data "huaweicloud_rds_flavors_v1" "flavor" {
    region = "%s"
	datastore_name = "PostgreSQL"
    datastore_version = "9.5.5"
    speccode = "rds.pg.s1.xlarge"
}
`, HW_REGION_NAME)
