package huaweicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccComputeV2FloatingIP_importBasic(t *testing.T) {
	resourceName := "huaweicloud_compute_floatingip_v2.fip_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDeprecated(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2FloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2FloatingIP_basic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
