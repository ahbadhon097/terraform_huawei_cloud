package subscriptions

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

var RequestOpts golangsdk.RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// CreateOpsBuilder is used for creating subscription parameters.
// any struct providing the parameters should implement this interface
type CreateOpsBuilder interface {
	ToSubscriptionCreateMap() (map[string]interface{}, error)
}

// CreateOps is a struct that contains all the parameters.
type CreateOps struct {
	//Message endpoint
	Endpoint string `json:"endpoint" required:"true"`
	//Protocol of the message endpoint
	Protocol string `json:"protocol" required:"true"`
	//Description of the subscription
	Remark string `json:"remark,omitempty"`
	// Extension config
	Extension *ExtensionSpec `json:"extension,omitempty"`
}

type ExtensionSpec struct {
	ClientID     string `json:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
	Keyword      string `json:"keyword,omitempty"`
	SignSecret   string `json:"sign_secret,omitempty"`
}

func (ops CreateOps) ToSubscriptionCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(ops, "")
}

// Create a subscription with given parameters.
func Create(client *golangsdk.ServiceClient, ops CreateOpsBuilder, topicUrn string) (r CreateResult) {
	b, err := ops.ToSubscriptionCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(createURL(client, topicUrn), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{201, 200},
		MoreHeaders: RequestOpts.MoreHeaders,
	})

	return
}

// delete a subscription via subscription urn
func Delete(client *golangsdk.ServiceClient, subscriptionUrn string) (r DeleteResult) {
	_, r.Err = client.Delete(deleteURL(client, subscriptionUrn), &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return
}

// list all the subscriptions
func List(client *golangsdk.ServiceClient) (r ListResult) {
	pages, err := pagination.NewPager(client, listURL(client),
		func(r pagination.PageResult) pagination.Page {
			p := SubscriptionPage{pagination.OffsetPageBase{PageResult: r}}
			return p
		}).AllPages()

	if err != nil {
		r.Err = err
		return
	}
	r.Body = pages.GetBody()
	return
}

// list all the subscriptions of a topic
func ListFromTopic(client *golangsdk.ServiceClient, topicUrn string) (r ListResult) {
	pages, err := pagination.NewPager(client, listFromTopicURL(client, topicUrn),
		func(r pagination.PageResult) pagination.Page {
			p := SubscriptionPage{pagination.OffsetPageBase{PageResult: r}}
			return p
		}).AllPages()

	if err != nil {
		r.Err = err
		return
	}

	r.Body = pages.GetBody()
	return
}
