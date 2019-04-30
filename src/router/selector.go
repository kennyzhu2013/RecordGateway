/*
@Time : 2019/4/18 11:02 
@Author : kenny zhu
@File : selector_byip
@Software: GoLand
@Others:
*/
package router

import (
	"github.com/micro/go-micro/selector"
	"math/rand"
	"time"
	"github.com/micro/go-micro/registry"
	"sync"
	"strconv"

	"github.com/kennyzhu/go-os/plugins/monitor"
)


var (
	// Set media-proxy tag, if any services discovered for et-cd cluster.
	serverTag = "media-proxy"
)


func init() {
	rand.Seed(time.Now().Unix())
}


//  use round select..
// server information must transfer here.
func roundBinSelect(services []*registry.Service) selector.Next {
	if len(services) == 0 {
		return func() (*registry.Node, error) {
			return nil, selector.ErrNotFound
		}
	}

	// flatten
	var nodes []*registry.Node = nil

	// Todo: must filter services by statics per call here..

	// Filter the nodes for serverTag marked by the server..
	for _, service := range services {
		for _, node := range service.Nodes {
			if node.Metadata["serverTag"] == serverTag && bHealthNodesByWeights(node) {
				nodes = append(nodes, node)
			}
		}
	}

	if len(nodes) == 0 {
		return func() (*registry.Node, error) {
			return nil, selector.ErrNotFound
		}
	}

	var i int = 0
	var mtx sync.Mutex

	// Round bin..
	return func() (*registry.Node, error) {
		mtx.Lock()
		defer mtx.Unlock()
		i++
		return nodes[i%len(nodes)], nil
	}
}

// add stats info to select.
// weights is set by clients.
// support 70%-90%, Bucket algorithm
func bHealthNodesByWeights(node *registry.Node) bool {

	// Todo: filter weights data.
	nodeMetas := node.Metadata
	timeNow := time.Now().Unix()
	if ts, ok := nodeMetas["timestamp"]; ok {
		timestamp,_ := strconv.ParseInt(ts, 10, 64)
		if timestamp + int64( monitor.HeartBeatTTL.Seconds() ) < timeNow {
			// long time no heartbeats
			return false
		}

	}

	if status, ok := nodeMetas[monitor.ServiceStatus]; ok {
		if status != monitor.NormalState{
			return false
		}
	}

	// default
	return true
}
