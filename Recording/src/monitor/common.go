/*
@Time : 2019/4/23 15:48 
@Author : kenny zhu
@File : common.go
@Software: GoLand
@Others:
*/
package monitor

import "github.com/micro/go-micro/registry"

type Heartbeat struct {
	id string// The ID of the heartbeat. Will generally be the Node ID.
	service *registry.Service // The service sending the heartbeat, support multi-service
	timestamp int64 // Unix time at which this was sent
	interval int64  // The interval at which this heartbeat is expected to be sent
	ttl int64  // The time to live for this heartbeat
	weights int32 // node weight.
}