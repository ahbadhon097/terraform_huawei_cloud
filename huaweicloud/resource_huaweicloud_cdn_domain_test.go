package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/chnsz/golangsdk/openstack/cdn/v1/domains"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccCdnDomain_basic(t *testing.T) {
	var domain domains.CdnDomain

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCDN(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCdnDomainV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCdnDomainV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCdnDomainV1Exists("huaweicloud_cdn_domain.domain_1", &domain),
					resource.TestCheckResourceAttr(
						"huaweicloud_cdn_domain.domain_1", "name", HW_CDN_DOMAIN_NAME),
				),
			},
		},
	})
}

func TestAccCdnDomain_cache(t *testing.T) {
	var domain domains.CdnDomain

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCDN(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCdnDomainV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCdnDomainV1_cache,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCdnDomainV1Exists("huaweicloud_cdn_domain.domain_1", &domain),
					resource.TestCheckResourceAttr(
						"huaweicloud_cdn_domain.domain_1", "name", HW_CDN_DOMAIN_NAME),
					resource.TestCheckResourceAttr(
						"huaweicloud_cdn_domain.domain_1", "cache_settings.0.rules.0.rule_type", "0"),
					resource.TestCheckResourceAttr(
						"huaweicloud_cdn_domain.domain_1", "cache_settings.0.rules.0.ttl", "180"),
					resource.TestCheckResourceAttr(
						"huaweicloud_cdn_domain.domain_1", "cache_settings.0.rules.0.ttl_type", "4"),
					resource.TestCheckResourceAttr(
						"huaweicloud_cdn_domain.domain_1", "cache_settings.0.rules.0.priority", "2"),
				),
			},
		},
	})
}

func TestAccCdnDomain_configs(t *testing.T) {
	var domain domains.CdnDomain

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCDN(t)
			testAccPreCheckCERT(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCdnDomainV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCdnDomainV1_configs,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCdnDomainV1Exists("huaweicloud_cdn_domain.domain_1", &domain),
					resource.TestCheckResourceAttr(
						"huaweicloud_cdn_domain.domain_1", "name", HW_CDN_DOMAIN_NAME),
					resource.TestCheckResourceAttr(
						"huaweicloud_cdn_domain.domain_1", "configs.0.origin_protocol", "http"),
					resource.TestCheckResourceAttr(
						"huaweicloud_cdn_domain.domain_1", "configs.0.https_settings.0.certificate_name", "terraform-test"),
					resource.TestCheckResourceAttr(
						"huaweicloud_cdn_domain.domain_1", "configs.0.https_settings.0.https_status", "on"),
					resource.TestCheckResourceAttr(
						"huaweicloud_cdn_domain.domain_1", "configs.0.https_settings.0.http2_status", "on"),
					resource.TestCheckResourceAttr(
						"huaweicloud_cdn_domain.domain_1", "configs.0.cache_url_parameter_filter.0.type", "ignore_url_params"),
					resource.TestCheckResourceAttr(
						"huaweicloud_cdn_domain.domain_1", "configs.0.retrieval_request_header.0.name", "test-name"),
					resource.TestCheckResourceAttr(
						"huaweicloud_cdn_domain.domain_1", "configs.0.url_signing.0.status", "off"),
					resource.TestCheckResourceAttr(
						"huaweicloud_cdn_domain.domain_1", "configs.0.compress.0.status", "off"),
					resource.TestCheckResourceAttr(
						"huaweicloud_cdn_domain.domain_1", "configs.0.force_redirect.0.status", "on"),
				),
			},
		},
	})
}

func testAccCheckCdnDomainV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	cdnClient, err := config.CdnV1Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud CDN Domain client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_cdn_domain" {
			continue
		}

		found, err := domains.Get(cdnClient, rs.Primary.ID, nil).Extract()
		if err == nil && found.DomainStatus != "deleting" {
			return fmtp.Errorf("Destroying CDN domain failed or domain still exists")
		}
	}

	return nil
}

func testAccCheckCdnDomainV1Exists(n string, domain *domains.CdnDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("CDN Domain Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		cdnClient, err := config.CdnV1Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud CDN Domain client: %s", err)
		}

		found, err := domains.Get(cdnClient, rs.Primary.ID, nil).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmtp.Errorf("CDN Domain not found")
		}

		*domain = *found
		return nil
	}
}

var testAccCdnDomainV1_basic = fmt.Sprintf(`
resource "huaweicloud_cdn_domain" "domain_1" {
  name   = "%s"
  type   = "wholeSite"
  enterprise_project_id = 0
  sources {
      active = 1
      origin = "100.254.53.75"
      origin_type  = "ipaddr"
  }
}
`, HW_CDN_DOMAIN_NAME)

var testAccCdnDomainV1_cache = fmt.Sprintf(`
resource "huaweicloud_cdn_domain" "domain_1" {
  name                  = "%s"
  type                  = "wholeSite"
  enterprise_project_id = 0

  sources {
    active      = 1
    origin      = "100.254.53.75"
    origin_type = "ipaddr"
  }

  cache_settings {
    rules {
      rule_type = 0
      ttl       = 180
      ttl_type  = 4
      priority  = 2
    }
  }
}
`, HW_CDN_DOMAIN_NAME)

var testAccCdnDomainV1_configs = fmt.Sprintf(`
resource "huaweicloud_cdn_domain" "domain_1" {
  name                  = "%s"
  type                  = "wholeSite"
  enterprise_project_id = 0

  sources {
    active      = 1
    origin      = "100.254.53.75"
    origin_type = "ipaddr"
  }

  configs {
	origin_protocol = "http"

    https_settings {
      certificate_name = "terraform-test"
      certificate_body = file("%s")
      http2_enabled    = true
      https_enabled    = true
      private_key      = file("%s")
    }

    cache_url_parameter_filter {
      type = "ignore_url_params"
    }

    retrieval_request_header {
      name   = "test-name"
      value  = "test-val"
      action = "set"
    }

    http_response_header {
      name   = "test-name"
      value  = "test-val"
      action = "set"
    }

    url_signing {
      enabled = false
    }

    compress {
      enabled = false
    }

    force_redirect {
      enabled = true
      type   = "http"
    }
  }
}
`, HW_CDN_DOMAIN_NAME, HW_CDN_CERT_PATH, HW_CDN_PRIVATE_KEY_PATH)
