package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk/openstack/iec/v1/security/groups"
)

func TestAccIecSecurityGroupResourceV1_basic(t *testing.T) {
	var group groups.RespSecurityGroupEntity
	resourceName := "huaweicloud_iec_security_group.my_group"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIecSecurityGroupV1Destory,
		Steps: []resource.TestStep{
			{
				Config: testAccIecSecurityGroupV1_Basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIecSecurityGroupV1Exists(resourceName, &group),
					testAccCheckIecSecurityGroupV1RuleCount(&group, 0),
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

func testAccCheckIecSecurityGroupV1RuleCount(group *groups.RespSecurityGroupEntity, count int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(group.SecurityGroupRules) == count {
			return nil
		}

		return fmt.Errorf("Unexpected number of rules in group %s. Expected %d, got %d",
			group.ID, count, len(group.SecurityGroupRules))
	}
}

func testAccCheckIecSecurityGroupV1Exists(n string, group *groups.RespSecurityGroupEntity) resource.TestCheckFunc {

	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID has been seted")
		}

		config := testAccProvider.Meta().(*Config)
		iecClient, err := config.IECV1Client(HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating Huaweicloud IEC client: %s", err)
		}

		found, err := groups.Get(iecClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("IEC Security group not found")
		}
		*group = *found
		return nil
	}
}

func testAccCheckIecSecurityGroupV1Destory(s *terraform.State) error {

	config := testAccProvider.Meta().(*Config)
	iecClient, err := config.IECV1Client(HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud IEC client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_iec_security_group" {
			continue
		}
		_, err := groups.Get(iecClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("IEC Security group still exists")
		}
	}

	return nil
}

func testAccIecSecurityGroupV1_Basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_iec_security_group" "my_group" {
  name = "my_test_group"
  description = "this is a test group"
}
`)
}
