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
			if node.Metadata["serverTag"] == serverTag {
				nodes = append(nodes, node)
			}
		}
	}

	var i int
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

