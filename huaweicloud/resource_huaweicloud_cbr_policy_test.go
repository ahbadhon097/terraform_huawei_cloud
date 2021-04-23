package huaweicloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/cbr/v3/policies"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccCBRV3Policy_basic(t *testing.T) {
	var asPolicy policies.Policy
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cbr_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckASV1PolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCBRV3Policy_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCBRV3PolicyExists(resourceName, &asPolicy),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "time_period", "20"),
					resource.TestCheckResourceAttr(resourceName, "backup_cycle.0.frequency", "WEEKLY"),
					resource.TestCheckResourceAttr(resourceName, "backup_cycle.0.days", "MO,TU"),
					resource.TestCheckResourceAttr(resourceName, "backup_cycle.0.execution_times.0", "06:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_cycle.0.execution_times.1", "18:00"),
				),
			},
			{
				Config: testCBRV3Policy_update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCBRV3PolicyExists(resourceName, &asPolicy),
					resource.TestCheckResourceAttr(resourceName, "name", rName+"-update"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "backup_quantity", "5"),
					resource.TestCheckResourceAttr(resourceName, "backup_cycle.0.frequency", "WEEKLY"),
					resource.TestCheckResourceAttr(resourceName, "backup_cycle.0.days", "SA,SU"),
					resource.TestCheckResourceAttr(resourceName, "backup_cycle.0.execution_times.0", "08:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_cycle.0.execution_times.1", "20:00"),
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

func TestAccCBRV3Policy_replication(t *testing.T) {
	var asPolicy policies.Policy
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cbr_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckDpsID(t)
			testAccPreCheckDestRegion(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckASV1PolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCBRV3Policy_replication(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCBRV3PolicyExists(resourceName, &asPolicy),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "replication"),
					resource.TestCheckResourceAttr(resourceName, "destination_project_id", HW_DEST_PROJECT_ID),
					resource.TestCheckResourceAttr(resourceName, "destination_region", HW_DEST_REGION),
					resource.TestCheckResourceAttr(resourceName, "time_period", "20"),
					resource.TestCheckResourceAttr(resourceName, "backup_cycle.0.frequency", "DAILY"),
					resource.TestCheckResourceAttr(resourceName, "backup_cycle.0.interval", "5"),
					resource.TestCheckResourceAttr(resourceName, "backup_cycle.0.execution_times.0", "06:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_cycle.0.execution_times.1", "18:00"),
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

func testAccCheckCBRV3PolicyDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	client, err := config.CbrV3Client(HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating Huaweicloud CBR client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_cbr_policy" {
			continue
		}

		_, err := policies.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("policy still exists")
		}
	}
	return nil
}

func testAccCheckCBRV3PolicyExists(n string, policy *policies.Policy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		client, err := config.CbrV3Client(HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating Huaweicloud CBR client: %s", err)
		}

		found, err := policies.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		log.Printf("[DEBUG] test found is: %#v", found)
		policy = found

		return nil
	}
}

func testCBRV3Policy_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cbr_policy" "test" {
  name            = "%s"
  protection_type = "backup"
  time_period     = 20

  backup_cycle {
    frequency       = "WEEKLY"
    days            = "MO,TU"
    execution_times = ["06:00", "18:00"]
  }
}
`, rName)
}

func testCBRV3Policy_update(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cbr_policy" "test" {
  name            = "%s-update"
  protection_type = "backup"
  backup_quantity = 5

  backup_cycle {
    frequency       = "WEEKLY"
    days            = "SA,SU"
    execution_times = ["08:00", "20:00"]
  }
}
`, rName)
}

func testCBRV3Policy_replication(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cbr_policy" "test" {
  name                   = "%s"
  protection_type        = "replication"
  destination_project_id = "%s"
  destination_region     = "%s"
  time_period            = 20

  backup_cycle {
    frequency       = "DAILY"
    interval        = 5
    execution_times = ["06:00", "18:00"]
  }
}
`, rName, HW_DEST_PROJECT_ID, HW_DEST_REGION)
}
