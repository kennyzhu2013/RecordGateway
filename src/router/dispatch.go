/*
@Time : 2019/4/18 15:14 
@Author : kenny zhu
@File : shard.go
@Software: GoLand
@Others:
*/
package router

import (
	"net/http"

	"web"
	"io"
	"github.com/kennyzhu/go-os/plugins/etcdv3"
)

type WebClient interface {
	// set destination
	Post(serviceUrl string, contentType string, destination string, body io.Reader) (resp *http.Response, err error)

	// for test
	TestGet(serviceUrl string) (resp *http.Response, err error)
}

// use http instead...
type dispatch struct {
	key string  // "X-Media-Server" header.
	client *http.Client
}

var DefaultRouter = &dispatch{
	key : "X-Media-Server",
	client: &http.Client{ Transport :
		web.NewRoundShardTripper( web.WithRegistry(etcdv3.DefaultEtcdRegistry),
			web.WithSelector(roundBinSelect) )},
}


// no wrap client..
/*
func (s *dispatch) DefaultClient() http.Client {
	rt := web.NewRoundTripper(
		web.WithRegistry(registry.DefaultRegistry),
		web.WithSelector(roundBinSelect),
	)

	s.client.Transport = rt
	return s.client
}*/

func (c *dispatch) Post(serviceUrl string, contentType string, destination string, body io.Reader) (resp *http.Response, err error) {
	// http://go.micro.api.iris/Preferences/GetPreferencesList?limit=2&index=1...
	httpServiceUrl := "http://" + serviceUrl
	req, err := http.NewRequest("POST", httpServiceUrl, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set(c.key, destination)
	return c.client.Do(req)
}

func (c *dispatch) TestGet(serviceUrl string)(resp *http.Response, err error)  {
	httpServiceUrl := "http://" + serviceUrl
	return c.client.Get(httpServiceUrl)
}

// NewClientWrapper is a wrapper which shards based on a header key value
func NewHttpRouter(key string, httpClient *http.Client) WebClient {
	return &dispatch{
		key:    key,
		client:  httpClient,
	}
}

