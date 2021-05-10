package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk/openstack/elb/v2/loadbalancers"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccELBV2LoadbalancerDataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccELBV2LoadbalancerDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckELBV2LoadbalancerDataSourceID("data.huaweicloud_lb_loadbalancer.test_by_name"),
					testAccCheckELBV2LoadbalancerDataSourceID("data.huaweicloud_lb_loadbalancer.test_by_description"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_lb_loadbalancer.test_by_name", "name", rName),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_lb_loadbalancer.test_by_description", "name", rName),
				),
			},
		},
	})
}

func testAccCheckELBV2LoadbalancerDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find elb load balancer data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("load balancer data source ID not set")
		}

		return nil
	}
}

func testAccCheckELBV2LoadbalancerDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	client, err := config.LoadBalancerClient(HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud load balancer client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_lb_loadbalancer" {
			continue
		}

		lb, err := loadbalancers.Get(client, rs.Primary.ID).Extract()
		if err == nil || lb.ID != "" {
			return fmt.Errorf("Load balancer still exists")
		}
	}

	return nil
}

func testAccELBV2LoadbalancerDataSource_basic(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "huaweicloud_lb_loadbalancer" "test" {
  name          = "%s"
  vip_subnet_id = data.huaweicloud_vpc_subnet.test.subnet_id
  description   = "test for load balancer data source"
}

data "huaweicloud_lb_loadbalancer" "test_by_name" {
  name = huaweicloud_lb_loadbalancer.test.name
}

data "huaweicloud_lb_loadbalancer" "test_by_description" {
  description = huaweicloud_lb_loadbalancer.test.description
}
`, rName)
}
