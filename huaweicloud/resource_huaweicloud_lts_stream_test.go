package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk/openstack/lts/huawei/logstreams"
)

func TestAccLogTankStreamV2_basic(t *testing.T) {
	var stream logstreams.LogStream
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLogTankStreamV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLogTankStreamV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogTankStreamV2Exists(
						"huaweicloud_lts_stream.testacc_stream", &stream),
					resource.TestCheckResourceAttr(
						"huaweicloud_lts_stream.testacc_stream", "stream_name", "testacc_stream"),
				),
			},
		},
	})
}

func testAccCheckLogTankStreamV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	ltsclient, err := config.ltsV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud LTS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_lts_stream" {
			continue
		}

		group_id := rs.Primary.Attributes["group_id"]
		_, err = logstreams.List(ltsclient, group_id).Extract()
		if err == nil {
			return fmt.Errorf("Log group (%s) still exists.", rs.Primary.ID)
		}

	}
	return nil
}

func testAccCheckLogTankStreamV2Exists(n string, stream *logstreams.LogStream) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		ltsclient, err := config.ltsV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud LTS client: %s", err)
		}

		group_id := rs.Primary.Attributes["group_id"]
		streams, err := logstreams.List(ltsclient, group_id).Extract()
		if err != nil {
			return fmt.Errorf("Log stream get list err: %s", err.Error())
		}
		for _, logstream := range streams.LogStreams {
			if logstream.ID == rs.Primary.ID {
				*stream = logstream
				return nil
			}
		}

		return fmt.Errorf("Error HuaweiCloud log stream %s: No Found", rs.Primary.ID)
	}
}

const testAccLogTankStreamV2_basic = `
resource "huaweicloud_lts_group" "testacc_group" {
	group_name  = "testacc_group"
	ttl_in_days = 1
}
resource "huaweicloud_lts_stream" "testacc_stream" {
  group_id = "${huaweicloud_lts_group.testacc_group.id}"
  stream_name = "testacc_stream"
}
`
