// Generated by PMS #289
package dew

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceDewKpsRunningTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDewKpsRunningTasksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"tasks": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of the running tasks.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the task.`,
						},
						"server_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the instance associated with the task.`,
						},
						"server_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the instance associated with the task.`,
						},
						"operate_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The operation type of the task.`,
						},
						"keypair_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the keypair associated with the task.`,
						},
						"task_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The start time of the task, in RFC3339 format.`,
						},
					},
				},
			},
		},
	}
}

type KpsRunningTasksDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newKpsRunningTasksDSWrapper(d *schema.ResourceData, meta interface{}) *KpsRunningTasksDSWrapper {
	return &KpsRunningTasksDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceDewKpsRunningTasksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newKpsRunningTasksDSWrapper(d, meta)
	listRunningTaskRst, err := wrapper.ListRunningTask()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listRunningTaskToSchema(listRunningTaskRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API KPS GET /v3/{project_id}/running-tasks
func (w *KpsRunningTasksDSWrapper) ListRunningTask() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "kms")
	if err != nil {
		return nil, err
	}

	uri := "/v3/{project_id}/running-tasks"
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		OffsetPager("tasks", "offset", "limit", 0).
		Request().
		Result()
}

func (w *KpsRunningTasksDSWrapper) listRunningTaskToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("tasks", schemas.SliceToList(body.Get("tasks"),
			func(tasks gjson.Result) any {
				return map[string]any{
					"id":           tasks.Get("task_id").Value(),
					"server_id":    tasks.Get("server_id").Value(),
					"server_name":  tasks.Get("server_name").Value(),
					"operate_type": tasks.Get("operate_type").Value(),
					"keypair_name": tasks.Get("keypair_name").Value(),
					"task_time":    w.setTasksTaskTime(tasks),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}

func (*KpsRunningTasksDSWrapper) setTasksTaskTime(data gjson.Result) string {
	return utils.FormatTimeStampRFC3339(convertStrToInt(data.Get("task_time").String())/1000, false)
}