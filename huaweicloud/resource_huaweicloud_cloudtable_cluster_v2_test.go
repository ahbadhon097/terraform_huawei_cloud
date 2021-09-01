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

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccCloudtableClusterV2_basic(t *testing.T) {
	resourceName := "huaweicloud_cloudtable_cluster_v2.test"
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCloudTable(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudtableClusterV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudtableClusterV2_basic(rName, HW_REGION_NAME,
					HW_CLOUDTABLE_AVAILABILITY_ZONE),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudtableClusterV2Exists(resourceName),
				),
			},
		},
	})
}

func testAccCloudtableClusterV2_basic(rName string, rRegion string, availabilityZone string) string {
	return fmt.Sprintf(`
resource "huaweicloud_networking_secgroup_v2" "secgroup" {
  region      = "%s"
  name        = "%s"
  description = "terraform security group acceptance test"
  timeouts {
    delete = "20m"
  }
}

resource "huaweicloud_vpc" "test" {
  region = "%s"
  name   = "%s"
  cidr   = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%s"
  cidr       = "192.168.0.0/20"
  vpc_id     = huaweicloud_vpc.test.id
  gateway_ip = "192.168.0.1"
}

resource "huaweicloud_cloudtable_cluster_v2" "test" {
  region            = "%s"
  availability_zone = "%s"
  name              = "%s"
  rs_num            = 2
  security_group_id = huaweicloud_networking_secgroup_v2.secgroup.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id
  storage_type      = "COMMON"
}`, rRegion, rName, rRegion, rName, rName, rRegion, availabilityZone, rName)
}

func testAccCheckCloudtableClusterV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	client, err := config.CloudtableV2Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating sdk client, err=%s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_cloudtable_cluster_v2" {
			continue
		}

		url, err := replaceVarsForTest(rs, "clusters/{id}")
		if err != nil {
			return err
		}
		url = client.ServiceURL(url)

		r := golangsdk.Result{}
		_, r.Err = client.Get(url, &r.Body, &golangsdk.RequestOpts{
			MoreHeaders: map[string]string{
				"Content-Type": "application/json",
				"X-Language":   "en-us",
			}})
		if r.Err != nil {
			if _, ok := r.Err.(golangsdk.ErrDefault404); ok {
				return nil
			}
			return fmtp.Errorf("huaweicloud_cloudtable_cluster_v2 query exception %s", url)
		}

		status, _ := navigateValue(r.Body, []string{"status"}, nil)
		if status == "303" {
			return nil
		}

		return fmtp.Errorf("huaweicloud_cloudtable_cluster_v2 still exists at %s", url)
	}

	return nil
}

func testAccCheckCloudtableClusterV2Exists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*config.Config)
		client, err := config.CloudtableV2Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating sdk client, err=%s", err)
		}

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmtp.Errorf("Error checking %s exist, err=not found this resource", resourceName)
		}

		url, err := replaceVarsForTest(rs, "clusters/{id}")
		if err != nil {
			return fmtp.Errorf("Error checking %s exist, err=building url failed: %s", resourceName, err)
		}
		url = client.ServiceURL(url)

		_, err2 := client.Get(url, nil, &golangsdk.RequestOpts{
			MoreHeaders: map[string]string{
				"Content-Type": "application/json",
				"X-Language":   "en-us",
			}})
		if err2 != nil {
			if _, ok := err2.(golangsdk.ErrDefault404); ok {
				return fmtp.Errorf("%s is not exist", resourceName)
			}
			return fmtp.Errorf("Error checking %s exist, err=send request failed: %s", resourceName, err)
		}
		return nil
	}
}
