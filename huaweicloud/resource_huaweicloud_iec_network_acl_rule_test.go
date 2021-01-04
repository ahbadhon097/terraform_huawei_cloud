package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/iec/v1/firewalls"
)

func TestAccIecNetworkACLRule_basic(t *testing.T) {
	rName := fmt.Sprintf("iec-acl-%s", acctest.RandString(5))
	resourceKey := "huaweicloud_iec_network_acl.acl_test"
	var fwGroup firewalls.Firewall
	checkMapBasic := make(map[string]string)
	checkMapBasic["durection"] = "ingress"
	checkMapBasic["protocol"] = "tcp"
	checkMapBasic["action"] = "allow"
	checkMapBasic["destPort"] = "445"
	checkMapUpdate := make(map[string]string)
	checkMapUpdate["durection"] = "ingress"
	checkMapUpdate["protocol"] = "udp"
	checkMapUpdate["action"] = "deny"
	checkMapUpdate["destPort"] = "23-30"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIecNetworkACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIecNetworkACLRule_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIecNetworkACLRuleExists(resourceKey, &fwGroup),
					testAccCheckIecNetworkACLRuleParameter(&fwGroup, checkMapBasic),
				),
			},
			{
				Config: testAccIecNetworkACLRule_basic_update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIecNetworkACLRuleExists(resourceKey, &fwGroup),
					testAccCheckIecNetworkACLRuleParameter(&fwGroup, checkMapUpdate),
				),
			},
		},
	})
}

func testAccCheckIecNetworkACLRuleDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	iecV1Client, err := config.IECV1Client(HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud IEC client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_iec_network_acl" {
			continue
		}

		_, err := firewalls.Get(iecV1Client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("IEC network acl still exists")
		}
	}

	return nil
}

func testAccCheckIecNetworkACLRuleExists(n string, resource *firewalls.Firewall) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		iecV1Client, err := config.IECV1Client(HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating Huaweicloud IEC client: %s", err)
		}

		found, err := firewalls.Get(iecV1Client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		if len(found.IngressFWPolicy.FirewallRules) != 0 || len(found.EgressFWPolicy.FirewallRules) != 0 {
			*resource = *found
		} else {
			return fmt.Errorf("IEC Network ACL Rule not found")
		}
		return nil
	}
}

func testAccCheckIecNetworkACLRuleParameter(resource *firewalls.Firewall, checkMap map[string]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(resource.IngressFWPolicy.FirewallRules) != 0 {
			if resource.IngressFWPolicy.FirewallRules[0].Protocol != checkMap["protocol"] {
				return fmt.Errorf("[%s]: The Protocol of IEC Network ACL Ingress Rule is not right.", resource.IngressFWPolicy.FirewallRules[0].Protocol)
			}
			if resource.IngressFWPolicy.FirewallRules[0].Action != checkMap["action"] {
				return fmt.Errorf("[%s]: The Action of IEC Network ACL Ingress Rule is not right.", resource.IngressFWPolicy.FirewallRules[0].Action)
			}
			if resource.IngressFWPolicy.FirewallRules[0].DstPort != checkMap["destPort"] {
				return fmt.Errorf("[%s]: The Destination Port of IEC Network ACL Ingress is not right.", resource.IngressFWPolicy.FirewallRules[0].DstPort)
			}
		} else if len(resource.EgressFWPolicy.FirewallRules) != 0 {
			if resource.EgressFWPolicy.FirewallRules[0].Protocol != checkMap["protocol"] {
				return fmt.Errorf("[%s]: The Protocol of IEC Network ACL Egress Rule is not right.", resource.EgressFWPolicy.FirewallRules[0].Protocol)
			}
			if resource.EgressFWPolicy.FirewallRules[0].Action != checkMap["action"] {
				return fmt.Errorf("[%s]: The Action of IEC Network ACL Egress is not right.", resource.EgressFWPolicy.FirewallRules[0].Action)
			}
			if resource.EgressFWPolicy.FirewallRules[0].DstPort != checkMap["destPort"] {
				return fmt.Errorf("[%s]: The Destination Port of IEC Network ACL Egress is not right.", resource.EgressFWPolicy.FirewallRules[0].DstPort)
			}
		} else {
			return fmt.Errorf("IEC Network ACL Rule not found")
		}
		return nil
	}
}

func testAccIecNetworkACLRule_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_iec_network_acl" "acl_test" {
  name = "%s"
}

resource "huaweicloud_iec_network_acl_rule" "rule_test" {
  network_acl_id         = huaweicloud_iec_network_acl.acl_test.id
  direction              = "ingress"
  protocol               = "tcp"
  action                 = "allow"
  source_ip_address      = "0.0.0.0/0"
  destination_ip_address = "0.0.0.0/0"
  destination_port       = "445"
  enabled                = true
}
`, rName)
}

func testAccIecNetworkACLRule_basic_update(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_iec_network_acl" "acl_test" {
  name = "%s"
}

resource "huaweicloud_iec_network_acl_rule" "rule_test" {
  network_acl_id         = huaweicloud_iec_network_acl.acl_test.id
  direction              = "ingress"
  protocol               = "udp"
  action                 = "deny"
  source_ip_address      = "0.0.0.0/0"
  destination_ip_address = "0.0.0.0/0"
  destination_port       = "23-30"
  enabled                = true
}
`, rName)
}
