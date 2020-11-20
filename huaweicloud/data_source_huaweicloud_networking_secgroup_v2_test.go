package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccHuaweiCloudNetworkingSecGroupV2DataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHuaweiCloudNetworkingSecGroupV2DataSource_group,
			},
			{
				Config: testAccHuaweiCloudNetworkingSecGroupV2DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSecGroupV2DataSourceID("data.huaweicloud_networking_secgroup_v2.secgroup_1"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_networking_secgroup_v2.secgroup_1", "name", "secgroup_1"),
				),
			},
		},
	})
}

func TestAccHuaweiCloudNetworkingSecGroupV2DataSource_secGroupID(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHuaweiCloudNetworkingSecGroupV2DataSource_group,
			},
			{
				Config: testAccHuaweiCloudNetworkingSecGroupV2DataSource_secGroupID,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSecGroupV2DataSourceID("data.huaweicloud_networking_secgroup_v2.secgroup_1"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_networking_secgroup_v2.secgroup_1", "name", "secgroup_1"),
				),
			},
		},
	})
}

func testAccCheckNetworkingSecGroupV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find security group data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Security group data source ID not set")
		}

		return nil
	}
}

const testAccHuaweiCloudNetworkingSecGroupV2DataSource_group = `
resource "huaweicloud_networking_secgroup_v2" "secgroup_1" {
        name        = "secgroup_1"
	description = "My neutron security group"
}
`

var testAccHuaweiCloudNetworkingSecGroupV2DataSource_basic = fmt.Sprintf(`
%s

data "huaweicloud_networking_secgroup_v2" "secgroup_1" {
	name = "${huaweicloud_networking_secgroup_v2.secgroup_1.name}"
}
`, testAccHuaweiCloudNetworkingSecGroupV2DataSource_group)

var testAccHuaweiCloudNetworkingSecGroupV2DataSource_secGroupID = fmt.Sprintf(`
%s

data "huaweicloud_networking_secgroup_v2" "secgroup_1" {
	secgroup_id = "${huaweicloud_networking_secgroup_v2.secgroup_1.id}"
}
`, testAccHuaweiCloudNetworkingSecGroupV2DataSource_group)
