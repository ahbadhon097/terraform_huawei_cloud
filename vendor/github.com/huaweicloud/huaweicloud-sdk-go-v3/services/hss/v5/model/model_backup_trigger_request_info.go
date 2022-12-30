package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 策略时间调度规则
type BackupTriggerRequestInfo struct {
	Properties *BackupTriggerPropertiesRequestInfo `json:"properties,omitempty"`
}

func (o BackupTriggerRequestInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BackupTriggerRequestInfo struct{}"
	}

	return strings.Join([]string{"BackupTriggerRequestInfo", string(data)}, " ")
}
