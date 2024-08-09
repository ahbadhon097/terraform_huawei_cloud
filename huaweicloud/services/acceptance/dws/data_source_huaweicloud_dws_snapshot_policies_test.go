package dws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSnapshotPolicies_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dws_snapshot_policies.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
			acceptance.TestAccPreCheckDwsSnapshotPolicyName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSnapshotPolicies_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_exist_snapshot_policy", "true"),
					resource.TestMatchResourceAttr(dataSource, "policies.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.strategy"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.backup_level"),
					resource.TestMatchResourceAttr(dataSource, "policies.0.next_fire_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(dataSource, "policies.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
		},
	})
}

func testDataSourceSnapshotPolicies_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dws_snapshot_policies" "test" {
  cluster_id = "%[1]s"
}

output "is_exist_snapshot_policy" {
  value  = contains(data.huaweicloud_dws_snapshot_policies.test.policies[*].name, "%[2]s")
}
`, acceptance.HW_DWS_CLUSTER_ID, acceptance.HW_DWS_SNAPSHOT_POLICY_NAME)
}
