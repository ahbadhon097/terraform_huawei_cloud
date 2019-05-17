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

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/huaweicloud/golangsdk"
)

func TestAccDwsCluster_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDws(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDwsClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDwsCluster_basic(acctest.RandString(10)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDwsClusterExists(),
				),
			},
		},
	})
}

func testAccDwsCluster_basic(val string) string {
	return fmt.Sprintf(`
resource "huaweicloud_networking_secgroup_v2" "secgroup" {
  name = "security_group_2%s"
  description = "terraform security group"
}

resource "huaweicloud_dws_cluster" "cluster" {
  node_type = "dws.m3.xlarge"
  number_of_node = 3
  network_id = "%s"
  vpc_id = "%s"
  security_group_id = "${huaweicloud_networking_secgroup_v2.secgroup.id}"
  availability_zone = "%s"
  name = "terraform_dws_cluster_test%s"
  user_name = "test_cluster_admin"
  user_pwd = "cluster123@!"

  timeouts {
    create = "30m"
    delete = "30m"
  }
}
	`, val, OS_NETWORK_ID, OS_VPC_ID, OS_AVAILABILITY_ZONE, val)
}

func testAccCheckDwsClusterDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.sdkClient(OS_REGION_NAME, "dws", serviceProjectLevel)
	if err != nil {
		return fmt.Errorf("Error creating sdk client, err=%s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_dws_cluster" {
			continue
		}

		url, err := replaceVarsForTest(rs, "clusters/{id}")
		if err != nil {
			return err
		}
		url = client.ServiceURL(url)

		_, err = client.Get(
			url, nil,
			&golangsdk.RequestOpts{MoreHeaders: map[string]string{"Content-Type": "application/json"}})
		if err == nil {
			return fmt.Errorf("huaweicloud_dws_cluster still exists at %s", url)
		}
	}

	return nil
}

func testAccCheckDwsClusterExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		client, err := config.sdkClient(OS_REGION_NAME, "dws", serviceProjectLevel)
		if err != nil {
			return fmt.Errorf("Error creating sdk client, err=%s", err)
		}

		rs, ok := s.RootModule().Resources["huaweicloud_dws_cluster.cluster"]
		if !ok {
			return fmt.Errorf("Error checking huaweicloud_dws_cluster.cluster exist, err=not found huaweicloud_dws_cluster.cluster")
		}

		url, err := replaceVarsForTest(rs, "clusters/{id}")
		if err != nil {
			return fmt.Errorf("Error checking huaweicloud_dws_cluster.cluster exist, err=building url failed: %s", err)
		}
		url = client.ServiceURL(url)

		_, err = client.Get(
			url, nil,
			&golangsdk.RequestOpts{MoreHeaders: map[string]string{"Content-Type": "application/json"}})
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return fmt.Errorf("huaweicloud_dws_cluster.cluster is not exist")
			}
			return fmt.Errorf("Error checking huaweicloud_dws_cluster.cluster exist, err=send request failed: %s", err)
		}
		return nil
	}
}
