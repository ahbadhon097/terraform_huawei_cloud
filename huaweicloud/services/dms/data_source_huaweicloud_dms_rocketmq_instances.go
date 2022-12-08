package dms

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceDmsRocketMQInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceDmsRocketMQInstancesRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of the DMS RocketMQ instance.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the rocketMQ instance.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the status of the DMS RocketMQ instance.`,
			},
			"exact_match_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: `Specifies whether to search for the instance that precisely matches a specified
instance name.`,
			},
			"instances": {
				Type:        schema.TypeList,
				Elem:        DmsRocketMQInstancesInstanceRefSchema(),
				Computed:    true,
				Description: `Indicates the list of DMS rocketMQ instances.`,
			},
		},
	}
}

func DmsRocketMQInstancesInstanceRefSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of the DMS RocketMQ instance.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the status of the DMS RocketMQ instance.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the description of the DMS RocketMQ instance.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the DMS RocketMQ instance type. Value: cluster.`,
			},
			"specification": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the instance specification.`,
			},
			"engine_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the version of the RocketMQ engine. Value: 4.8.0.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of a VPC.`,
			},
			"flavor_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates a product ID.`,
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of a security group.`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of a subnet.`,
			},
			"availability_zones": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Description: `Indicates the list of availability zone names, where
`,
			},
			"maintain_begin": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the time at which the maintenance window starts. The format is HH:mm:ss.`,
			},
			"maintain_end": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the time at which the maintenance window ends. The format is HH:mm:ss.`,
			},
			"storage_space": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the message storage capacity, unit is GB.`,
			},
			"used_storage_space": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the used message storage space. Unit: GB.`,
			},
			"enable_publicip": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether to enable public access.`,
			},
			"publicip_id": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `Indicates the ID of the EIP bound to the instance.
`,
			},
			"publicip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the public IP address.`,
			},
			"ssl_enable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the RocketMQ SASL_SSL is enabled. Defaults to false.`,
			},
			"cross_vpc_accesses": {
				Type:        schema.TypeList,
				Elem:        DmsRocketMQInstancesInstanceRefCrossVpcInfoRefSchema(),
				Computed:    true,
				Description: `Indicates the Cross-VPC access information.`,
			},
			"storage_spec_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the storage I/O specification.`,
			},
			"ipv6_enable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether to support IPv6. Defaults to false.`,
			},
			"node_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the node quantity.`,
			},
			"new_spec_billing_enable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates the whether billing based on new specifications is enabled.`,
			},
			"enable_acl": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether access control is enabled.`,
			},
			"broker_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Specifies the broker numbers. Defaults to 1.`,
			},
			"namesrv_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the metadata address.`,
			},
			"broker_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the service data address.`,
			},
			"public_namesrv_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the public network metadata address.`,
			},
			"public_broker_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the public network service data address.`,
			},
			"resource_spec_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the resource specifications.`,
			},
		},
	}
	return &sc
}

func DmsRocketMQInstancesInstanceRefCrossVpcInfoRefSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"lisenter_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"advertised_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"port_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func resourceDmsRocketMQInstancesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)

	var mErr *multierror.Error

	// getRocketmqInstances: Query DMS rocketmq instances list
	var (
		getRocketmqInstancesHttpUrl = "v2/{project_id}/instances"
		getRocketmqInstancesProduct = "dms"
	)
	getRocketmqInstancesClient, err := config.NewServiceClient(getRocketmqInstancesProduct, region)
	if err != nil {
		return diag.Errorf("error creating DmsRocketMQInstances Client: %s", err)
	}

	getRocketmqInstancesPath := getRocketmqInstancesClient.Endpoint + getRocketmqInstancesHttpUrl
	getRocketmqInstancesPath = strings.ReplaceAll(getRocketmqInstancesPath, "{project_id}", getRocketmqInstancesClient.ProjectID)

	getRocketmqInstancesQueryParams := buildGetRocketmqInstancesQueryParams(d)
	getRocketmqInstancesPath = getRocketmqInstancesPath + getRocketmqInstancesQueryParams

	getRocketmqInstancesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getRocketmqInstancesResp, err := getRocketmqInstancesClient.Request("GET", getRocketmqInstancesPath, &getRocketmqInstancesOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DmsRocketMQInstances")
	}

	getRocketmqInstancesRespBody, err := utils.FlattenResponse(getRocketmqInstancesResp)
	if err != nil {
		return diag.FromErr(err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instances", flattenGetRocketmqInstancesResponseBodyInstanceRef(getRocketmqInstancesRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetRocketmqInstancesResponseBodyInstanceRef(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("instances", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))

	crossVpcInfo := utils.PathSearch("cross_vpc_info", curJson, nil)
	var crossVpcAccess []map[string]interface{}
	if crossVpcInfo != nil {
		crossVpcAccess, _ = flattenConnectPorts(crossVpcInfo.(string))
		// if err != nil {
		// 	return err
		// }
	}

	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"name":               utils.PathSearch("name", v, nil),
			"status":             utils.PathSearch("status", v, nil),
			"description":        utils.PathSearch("description", v, nil),
			"type":               utils.PathSearch("type", v, nil),
			"specification":      utils.PathSearch("specification", v, nil),
			"engine_version":     utils.PathSearch("engine_version", v, nil),
			"vpc_id":             utils.PathSearch("vpc_id", v, nil),
			"flavor_id":          utils.PathSearch("product_id", v, nil),
			"security_group_id":  utils.PathSearch("security_group_id", v, nil),
			"subnet_id":          utils.PathSearch("subnet_id", v, nil),
			"availability_zones": utils.PathSearch("available_zones", v, nil),
			"maintain_begin":     utils.PathSearch("maintain_begin", v, nil),
			"maintain_end":       utils.PathSearch("maintain_end", v, nil),
			"storage_space":      utils.PathSearch("total_storage_space", v, nil),
			"used_storage_space": utils.PathSearch("used_storage_space", v, nil),
			"enable_publicip":    utils.PathSearch("enable_publicip", v, nil),
			"publicip_id":        utils.PathSearch("publicip_id", v, nil),
			"publicip_address":   utils.PathSearch("publicip_address", v, nil),
			"ssl_enable":         utils.PathSearch("ssl_enable", v, nil),
			// "cross_vpc_accesses":      flattenInstanceRefCrossVpcAccesses(v),
			"storage_spec_code":       utils.PathSearch("storage_spec_code", v, nil),
			"ipv6_enable":             utils.PathSearch("ipv6_enable", v, nil),
			"node_num":                utils.PathSearch("node_num", v, nil),
			"new_spec_billing_enable": utils.PathSearch("new_spec_billing_enable", v, nil),
			"enable_acl":              utils.PathSearch("enable_acl", v, nil),
			"broker_num":              utils.PathSearch("broker_num", v, nil),
			"namesrv_address":         utils.PathSearch("namesrv_address", v, nil),
			"broker_address":          utils.PathSearch("broker_address", v, nil),
			"public_namesrv_address":  utils.PathSearch("public_namesrv_address", v, nil),
			"public_broker_address":   utils.PathSearch("public_broker_address", v, nil),
			"resource_spec_code":      utils.PathSearch("resource_spec_code", v, nil),
			"cross_vpc_accesses":      crossVpcAccess,
		})
	}
	return rst
}

func buildGetRocketmqInstancesQueryParams(d *schema.ResourceData) string {
	res := ""
	res = fmt.Sprintf("%s&engine=%v", res, "reliability")

	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}

	if v, ok := d.GetOk("instance_id"); ok {
		res = fmt.Sprintf("%s&instance_id=%v", res, v)
	}

	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%v", res, v)
	}

	if v, ok := d.GetOk("exact_match_name"); ok {
		res = fmt.Sprintf("%s&instance_id=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
