package cce

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/cce/v3/nodes"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func resourceNodeExtendParamsSchema(conflictList []string) *schema.Schema {
	return &schema.Schema{
		Type:          schema.TypeList,
		Optional:      true,
		ForceNew:      true,
		MaxItems:      1,
		ConflictsWith: conflictList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"max_pods": {
					Type:     schema.TypeInt,
					Optional: true,
					ForceNew: true,
				},
				"docker_base_size": {
					Type:     schema.TypeInt,
					Optional: true,
					ForceNew: true,
				},
				"preinstall": {
					Type:      schema.TypeString,
					Optional:  true,
					ForceNew:  true,
					StateFunc: utils.DecodeHashAndHexEncode,
				},
				"postinstall": {
					Type:      schema.TypeString,
					Optional:  true,
					ForceNew:  true,
					StateFunc: utils.DecodeHashAndHexEncode,
				},
				"node_image_id": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
				},
				"node_multi_queue": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
				},
				"nic_threshold": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
				},
				"agency_name": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
				},
				"kube_reserved_mem": {
					Type:     schema.TypeInt,
					Optional: true,
					ForceNew: true,
				},
				"system_reserved_mem": {
					Type:     schema.TypeInt,
					Optional: true,
					ForceNew: true,
				},
			},
		},
	}
}

func buildResourceNodeExtendParam(d *schema.ResourceData) map[string]interface{} {
	extendParam := make(map[string]interface{})
	if v, ok := d.GetOk("extend_param"); ok {
		for key, val := range v.(map[string]interface{}) {
			extendParam[key] = val.(string)
		}
		if v, ok := extendParam["periodNum"]; ok {
			periodNum, err := strconv.Atoi(v.(string))
			if err != nil {
				log.Printf("[WARNING] PeriodNum %s invalid, Type conversion error: %s", v.(string), err)
			}
			extendParam["periodNum"] = periodNum
		}
	}

	if v, ok := d.GetOk("ecs_performance_type"); ok {
		extendParam["ecs:performancetype"] = v.(string)
	}
	if v, ok := d.GetOk("max_pods"); ok {
		extendParam["maxPods"] = v.(int)
	}
	if v, ok := d.GetOk("order_id"); ok {
		extendParam["orderID"] = v.(string)
	}
	if v, ok := d.GetOk("product_id"); ok {
		extendParam["productID"] = v.(string)
	}
	if v, ok := d.GetOk("public_key"); ok {
		extendParam["publicKey"] = v.(string)
	}
	if v, ok := d.GetOk("preinstall"); ok {
		extendParam["alpha.cce/preInstall"] = utils.TryBase64EncodeString(v.(string))
	}
	if v, ok := d.GetOk("postinstall"); ok {
		extendParam["alpha.cce/postInstall"] = utils.TryBase64EncodeString(v.(string))
	}

	return extendParam
}

func buildResourceNodeExtendParams(extendParamsRaw []interface{}) map[string]interface{} {
	if len(extendParamsRaw) != 1 {
		return nil
	}

	if extendParams, ok := extendParamsRaw[0].(map[string]interface{}); ok {
		res := map[string]interface{}{
			"maxPods":               utils.ValueIngoreEmpty(extendParams["max_pods"]),
			"dockerBaseSize":        utils.ValueIngoreEmpty(extendParams["docker_base_size"]),
			"alpha.cce/preInstall":  utils.ValueIngoreEmpty(utils.TryBase64EncodeString(extendParams["preinstall"].(string))),
			"alpha.cce/postInstall": utils.ValueIngoreEmpty(utils.TryBase64EncodeString(extendParams["postinstall"].(string))),
			"alpha.cce/NodeImageID": utils.ValueIngoreEmpty(extendParams["node_image_id"]),
			"nicMultiqueue":         utils.ValueIngoreEmpty(extendParams["node_multi_queue"]),
			"nicThreshold":          utils.ValueIngoreEmpty(extendParams["nic_threshold"]),
			"agency_name":           utils.ValueIngoreEmpty(extendParams["agency_name"]),
			"kube-reserved-mem":     utils.ValueIngoreEmpty(extendParams["kube_reserved_mem"]),
			"system-reserved-mem":   utils.ValueIngoreEmpty(extendParams["system_reserved_mem"]),
		}

		return res
	}

	return nil
}

func buildExtendParams(d *schema.ResourceData) map[string]interface{} {
	res := make(map[string]interface{})
	extendParam := buildResourceNodeExtendParam(d)
	extendParams := buildResourceNodeExtendParams(d.Get("extend_params").([]interface{}))

	// defaults to use extend_params
	if len(extendParam) != 0 {
		for k, v := range extendParam {
			res[k] = v
		}
	} else {
		for k, v := range extendParams {
			res[k] = v
		}
	}

	// assemble the charge info
	var isPrePaid bool
	var billingMode int

	if v, ok := d.GetOk("charging_mode"); ok && v.(string) == "prePaid" {
		isPrePaid = true
	}
	if v, ok := d.GetOk("billing_mode"); ok {
		billingMode = v.(int)
	}
	if isPrePaid || billingMode == 1 {
		res["chargingMode"] = 1
		res["isAutoRenew"] = "false"
		res["isAutoPay"] = common.GetAutoPay(d)
	}

	if v, ok := d.GetOk("period_unit"); ok {
		res["periodType"] = v.(string)
	}
	if v, ok := d.GetOk("period"); ok {
		res["periodNum"] = v.(int)
	}
	if v, ok := d.GetOk("auto_renew"); ok {
		res["isAutoRenew"] = v.(string)
	}

	return utils.RemoveNil(res)
}

func resourceNodeRootVolume() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"size": {
					Type:     schema.TypeInt,
					Required: true,
					ForceNew: true,
				},
				"volumetype": {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
				},
				"extend_params": {
					Type:     schema.TypeMap,
					Optional: true,
					ForceNew: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
				"kms_key_id": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
				"dss_pool_id": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},

				// Internal parameters
				"hw_passthrough": {
					Type:        schema.TypeBool,
					Optional:    true,
					ForceNew:    true,
					Description: "schema: Internal",
				},

				// Deprecated parameters
				"extend_param": {
					Type:       schema.TypeString,
					Optional:   true,
					ForceNew:   true,
					Deprecated: "use extend_params instead",
				},
			},
		},
	}
}

func resourceNodeDataVolume() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		ForceNew: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"size": {
					Type:     schema.TypeInt,
					Required: true,
					ForceNew: true,
				},
				"volumetype": {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
				},
				"extend_params": {
					Type:     schema.TypeMap,
					Optional: true,
					ForceNew: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
				"kms_key_id": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
				"dss_pool_id": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},

				// Internal parameters
				"hw_passthrough": {
					Type:        schema.TypeBool,
					Optional:    true,
					ForceNew:    true,
					Description: "schema: Internal",
				},

				// Deprecated parameters
				"extend_param": {
					Type:       schema.TypeString,
					Optional:   true,
					ForceNew:   true,
					Deprecated: "use extend_params instead",
				},
			},
		},
	}
}

func buildResourceNodeRootVolume(d *schema.ResourceData) nodes.VolumeSpec {
	var root nodes.VolumeSpec
	volumeRaw := d.Get("root_volume").([]interface{})
	if len(volumeRaw) == 1 {
		rawMap := volumeRaw[0].(map[string]interface{})
		root.Size = rawMap["size"].(int)
		root.VolumeType = rawMap["volumetype"].(string)
		root.HwPassthrough = rawMap["hw_passthrough"].(bool)
		root.ExtendParam = rawMap["extend_params"].(map[string]interface{})

		if rawMap["kms_key_id"].(string) != "" {
			metadata := nodes.VolumeMetadata{
				SystemEncrypted: "1",
				SystemCmkid:     rawMap["kms_key_id"].(string),
			}
			root.Metadata = &metadata
		}

		if rawMap["dss_pool_id"].(string) != "" {
			root.ClusterID = rawMap["dss_pool_id"].(string)
			root.ClusterType = "dss"
		}
	}

	return root
}

func buildResourceNodeDataVolume(d *schema.ResourceData) []nodes.VolumeSpec {
	volumeRaw := d.Get("data_volumes").([]interface{})
	volumes := make([]nodes.VolumeSpec, len(volumeRaw))
	for i, raw := range volumeRaw {
		rawMap := raw.(map[string]interface{})
		volumes[i] = nodes.VolumeSpec{
			Size:          rawMap["size"].(int),
			VolumeType:    rawMap["volumetype"].(string),
			HwPassthrough: rawMap["hw_passthrough"].(bool),
			ExtendParam:   rawMap["extend_params"].(map[string]interface{}),
		}
		if rawMap["kms_key_id"].(string) != "" {
			metadata := nodes.VolumeMetadata{
				SystemEncrypted: "1",
				SystemCmkid:     rawMap["kms_key_id"].(string),
			}
			volumes[i].Metadata = &metadata
		}

		if rawMap["dss_pool_id"].(string) != "" {
			volumes[i].ClusterID = rawMap["dss_pool_id"].(string)
			volumes[i].ClusterType = "dss"
		}
	}
	return volumes
}

func flattenResourceNodeRootVolume(rootVolume nodes.VolumeSpec) []map[string]interface{} {
	res := []map[string]interface{}{
		{
			"size":           rootVolume.Size,
			"volumetype":     rootVolume.VolumeType,
			"hw_passthrough": rootVolume.HwPassthrough,
			"extend_params":  rootVolume.ExtendParam,
			"extend_param":   "",
			"dss_pool_id":    rootVolume.ClusterID,
		},
	}
	if rootVolume.Metadata != nil {
		res[0]["kms_key_id"] = rootVolume.Metadata.SystemCmkid
	}

	return res
}

func flattenResourceNodeDataVolume(dataVolumes []nodes.VolumeSpec) []map[string]interface{} {
	if len(dataVolumes) == 0 {
		return nil
	}

	res := make([]map[string]interface{}, len(dataVolumes))
	for i, v := range dataVolumes {
		res[i] = map[string]interface{}{
			"size":           v.Size,
			"volumetype":     v.VolumeType,
			"hw_passthrough": v.HwPassthrough,
			"extend_params":  v.ExtendParam,
			"extend_param":   "",
			"dss_pool_id":    v.ClusterID,
		}

		if v.Metadata != nil {
			res[i]["kms_key_id"] = v.Metadata.SystemCmkid
		}
	}

	return res
}
