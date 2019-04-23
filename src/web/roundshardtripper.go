/*
@Time : 2019/4/19 16:28 
@Author : kenny zhu
@File : roundShardTripper.go
@Software: GoLand
@Others:
*/
package web

import (
	"net/http"
	"errors"
	"fmt"
	"github.com/micro/go-micro/registry"
	"strings"
	"strconv"
	"github.com/kennyzhu/go-os/log"
)

// st   selector.Strategy
type roundShardTripper struct {
	rt   http.RoundTripper
	opts Options
}

// rewrite send request and receive response.
// req.URL.Host = er
func (r *roundShardTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	s, err := r.opts.Registry.GetService(req.URL.Host)
	if err != nil {
		return nil, err
	}

	// get destination header.
	val := req.Header[r.opts.Destination]
	if len(val) < 1 {
		// no header, defer to client
		return r.roundTrip(s, req)
	}

	var ipSets string
	for _, temp := range  val {
		ipSets += temp
	}

	// split ip:port. with first one.
	address := strings.Split(ipSets, ":")
	if len(address) < 2 {
		// valid ip and port, defer to client
		return r.roundTrip(s, req)
	}
	proxyIp := address[0]
	proxyPort,_ := strconv.Atoi( address[1] )

	// Filter the nodes for serverTag marked by the server..
	var nodeResult *registry.Node
	for _, service := range s {
		for _, node := range service.Nodes {
			if node.Address == proxyIp && node.Port == proxyPort {
				nodeResult = node
			}
		}
	}

	if nil == nodeResult {
		// not find any service
		log.Errorf("RoundTrip not find destination service:%v", address)
		return r.roundTrip(s, req)
	}

	// need retry ?...
	req.URL.Host = fmt.Sprintf("%s:%d", proxyIp, proxyPort)
	w, err := r.rt.RoundTrip(req)
	if err != nil {
		log.Error("RoundTrip failed request")
		return nil, errors.New("RoundTrip failed request")
	}
	return w, nil
}

func (r *roundShardTripper)  roundTrip(s []*registry.Service, req *http.Request) (*http.Response, error) {
	// select the one with roundBinSelect
	next := r.opts.Selector(s)

	// rudimentary retry 3 times , may be the same one.
	for i := 0; i < 3; i++ {
		n, err := next()
		if err != nil {
			continue
		}
		if nil == n {
			log.Error("roundTrip failed not found any normal node")
			return nil, errors.New("roundTrip failed not found any normal node")
		}
		log.Infof("roundTrip found node with ip:%v, port:%v", n.Address, n.Port)
		req.URL.Host = fmt.Sprintf("%s:%d", n.Address, n.Port)

		// w, err := r.rt.RoundTrip(req)
		w, err := r.rt.RoundTrip(req)
		if err != nil {
			continue
		}
		return w, nil
	}

	log.Error("roundTrip failed request")
	return nil, errors.New("roundTrip failed request")
}