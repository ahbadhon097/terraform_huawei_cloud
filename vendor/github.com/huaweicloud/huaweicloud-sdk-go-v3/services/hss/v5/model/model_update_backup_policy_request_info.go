package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 备份策略
type UpdateBackupPolicyRequestInfo struct {

	// 策略是否启用
	Enabled *bool `json:"enabled,omitempty"`

	// 策略ID
	PolicyId *string `json:"policy_id,omitempty"`

	OperationDefinition *OperationDefinitionRequestInfo `json:"operation_definition,omitempty"`

	Trigger *BackupTriggerRequestInfo `json:"trigger,omitempty"`
}

func (o UpdateBackupPolicyRequestInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateBackupPolicyRequestInfo struct{}"
	}

	return strings.Join([]string{"UpdateBackupPolicyRequestInfo", string(data)}, " ")
}
