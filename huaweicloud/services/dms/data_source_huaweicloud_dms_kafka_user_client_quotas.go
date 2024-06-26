// Generated by PMS #204
package dms

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/filters"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
)

func DataSourceDmsKafkaUserClientQuotas() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDmsKafkaUserClientQuotasRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the instance ID.`,
			},
			"user": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the user name.`,
			},
			"client": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the client ID.`,
			},
			"quotas": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the client quotas.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the username.`,
						},
						"user_default": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether to use the default user settings.`,
						},
						"client": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the client ID.`,
						},
						"client_default": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether to use the default client settings.`,
						},
						"producer_byte_rate": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the production rate limit. The unit is byte/s.`,
						},
						"consumer_byte_rate": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the consumption rate limit. The unit is byte/s.`,
						},
					},
				},
			},
		},
	}
}

type KafkaUserClientQuotasDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newKafkaUserClientQuotasDSWrapper(d *schema.ResourceData, meta interface{}) *KafkaUserClientQuotasDSWrapper {
	return &KafkaUserClientQuotasDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceDmsKafkaUserClientQuotasRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newKafkaUserClientQuotasDSWrapper(d, meta)
	shoKafUseCliQuoRst, err := wrapper.ShowKafkaUserClientQuota()
	if err != nil {
		return diag.FromErr(err)
	}

	err = wrapper.showKafkaUserClientQuotaToSchema(shoKafUseCliQuoRst)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)
	return nil
}

// @API Kafka GET /v2/kafka/{project_id}/instances/{instance_id}/kafka-user-client-quota
func (w *KafkaUserClientQuotasDSWrapper) ShowKafkaUserClientQuota() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "dmsv2")
	if err != nil {
		return nil, err
	}

	uri := "/v2/kafka/{project_id}/instances/{instance_id}/kafka-user-client-quota"
	uri = strings.ReplaceAll(uri, "{instance_id}", w.Get("instance_id").(string))
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Filter(
			filters.New().From("quotas").
				Where("user", "contains", w.Get("user")).
				Where("client", "contains", w.Get("client")),
		).
		Request().
		Result()
}

func (w *KafkaUserClientQuotasDSWrapper) showKafkaUserClientQuotaToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("quotas", schemas.SliceToList(body.Get("quotas"),
			func(quotas gjson.Result) any {
				return map[string]any{
					"user":               quotas.Get("user").Value(),
					"user_default":       quotas.Get("user-default").Value(),
					"client":             quotas.Get("client").Value(),
					"client_default":     quotas.Get("client-default").Value(),
					"producer_byte_rate": quotas.Get("producer-byte-rate").Value(),
					"consumer_byte_rate": quotas.Get("consumer-byte-rate").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}