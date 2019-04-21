/*
@Time : 2019/4/18 15:14 
@Author : kenny zhu
@File : shard.go
@Software: GoLand
@Others:
*/
package router

import (
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/selector"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/metadata"
	"context"
	"strings"
	"strconv"
	"net/http"

	"web"
)

// use http instead...
type dispatch struct {
	key string  // "X-Media-Server" header.
	client *http.Client
}

var DefaultDispatch = &dispatch{
	key : "X-Media-Server",
	client: &http.Client{ Transport :
		web.NewRoundTripper( web.WithRegistry(registry.DefaultRegistry),
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

// call rtp-proxy with client..
// req and rsp is string here ...
func (s *dispatch) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	// get headers
	md, ok := metadata.FromContext(ctx)
	if !ok {
		// no header, defer to client
		nOpts := append(opts, client.WithSelectOption(
			selector.WithStrategy(roundBinSelect),
		))

		s.client.Post()
		return s.Client.Call(ctx, req, rsp, nOpts...)
	}

	// get key val
	val := md[s.key]

	// noop on nil value
	if len(val) == 0 {
		// no header or value, defer to client
		nOpts := append(opts, client.WithSelectOption(
			// create a selector strategy
			selector.WithStrategy(roundBinSelect),
		) )
		return s.Client.Call(ctx, req, rsp, nOpts...)
	}

	// split ip:port.
	address := strings.Split(val, ":")
	if len(address) < 2 {
		// valid ip and port, defer to client
		nOpts := append(opts, client.WithSelectOption( selector.WithStrategy(roundBinSelect) ) )
		return s.Client.Call(ctx, req, rsp, nOpts...)
	}
	proxyIp := address[0]
	proxyPort,_ := strconv.Atoi( address[1] )

	nOpts := append(opts, client.WithSelectOption(
		// create a selector strategy
		selector.WithStrategy(func(services []*registry.Service) selector.Next {
			// flatten
			var nodeResult *registry.Node

			// create the next func that always returns our node
			return func() (*registry.Node, error) {
				// Filter the nodes for serverTag marked by the server..
				for _, service := range services {
					for _, node := range service.Nodes {
						if node.Address == proxyIp && node.Port == proxyPort {
							nodeResult = node
						}
					}
				}

				if nil == nodeResult {
					return nil, selector.ErrNoneAvailable
				}

				return nodeResult, nil
			}
		}),
	))

	return s.Client.Call(ctx, req, rsp, nOpts...)
}

// NewClientWrapper is a wrapper which shards based on a header key value
func NewClientWrapper(key string) *http.Client {
	return &dispatch{
		key:    key,
		client:  http.DefaultClient,
	}
}

