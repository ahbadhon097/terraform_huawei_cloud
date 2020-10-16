// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    AUTO GENERATED CODE     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file at
//     https://www.github.com/huaweicloud/magic-modules
//
// ----------------------------------------------------------------------------

package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk"
)

func TestAccCssClusterV1_basic(t *testing.T) {
	randName := acctest.RandString(6)
	resourceName := "huaweicloud_css_cluster_v1.cluster"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCssClusterV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCssClusterV1_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCssClusterV1Exists(),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("terraform_test_cluster%s", randName)),
					resource.TestCheckResourceAttr(resourceName, "expect_node_num", "1"),
					resource.TestCheckResourceAttr(resourceName, "engine_type", "elasticsearch"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
			{
				Config: testAccCssClusterV1_update(randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCssClusterV1Exists(),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar_update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key_update", "value"),
				),
			},
		},
	})
}

func testAccCheckCssClusterV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.sdkClient(OS_REGION_NAME, "css", serviceProjectLevel)
	if err != nil {
		return fmt.Errorf("Error creating sdk client, err=%s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_css_cluster_v1" {
			continue
		}

		url, err := replaceVarsForTest(rs, "clusters/{id}")
		if err != nil {
			return err
		}
		url = client.ServiceURL(url)

		_, err = client.Get(url, nil, &golangsdk.RequestOpts{
			MoreHeaders: map[string]string{"Content-Type": "application/json"}})
		if err == nil {
			return fmt.Errorf("huaweicloud_css_cluster_v1 still exists at %s", url)
		}
	}

	return nil
}

func testAccCheckCssClusterV1Exists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		client, err := config.sdkClient(OS_REGION_NAME, "css", serviceProjectLevel)
		if err != nil {
			return fmt.Errorf("Error creating sdk client, err=%s", err)
		}

		rs, ok := s.RootModule().Resources["huaweicloud_css_cluster_v1.cluster"]
		if !ok {
			return fmt.Errorf("Error checking huaweicloud_css_cluster_v1.cluster exist, err=not found this resource")
		}

		url, err := replaceVarsForTest(rs, "clusters/{id}")
		if err != nil {
			return fmt.Errorf("Error checking huaweicloud_css_cluster_v1.cluster exist, err=building url failed: %s", err)
		}
		url = client.ServiceURL(url)

		_, err = client.Get(url, nil, &golangsdk.RequestOpts{
			MoreHeaders: map[string]string{"Content-Type": "application/json"}})
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return fmt.Errorf("huaweicloud_css_cluster_v1.cluster is not exist")
			}
			return fmt.Errorf("Error checking huaweicloud_css_cluster_v1.cluster exist, err=send request failed: %s", err)
		}
		return nil
	}
}

func testAccCssClusterV1_basic(val string) string {
	return fmt.Sprintf(`
resource "huaweicloud_networking_secgroup_v2" "secgroup" {
  name = "terraform_test_security_group%s"
  description = "terraform security group acceptance test"
}

resource "huaweicloud_css_cluster_v1" "cluster" {
  name = "terraform_test_cluster%s"
  engine_version  = "7.1.1"
  expect_node_num = 1

  node_config {
    flavor = "ess.spec-4u16g"
    network_info {
      security_group_id = huaweicloud_networking_secgroup_v2.secgroup.id
      subnet_id = "%s"
      vpc_id = "%s"
    }
    volume {
      volume_type = "HIGH"
      size = 40
    }
    availability_zone = "%s"
  }
  tags = {
    foo = "bar"
    key = "value"
  }
}
	`, val, val, OS_NETWORK_ID, OS_VPC_ID, OS_AVAILABILITY_ZONE)
}

func testAccCssClusterV1_update(val string) string {
	return fmt.Sprintf(`
resource "huaweicloud_networking_secgroup_v2" "secgroup" {
  name = "terraform_test_security_group%s"
  description = "terraform security group acceptance test"
}

resource "huaweicloud_css_cluster_v1" "cluster" {
  name = "terraform_test_cluster%s"
  engine_version  = "7.1.1"
  expect_node_num = 1

  node_config {
    flavor = "ess.spec-4u16g"
    network_info {
      security_group_id = huaweicloud_networking_secgroup_v2.secgroup.id
      subnet_id = "%s"
      vpc_id = "%s"
    }
    volume {
      volume_type = "HIGH"
      size = 40
    }
    availability_zone = "%s"
  }
  tags = {
    foo = "bar_update"
    key_update = "value"
  }
}
	`, val, val, OS_NETWORK_ID, OS_VPC_ID, OS_AVAILABILITY_ZONE)
}
