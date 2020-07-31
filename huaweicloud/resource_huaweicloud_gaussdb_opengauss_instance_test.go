package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk/openstack/opengauss/v3/instances"
)

func TestOpenGaussInstance_basic(t *testing.T) {
	var instance instances.GaussDBInstance

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_gaussdb_opengauss_instance.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOpenGaussInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenGaussInstanceConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOpenGaussInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
		},
	})
}

func testAccCheckOpenGaussInstanceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.initServiceClient("gaussdb", OS_REGION_NAME, "opengauss/v3")
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud GaussDB client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_gaussdb_opengauss_instance" {
			continue
		}

		v, err := instances.GetInstanceByID(client, rs.Primary.ID)
		if err == nil && v.Id == rs.Primary.ID {
			return fmt.Errorf("Instance <%s> still exists.", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckOpenGaussInstanceExists(n string, instance *instances.GaussDBInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s.", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set.")
		}

		config := testAccProvider.Meta().(*Config)
		client, err := config.initServiceClient("gaussdb", OS_REGION_NAME, "opengauss/v3")
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud GaussDB client: %s", err)
		}

		found, err := instances.GetInstanceByID(client, rs.Primary.ID)
		if err != nil {
			return err
		}
		if found.Id != rs.Primary.ID {
			return fmt.Errorf("Instance <%s> not found.", rs.Primary.ID)
		}
		instance = &found

		return nil
	}
}

func testAccOpenGaussInstanceConfig_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_networking_secgroup_v2" "test" {
  name = "default"
}

resource "huaweicloud_gaussdb_opengauss_instance" "test" {
  name        = "%s"
  password    = "Test@123"
  flavor      = "gaussdb.opengauss.ee.dn.m6.2xlarge.8.in"
  vpc_id      = huaweicloud_vpc_v1.test.id
  subnet_id   = huaweicloud_vpc_subnet_v1.test.id

  availability_zone = "cn-north-4a,cn-north-4a,cn-north-4a"
  security_group_id = data.huaweicloud_networking_secgroup_v2.test.id

  ha {
    mode             = "enterprise"
    replication_mode = "sync"
    consistency      = "strong"
  }
  volume {
    type = "ULTRAHIGH"
    size = 40
  }

  sharding_num = 1
  coordinator_num = 2
}
`, testAccVpcConfig_Base(rName), rName)
}
