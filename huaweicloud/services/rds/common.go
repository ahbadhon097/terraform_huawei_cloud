package rds

import (
	"encoding/json"
	"fmt"

	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"
)

// The RDS instance is limited to only one operation at a time.
// In addition to locking and waiting between multiple operations, a retry method is required to ensure that the
// request can be executed correctly.
func handleMultiOperationsError(err error) (bool, error) {
	if err == nil {
		// The operation was executed successfully and does not need to be executed again.
		return false, nil
	}
	if errCode, ok := err.(golangsdk.ErrUnexpectedResponseCode); ok && errCode.Actual == 409 {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return false, fmt.Errorf("unmarshal the response body failed: %s", jsonErr)
		}

		// Some error codes that need to be retried coming from https://console-intl.huaweicloud.com/apiexplorer/#/errorcenter/RDS.
		retryErrCodes := map[string]struct{}{
			"DBS.201202": struct{}{},
			"DBS.200011": struct{}{},
			"DBS.200019": struct{}{},
			"DBS.200047": struct{}{},
			"DBS.200080": struct{}{},
			"DBS.201015": struct{}{},
			"DBS.201206": struct{}{},
			"DBS.212033": struct{}{},
			"DBS.280011": struct{}{},
			"DBS.280816": struct{}{},
		}

		errorCode, errorCodeErr := jmespath.Search("errCode||error_code", apiError)
		if errorCodeErr != nil {
			return false, fmt.Errorf("error parse errorCode from response body: %s", errorCodeErr)
		}

		if _, ok = retryErrCodes[errorCode.(string)]; ok {
			// The operation failed to execute and needs to be executed again, because other operations are
			// currently in progress.
			return true, err
		}
	}
	// Operation execution failed due to some resource or server issues, no need to try again.
	return false, err
}
