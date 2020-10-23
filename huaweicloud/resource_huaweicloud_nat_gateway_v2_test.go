package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/layer3/routers"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/natgateways"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/networks"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/subnets"
)

func TestAccNatGateway_basic(t *testing.T) {
	var network networks.Network
	var router routers.Router
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckNat(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNatV2GatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNatV2Gateway_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkExists("huaweicloud_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2SubnetExists("huaweicloud_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2RouterExists("huaweicloud_networking_router_v2.router_1", &router),
					testAccCheckNetworkingV2RouterInterfaceExists("huaweicloud_networking_router_interface_v2.int_1"),
					testAccCheckNatV2GatewayExists("huaweicloud_nat_gateway_v2.nat_1"),
				),
			},
			{
				Config: testAccNatV2Gateway_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("huaweicloud_nat_gateway_v2.nat_1", "name", "nat_1_updated"),
					resource.TestCheckResourceAttr("huaweicloud_nat_gateway_v2.nat_1", "description", "nat_1 updated"),
					resource.TestCheckResourceAttr("huaweicloud_nat_gateway_v2.nat_1", "spec", "2"),
				),
			},
		},
	})
}

func TestAccNatGateway_withEpsId(t *testing.T) {
	natgateway := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckNat(t); testAccPreCheckEpsID(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNatV2GatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNatV2Gateway_epsId(natgateway),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatV2GatewayExists("huaweicloud_nat_gateway.nat_1"),
					resource.TestCheckResourceAttr("huaweicloud_nat_gateway.nat_1", "enterprise_project_id", OS_ENTERPRISE_PROJECT_ID),
				),
			},
		},
	})
}

func testAccCheckNatV2GatewayDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	natClient, err := config.natGatewayV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud nat client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_nat_gateway_v2" {
			continue
		}

		_, err := natgateways.Get(natClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Nat gateway still exists")
		}
	}

	return nil
}

func testAccCheckNatV2GatewayExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		natClient, err := config.natGatewayV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud nat client: %s", err)
		}

		found, err := natgateways.Get(natClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Nat gateway not found")
		}

		return nil
	}
}

const testAccNatV2Gateway_basic = `
resource "huaweicloud_networking_router_v2" "router_1" {
  name = "router_1"
  admin_state_up = "true"
}

resource "huaweicloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloud_networking_subnet_v2" "subnet_1" {
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
}

resource "huaweicloud_networking_router_interface_v2" "int_1" {
  subnet_id = "${huaweicloud_networking_subnet_v2.subnet_1.id}"
  router_id = "${huaweicloud_networking_router_v2.router_1.id}"
}

resource "huaweicloud_nat_gateway_v2" "nat_1" {
  name   = "nat_1"
  description = "test for terraform"
  spec = "1"
  internal_network_id = "${huaweicloud_networking_network_v2.network_1.id}"
  router_id = "${huaweicloud_networking_router_v2.router_1.id}"
  depends_on = ["huaweicloud_networking_router_interface_v2.int_1"]
}
`

func testAccNatV2Gateway_epsId(name string) string {
	return fmt.Sprintf(`
	resource "huaweicloud_vpc" "vpc_1" {
	  name = "%s"
	  cidr = "192.168.0.0/16"
	}
	
	resource "huaweicloud_vpc_subnet" "subnet_1" {
	  name       = "%s"
	  cidr       = "192.168.199.0/24"
	  gateway_ip = "192.168.199.1"
	  vpc_id     = huaweicloud_vpc.vpc_1.id
	}
	
	resource "huaweicloud_nat_gateway" "nat_1" {
	  name   = "%s"
	  description = "test for terraform"
	  spec = "1"
	  internal_network_id = huaweicloud_vpc_subnet.subnet_1.id
	  router_id = huaweicloud_vpc.vpc_1.id
	  enterprise_project_id = "%s"
	}
	`, name, name, name, OS_ENTERPRISE_PROJECT_ID)
}

const testAccNatV2Gateway_update = `
resource "huaweicloud_networking_router_v2" "router_1" {
  name = "router_1"
  admin_state_up = "true"
}

resource "huaweicloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloud_networking_subnet_v2" "subnet_1" {
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
}

resource "huaweicloud_networking_router_interface_v2" "int_1" {
  subnet_id = "${huaweicloud_networking_subnet_v2.subnet_1.id}"
  router_id = "${huaweicloud_networking_router_v2.router_1.id}"
}

resource "huaweicloud_nat_gateway_v2" "nat_1" {
  name   = "nat_1_updated"
  description = "nat_1 updated"
  spec = "2"
  internal_network_id = "${huaweicloud_networking_network_v2.network_1.id}"
  router_id = "${huaweicloud_networking_router_v2.router_1.id}"
  depends_on = ["huaweicloud_networking_router_interface_v2.int_1"]
}
`
