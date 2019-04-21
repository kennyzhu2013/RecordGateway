/*
@Time : 2019/4/19 16:28 
@Author : kenny zhu
@File : roundtripper.go
@Software: GoLand
@Others:
*/
package web

import (
	"net/http"
	"errors"
	"fmt"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/selector"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/metadata"
	"strings"
	"strconv"
)

// st   selector.Strategy
type roundTripper struct {
	rt   http.RoundTripper
	opts Options
}

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

// rewrite send request and receive response.
func (r *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	s, err := r.opts.Registry.GetService(req.URL.Host)
	if err != nil {
		return nil, err
	}

	// select the one
	next := r.opts.Selector(s)

	// rudimentary retry 3 times , may be the same one.
	for i := 0; i < 3; i++ {
		n, err := next()
		if err != nil {
			continue
		}
		req.URL.Host = fmt.Sprintf("%s:%d", n.Address, n.Port)
		w, err := r.rt.RoundTrip(req)
		if err != nil {
			continue
		}
		return w, nil
	}

	return nil, errors.New("failed request")
}
