package main

import (
	"fmt"
	"huaweicloud-sdk-go-v3/core/auth/basic"
	sms "huaweicloud-sdk-go-v3/services/sms/v3"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/sms/v3/model"
	region "huaweicloud-sdk-go-v3/services/sms/v3/region"
)

func main() {
	ak := "OYFWTNGHT3HG0V0MTXBZ"
	sk := "TYkhq0NqrCmOjAckdVC4fhnCW0AgGdJFAZ2q1oCQ"
	

    auth := basic.NewCredentialsBuilder().
        WithAk(ak).
        WithSk(sk).
        Build()

    client := sms.NewSmsClient(
        sms.SmsClientBuilder().
            WithRegion(region.ValueOf("ap-southeast-1")).
            WithCredential(auth).
            Build())

    request := &model.ListTemplatesRequest{}
	response, err := client.ListTemplates(request)
	if err == nil {
        fmt.Printf("%+v\n", response)
    } else {
        fmt.Println(err)
    }
}