/*
@Time : 2019/4/19 16:23 
@Author : kenny zhu
@File : web.go
@Software: GoLand
@Others:
*/
package web

import (
	"net/http"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/selector"
)

// https://stackoverflow.com/questions/26707941/go-roundtripper-and-transport.
func NewRoundShardTripper(opts ...Option) http.RoundTripper {
	options := Options{
		Registry: registry.DefaultRegistry,
		Selector: selector.Random,
		Destination: "X-Media-Server",
	}
	for _, o := range opts {
		o(&options)
	}

	return &roundShardTripper{
		rt:   http.DefaultTransport,
		opts: options,
	}
}