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
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func ResourceDisStreamV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceDisStreamV2Create,
		Read:   resourceDisStreamV2Read,
		Update: resourceDisStreamV2Update,
		Delete: resourceDisStreamV2Delete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"partition_count": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"stream_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"auto_scale_max_partition_count": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},

			"auto_scale_min_partition_count": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},

			"compression_format": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},

			"csv_delimiter": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},

			"data_schema": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},

			"data_type": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},

			"retention_period": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},

			"stream_type": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},

			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
							ForceNew: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},

			"created": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"readable_partition_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"writable_partition_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceDisStreamV2UserInputParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"terraform_resource_data":        d,
		"auto_scale_max_partition_count": d.Get("auto_scale_max_partition_count"),
		"auto_scale_min_partition_count": d.Get("auto_scale_min_partition_count"),
		"compression_format":             d.Get("compression_format"),
		"csv_delimiter":                  d.Get("csv_delimiter"),
		"data_schema":                    d.Get("data_schema"),
		"data_type":                      d.Get("data_type"),
		"partition_count":                d.Get("partition_count"),
		"retention_period":               d.Get("retention_period"),
		"stream_name":                    d.Get("stream_name"),
		"stream_type":                    d.Get("stream_type"),
		"tags":                           d.Get("tags"),
	}
}

func resourceDisStreamV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.DisV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating sdk client, err=%s", err)
	}

	opts := resourceDisStreamV2UserInputParams(d)

	params, err := buildDisStreamV2CreateParameters(opts, nil)
	if err != nil {
		return fmt.Errorf("Error building the request body of api(create), err=%s", err)
	}
	_, err = sendDisStreamV2CreateRequest(d, params, client)
	if err != nil {
		return fmt.Errorf("Error creating DisStreamV2, err=%s", err)
	}

	d.SetId(opts["stream_name"].(string))

	return resourceDisStreamV2Read(d, meta)
}

func resourceDisStreamV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.DisV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating sdk client, err=%s", err)
	}

	res := make(map[string]interface{})

	v, err := sendDisStreamV2ReadRequest(d, client)
	if err != nil {
		return err
	}
	res["read"] = fillDisStreamV2ReadRespBody(v)

	return setDisStreamV2Properties(d, res)
}

func resourceDisStreamV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.DisV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating sdk client, err=%s", err)
	}

	opts := resourceDisStreamV2UserInputParams(d)

	params, err := buildDisStreamV2UpdateParameters(opts, nil)
	if err != nil {
		return fmt.Errorf("Error building the request body of api(update), err=%s", err)
	}
	if e, _ := isEmptyValue(reflect.ValueOf(params)); !e {
		_, err = sendDisStreamV2UpdateRequest(d, params, client)
		if err != nil {
			return fmt.Errorf("Error updating (DisStreamV2: %v), err=%s", d.Id(), err)
		}
	}

	return resourceDisStreamV2Read(d, meta)
}

func resourceDisStreamV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.DisV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating sdk client, err=%s", err)
	}

	url, err := replaceVars(d, "streams/{id}", nil)
	if err != nil {
		return err
	}
	url = client.ServiceURL(url)

	log.Printf("[DEBUG] Deleting Stream %q", d.Id())
	r := golangsdk.Result{}
	_, r.Err = client.Delete(url, &golangsdk.RequestOpts{
		OkCodes:      successHTTPCodes,
		JSONBody:     nil,
		JSONResponse: nil,
		MoreHeaders:  map[string]string{"Content-Type": "application/json"},
	})
	if r.Err != nil {
		return fmt.Errorf("Error deleting Stream %q, err=%s", d.Id(), r.Err)
	}

	return nil
}

func buildDisStreamV2CreateParameters(opts map[string]interface{}, arrayIndex map[string]int) (interface{}, error) {
	params := make(map[string]interface{})

	v, err := expandDisStreamV2CreateAutoCaleEnable(opts, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		params["auto_scale_enabled"] = v
	}

	v, err = navigateValue(opts, []string{"auto_scale_max_partition_count"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		params["auto_scale_max_partition_count"] = v
	}

	v, err = navigateValue(opts, []string{"auto_scale_min_partition_count"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		params["auto_scale_min_partition_count"] = v
	}

	v, err = navigateValue(opts, []string{"compression_format"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		params["compression_format"] = v
	}

	v, err = expandDisStreamV2CreateCsvProperties(opts, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		params["csv_properties"] = v
	}

	v, err = navigateValue(opts, []string{"retention_period"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		params["data_duration"] = v
	}

	v, err = navigateValue(opts, []string{"data_schema"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		params["data_schema"] = v
	}

	v, err = navigateValue(opts, []string{"data_type"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		params["data_type"] = v
	}

	v, err = navigateValue(opts, []string{"partition_count"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		params["partition_count"] = v
	}

	v, err = navigateValue(opts, []string{"stream_name"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		params["stream_name"] = v
	}

	v, err = navigateValue(opts, []string{"stream_type"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		params["stream_type"] = v
	}

	v, err = expandDisStreamV2CreateTags(opts, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		params["tags"] = v
	}

	return params, nil
}

func expandDisStreamV2CreateCsvProperties(d interface{}, arrayIndex map[string]int) (interface{}, error) {
	req := make(map[string]interface{})

	v, err := navigateValue(d, []string{"csv_delimiter"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		req["delimiter"] = v
	}

	return req, nil
}

func expandDisStreamV2CreateTags(d interface{}, arrayIndex map[string]int) (interface{}, error) {
	newArrayIndex := make(map[string]int)
	if arrayIndex != nil {
		for k, v := range arrayIndex {
			newArrayIndex[k] = v
		}
	}

	val, err := navigateValue(d, []string{"tags"}, newArrayIndex)
	if err != nil {
		return nil, err
	}
	n := 0
	if val1, ok := val.([]interface{}); ok && len(val1) > 0 {
		n = len(val1)
	} else {
		return nil, nil
	}
	req := make([]interface{}, 0, n)
	for i := 0; i < n; i++ {
		newArrayIndex["tags"] = i
		transformed := make(map[string]interface{})

		v, err := navigateValue(d, []string{"tags", "key"}, newArrayIndex)
		if err != nil {
			return nil, err
		}
		if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
			return nil, err
		} else if !e {
			transformed["key"] = v
		}

		v, err = navigateValue(d, []string{"tags", "value"}, newArrayIndex)
		if err != nil {
			return nil, err
		}
		if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
			return nil, err
		} else if !e {
			transformed["value"] = v
		}

		if len(transformed) > 0 {
			req = append(req, transformed)
		}
	}

	return req, nil
}

func sendDisStreamV2CreateRequest(d *schema.ResourceData, params interface{},
	client *golangsdk.ServiceClient) (interface{}, error) {
	url := client.ServiceURL("streams")

	r := golangsdk.Result{}
	_, r.Err = client.Post(url, params, nil, &golangsdk.RequestOpts{
		OkCodes: successHTTPCodes,
	})
	if r.Err != nil {
		return nil, fmt.Errorf("Error running api(create), err=%s", r.Err)
	}
	return r.Body, nil
}

func buildDisStreamV2UpdateParameters(opts map[string]interface{}, arrayIndex map[string]int) (interface{}, error) {
	params := make(map[string]interface{})

	v, err := navigateValue(opts, []string{"stream_name"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		params["stream_name"] = v
	}

	v, err = navigateValue(opts, []string{"partition_count"}, arrayIndex)
	if err != nil {
		return nil, err
	}
	if e, err := isEmptyValue(reflect.ValueOf(v)); err != nil {
		return nil, err
	} else if !e {
		params["target_partition_count"] = v
	}

	return params, nil
}

func sendDisStreamV2UpdateRequest(d *schema.ResourceData, params interface{},
	client *golangsdk.ServiceClient) (interface{}, error) {
	url, err := replaceVars(d, "streams/{id}", nil)
	if err != nil {
		return nil, err
	}
	url = client.ServiceURL(url)

	r := golangsdk.Result{}
	_, r.Err = client.Put(url, params, nil, &golangsdk.RequestOpts{
		OkCodes: successHTTPCodes,
	})
	if r.Err != nil {
		return nil, fmt.Errorf("Error running api(update), err=%s", r.Err)
	}
	return r.Body, nil
}

func sendDisStreamV2ReadRequest(d *schema.ResourceData, client *golangsdk.ServiceClient) (interface{}, error) {
	url, err := replaceVars(d, "streams/{id}", nil)
	if err != nil {
		return nil, err
	}
	url = client.ServiceURL(url)

	r := golangsdk.Result{}
	_, r.Err = client.Get(url, &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: map[string]string{"Content-Type": "application/json"}})
	if r.Err != nil {
		return nil, fmt.Errorf("Error running api(read) for resource(DisStreamV2), err=%s", r.Err)
	}

	return r.Body, nil
}

func fillDisStreamV2ReadRespBody(body interface{}) interface{} {
	result := make(map[string]interface{})
	val, ok := body.(map[string]interface{})
	if !ok {
		val = make(map[string]interface{})
	}

	if v, ok := val["auto_scale_max_partition_count"]; ok {
		result["auto_scale_max_partition_count"] = v
	} else {
		result["auto_scale_max_partition_count"] = nil
	}

	if v, ok := val["auto_scale_min_partition_count"]; ok {
		result["auto_scale_min_partition_count"] = v
	} else {
		result["auto_scale_min_partition_count"] = nil
	}

	if v, ok := val["compression_format"]; ok {
		result["compression_format"] = v
	} else {
		result["compression_format"] = nil
	}

	if v, ok := val["create_time"]; ok {
		result["create_time"] = v
	} else {
		result["create_time"] = nil
	}

	if v, ok := val["csv_properties"]; ok {
		result["csv_properties"] = fillDisStreamV2ReadRespCsvProperties(v)
	} else {
		result["csv_properties"] = nil
	}

	if v, ok := val["data_schema"]; ok {
		result["data_schema"] = v
	} else {
		result["data_schema"] = nil
	}

	if v, ok := val["data_type"]; ok {
		result["data_type"] = v
	} else {
		result["data_type"] = nil
	}

	if v, ok := val["readable_partition_count"]; ok {
		result["readable_partition_count"] = v
	} else {
		result["readable_partition_count"] = nil
	}

	if v, ok := val["retention_period"]; ok {
		result["retention_period"] = v
	} else {
		result["retention_period"] = nil
	}

	if v, ok := val["stream_name"]; ok {
		result["stream_name"] = v
	} else {
		result["stream_name"] = nil
	}

	if v, ok := val["stream_type"]; ok {
		result["stream_type"] = v
	} else {
		result["stream_type"] = nil
	}

	if v, ok := val["tags"]; ok {
		result["tags"] = fillDisStreamV2ReadRespTags(v)
	} else {
		result["tags"] = nil
	}

	if v, ok := val["writable_partition_count"]; ok {
		result["writable_partition_count"] = v
	} else {
		result["writable_partition_count"] = nil
	}

	return result
}

func fillDisStreamV2ReadRespCsvProperties(value interface{}) interface{} {
	if value == nil {
		return nil
	}

	value1, ok := value.(map[string]interface{})
	if !ok {
		value1 = make(map[string]interface{})
	}
	result := make(map[string]interface{})

	if v, ok := value1["delimiter"]; ok {
		result["delimiter"] = v
	} else {
		result["delimiter"] = nil
	}

	return result
}

func fillDisStreamV2ReadRespTags(value interface{}) interface{} {
	if value == nil {
		return nil
	}

	value1, ok := value.([]interface{})
	if !ok || len(value1) == 0 {
		return nil
	}

	n := len(value1)
	result := make([]interface{}, n, n)
	for i := 0; i < n; i++ {
		val := make(map[string]interface{})
		item := value1[i].(map[string]interface{})

		if v, ok := item["key"]; ok {
			val["key"] = v
		} else {
			val["key"] = nil
		}

		if v, ok := item["value"]; ok {
			val["value"] = v
		} else {
			val["value"] = nil
		}

		result[i] = val
	}

	return result
}

func setDisStreamV2Properties(d *schema.ResourceData, response map[string]interface{}) error {
	opts := resourceDisStreamV2UserInputParams(d)

	v, err := navigateValue(response, []string{"read", "auto_scale_max_partition_count"}, nil)
	if err != nil {
		return fmt.Errorf("Error reading Stream:auto_scale_max_partition_count, err: %s", err)
	}
	if err = d.Set("auto_scale_max_partition_count", v); err != nil {
		return fmt.Errorf("Error setting Stream:auto_scale_max_partition_count, err: %s", err)
	}

	v, err = navigateValue(response, []string{"read", "auto_scale_min_partition_count"}, nil)
	if err != nil {
		return fmt.Errorf("Error reading Stream:auto_scale_min_partition_count, err: %s", err)
	}
	if err = d.Set("auto_scale_min_partition_count", v); err != nil {
		return fmt.Errorf("Error setting Stream:auto_scale_min_partition_count, err: %s", err)
	}

	v, err = navigateValue(response, []string{"read", "compression_format"}, nil)
	if err != nil {
		return fmt.Errorf("Error reading Stream:compression_format, err: %s", err)
	}
	if err = d.Set("compression_format", v); err != nil {
		return fmt.Errorf("Error setting Stream:compression_format, err: %s", err)
	}

	v, err = navigateValue(response, []string{"read", "create_time"}, nil)
	if err != nil {
		return fmt.Errorf("Error reading Stream:created, err: %s", err)
	}
	if err = d.Set("created", v); err != nil {
		return fmt.Errorf("Error setting Stream:created, err: %s", err)
	}

	v, err = navigateValue(response, []string{"read", "csv_properties", "delimiter"}, nil)
	if err != nil {
		return fmt.Errorf("Error reading Stream:csv_delimiter, err: %s", err)
	}
	if err = d.Set("csv_delimiter", v); err != nil {
		return fmt.Errorf("Error setting Stream:csv_delimiter, err: %s", err)
	}

	v, err = navigateValue(response, []string{"read", "data_schema"}, nil)
	if err != nil {
		return fmt.Errorf("Error reading Stream:data_schema, err: %s", err)
	}
	if err = d.Set("data_schema", v); err != nil {
		return fmt.Errorf("Error setting Stream:data_schema, err: %s", err)
	}

	v, err = navigateValue(response, []string{"read", "data_type"}, nil)
	if err != nil {
		return fmt.Errorf("Error reading Stream:data_type, err: %s", err)
	}
	if err = d.Set("data_type", v); err != nil {
		return fmt.Errorf("Error setting Stream:data_type, err: %s", err)
	}

	v, err = navigateValue(response, []string{"read", "readable_partition_count"}, nil)
	if err != nil {
		return fmt.Errorf("Error reading Stream:readable_partition_count, err: %s", err)
	}
	if err = d.Set("readable_partition_count", v); err != nil {
		return fmt.Errorf("Error setting Stream:readable_partition_count, err: %s", err)
	}

	v, err = navigateValue(response, []string{"read", "retention_period"}, nil)
	if err != nil {
		return fmt.Errorf("Error reading Stream:retention_period, err: %s", err)
	}
	if err = d.Set("retention_period", v); err != nil {
		return fmt.Errorf("Error setting Stream:retention_period, err: %s", err)
	}

	v, err = navigateValue(response, []string{"read", "stream_name"}, nil)
	if err != nil {
		return fmt.Errorf("Error reading Stream:stream_name, err: %s", err)
	}
	if err = d.Set("stream_name", v); err != nil {
		return fmt.Errorf("Error setting Stream:stream_name, err: %s", err)
	}

	v, err = navigateValue(response, []string{"read", "stream_type"}, nil)
	if err != nil {
		return fmt.Errorf("Error reading Stream:stream_type, err: %s", err)
	}
	if err = d.Set("stream_type", v); err != nil {
		return fmt.Errorf("Error setting Stream:stream_type, err: %s", err)
	}

	v, _ = opts["tags"]
	v, err = flattenDisStreamV2Tags(response, nil, v)
	if err != nil {
		return fmt.Errorf("Error reading Stream:tags, err: %s", err)
	}
	if err = d.Set("tags", v); err != nil {
		return fmt.Errorf("Error setting Stream:tags, err: %s", err)
	}

	v, err = navigateValue(response, []string{"read", "writable_partition_count"}, nil)
	if err != nil {
		return fmt.Errorf("Error reading Stream:writable_partition_count, err: %s", err)
	}
	if err = d.Set("writable_partition_count", v); err != nil {
		return fmt.Errorf("Error setting Stream:writable_partition_count, err: %s", err)
	}

	return nil
}

func flattenDisStreamV2Tags(d interface{}, arrayIndex map[string]int, currentValue interface{}) (interface{}, error) {
	n := 0
	hasInitValue := true
	result, ok := currentValue.([]interface{})
	if !ok || len(result) == 0 {
		v, err := navigateValue(d, []string{"read", "tags"}, arrayIndex)
		if err != nil {
			return nil, err
		}
		if v1, ok := v.([]interface{}); ok && len(v1) > 0 {
			n = len(v1)
		} else {
			return currentValue, nil
		}
		result = make([]interface{}, 0, n)
		hasInitValue = false
	} else {
		n = len(result)
	}

	newArrayIndex := make(map[string]int)
	if arrayIndex != nil {
		for k, v := range arrayIndex {
			newArrayIndex[k] = v
		}
	}

	for i := 0; i < n; i++ {
		newArrayIndex["read.tags"] = i

		var r map[string]interface{}
		if len(result) >= (i+1) && result[i] != nil {
			r = result[i].(map[string]interface{})
		} else {
			r = make(map[string]interface{})
		}

		v, err := navigateValue(d, []string{"read", "tags", "key"}, newArrayIndex)
		if err != nil {
			return nil, fmt.Errorf("Error reading Stream:key, err: %s", err)
		}
		r["key"] = v

		v, err = navigateValue(d, []string{"read", "tags", "value"}, newArrayIndex)
		if err != nil {
			return nil, fmt.Errorf("Error reading Stream:value, err: %s", err)
		}
		r["value"] = v

		if len(result) >= (i + 1) {
			if result[i] == nil {
				result[i] = r
			}
		} else {
			for _, v := range r {
				if v != nil {
					result = append(result, r)
					break
				}
			}
		}
	}

	if hasInitValue || len(result) > 0 {
		return result, nil
	}
	return currentValue, nil
}
