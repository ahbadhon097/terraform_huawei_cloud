package dli

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/dli/v1/queues"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	act "github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dli"
)

func getDliQueueResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.DliV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Dli v1 client, err=%s", err)
	}

	result := queues.Get(client, state.Primary.Attributes["name"])
	return result.Body, result.Err
}

func TestAccDliQueue_basic(t *testing.T) {
	var (
		obj   interface{}
		rName = act.RandomAccResourceName()

		typeSQL     = "huaweicloud_dli_queue.default"
		typeGeneral = "huaweicloud_dli_queue.general"

		rcForTypeSQL     = acceptance.InitResourceCheck(typeSQL, &obj, getDliQueueResourceFunc)
		rcForTypeGeneral = acceptance.InitResourceCheck(typeGeneral, &obj, getDliQueueResourceFunc)
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { act.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rcForTypeSQL.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDliQueue_basic(rName, dli.CU16),
				Check: resource.ComposeTestCheckFunc(
					rcForTypeSQL.CheckResourceExists(),
					resource.TestCheckResourceAttr(typeSQL, "name", rName+"_sql"),
					resource.TestCheckResourceAttr(typeSQL, "queue_type", dli.QueueTypeSQL),
					resource.TestCheckResourceAttr(typeSQL, "cu_count", fmt.Sprintf("%d", dli.CU16)),
					resource.TestCheckResourceAttr(typeSQL, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(typeSQL, "resource_mode", "1"),
					resource.TestCheckResourceAttrSet(typeSQL, "create_time"),
					rcForTypeGeneral.CheckResourceExists(),
					resource.TestCheckResourceAttr(typeGeneral, "name", rName+"_general"),
					resource.TestCheckResourceAttr(typeGeneral, "queue_type", dli.QueueTypeGeneral),
					resource.TestCheckResourceAttr(typeGeneral, "cu_count", fmt.Sprintf("%d", dli.CU16)),
					resource.TestCheckResourceAttr(typeSQL, "resource_mode", "1"),
					resource.TestCheckResourceAttrSet(typeGeneral, "create_time"),
					waitForDeletionCooldownComplete(),
				),
			},
			{
				ResourceName:      typeSQL,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccQueueImportStateFunc(typeSQL),
				ImportStateVerifyIgnore: []string{
					"tags",
				},
			},
		},
	})
}

func testAccQueueImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}
		name := rs.Primary.Attributes["name"]
		if name == "" {
			return "", fmt.Errorf("the queue name is incorrect, got '%s'", name)
		}
		return name, nil
	}
}

func testAccDliQueue_base(rName string, cuCount int) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dli_elastic_resource_pool" "test" {
  name   = "%[2]s"
  min_cu = %[3]d
  max_cu = %[3]d
  cidr   = cidrsubnet(huaweicloud_vpc.test.cidr, 3, 1)

  enterprise_project_id = "0"
}
`, common.TestVpc(rName), rName, cuCount)
}

func testAccDliQueue_basic(rName string, cuCount int) string {
	return fmt.Sprintf(`
%[1]s

# The default type is SQL
resource "huaweicloud_dli_queue" "default" {
  elastic_resource_pool_name = huaweicloud_dli_elastic_resource_pool.test.name
  resource_mode              = 1

  name     = "%[2]s_sql"
  cu_count = %[3]d

  tags = {
    foo = "bar"
  }
}

resource "huaweicloud_dli_queue" "general" {
  elastic_resource_pool_name = huaweicloud_dli_elastic_resource_pool.test.name
  resource_mode              = 1

  name       = "%[2]s_general"
  cu_count   = %[3]d
  queue_type = "general"
}
`, testAccDliQueue_base(rName, cuCount*4), // Prevent resource creation failure due to insufficient total resource CUs
		rName, cuCount)
}

func TestAccDliQueue_cidr(t *testing.T) {
	rName := act.RandomAccResourceName()
	resourceName := "huaweicloud_dli_queue.test"

	var obj queues.CreateOpts
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDliQueueResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { act.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDliQueue_cidr(rName, "172.16.0.0/21"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "queue_type", dli.QueueTypeSQL),
					resource.TestCheckResourceAttr(resourceName, "cu_count", "16"),
					resource.TestCheckResourceAttr(resourceName, "resource_mode", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpc_cidr", "172.16.0.0/21"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
				),
			},
			{

				Config: testAccDliQueue_cidr(rName, "172.16.0.0/18"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "queue_type", dli.QueueTypeSQL),
					resource.TestCheckResourceAttr(resourceName, "cu_count", "16"),
					resource.TestCheckResourceAttr(resourceName, "resource_mode", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpc_cidr", "172.16.0.0/18"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccQueueImportStateFunc(resourceName),
				ImportStateVerifyIgnore: []string{
					"tags",
				},
			},
		},
	})
}

func testAccDliQueue_cidr(rName string, cidr string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dli_queue" "test" {
  name          = "%s"
  cu_count      = 16
  resource_mode = 1
  vpc_cidr      = "%s"

  tags = {
    foo = "bar"
  }
}`, rName, cidr)
}

func TestAccDliQueue_scalingPolicies(t *testing.T) {
	// Creating a queue will create an elastic resource pool with the same name.
	elasticResourcePoolName := acceptance.RandomAccResourceName()
	queueName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_dli_queue.test"
	var obj queues.CreateOpts
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDliQueueResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccQueue_scalingPolicies_step1(elasticResourcePoolName, queueName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "elastic_resource_pool_name", elasticResourcePoolName),
					resource.TestCheckResourceAttr(resourceName, "name", queueName),
					resource.TestCheckResourceAttr(resourceName, "scaling_policies.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "scaling_policies.0.%", "5"),
				),
			},
			{
				Config: testAccQueue_scalingPolicies_step2(elasticResourcePoolName, queueName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "scaling_policies.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "scaling_policies.0.priority", "1"),
					resource.TestCheckResourceAttr(resourceName, "scaling_policies.0.impact_start_time", "00:00"),
					resource.TestCheckResourceAttr(resourceName, "scaling_policies.0.impact_stop_time", "24:00"),
					resource.TestCheckResourceAttr(resourceName, "scaling_policies.0.min_cu", "16"),
					resource.TestCheckResourceAttr(resourceName, "scaling_policies.0.max_cu", "16"),
					waitForDeletionCooldownComplete(),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccQueueImportStateFunc(resourceName),
			},
		},
	})
}

func testAccQueue_scalingPolicies_step1(elasticResourcePoolName string, queueName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dli_elastic_resource_pool" "test" {
  name                  = "%[1]s"
  max_cu                = 64
  min_cu                = 64
  enterprise_project_id = "0"
}

resource "huaweicloud_dli_queue" "test" {
  name                       = "%[2]s"
  cu_count                   = 16
  resource_mode              = 1
  elastic_resource_pool_name = huaweicloud_dli_elastic_resource_pool.test.name

  scaling_policies {
    priority          = 1
    impact_start_time = "00:00"
    impact_stop_time  = "24:00"
    min_cu            = 16
    max_cu            = 16
  }
  
  scaling_policies {
    priority          = 2
    impact_start_time = "00:00"
    impact_stop_time  = "01:00"
    min_cu            = 20
    max_cu            = 28
  }
}`, elasticResourcePoolName, queueName)
}

func testAccQueue_scalingPolicies_step2(elasticResourcePoolName string, queueName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dli_elastic_resource_pool" "test" {
  name                  = "%[1]s"
  max_cu                = 64
  min_cu                = 64
  enterprise_project_id = "0"
}

resource "huaweicloud_dli_queue" "test" {
  name                       = "%[2]s"
  cu_count                   = 16
  resource_mode              = 1
  elastic_resource_pool_name = huaweicloud_dli_elastic_resource_pool.test.name

  scaling_policies {
    priority          = 1
    impact_start_time = "00:00"
    impact_stop_time  = "24:00"
    max_cu            = 16
    min_cu            = 16
  }
}`, elasticResourcePoolName, queueName)
}

func TestAccDliQueue_sparkDriver(t *testing.T) {
	rName := act.RandomAccResourceName()
	resourceName := "huaweicloud_dli_queue.test"

	var obj queues.CreateOpts
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDliQueueResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { act.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDliQueue_sparkDriver_step1(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "spark_driver.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "spark_driver.0.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "spark_driver.0.max_instance", "2"),
					resource.TestCheckResourceAttr(resourceName, "spark_driver.0.max_concurrent", "1"),
					resource.TestCheckResourceAttr(resourceName, "spark_driver.0.max_prefetch_instance", "0"),
				),
			},
			{
				Config: testAccDliQueue_sparkDriver_step2(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "spark_driver.0.max_prefetch_instance", "4"),
				),
			},
			{
				Config: testAccDliQueue_sparkDriver_step3(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "scaling_policies.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccQueueImportStateFunc(resourceName),
				ImportStateVerifyIgnore: []string{
					"tags",
				},
			},
		},
	})
}

func testAccDliQueue_sparkDriver_step1(queueName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dli_queue" "test" {
  name     = "%s"
  cu_count = 64

  spark_driver {
    max_instance          = 2
    max_concurrent        = 1
    max_prefetch_instance = "0"
  }
}`, queueName)
}

// Modify "max_prefetch_instance" parameter, and remove the "max_instance" and "max_concurrent" parameters。
func testAccDliQueue_sparkDriver_step2(queueName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dli_queue" "test" {
  name     = "%s"
  cu_count = 64

  spark_driver {
    max_prefetch_instance = "4"
  }
}`, queueName)
}

// Remove spark_driver parameters
func testAccDliQueue_sparkDriver_step3(queueName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dli_queue" "test" {
  name     = "%s"
  cu_count = 64
}`, queueName)
}
