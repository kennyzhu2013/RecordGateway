// 端口的分配管理模块

package core

import (
	"config"
	"errors"
	"sync"
)

type PortSuite struct {
	PortRtp  int
	PortRtcp int
}

type PortSuiteHelper struct {
	startPort int
	capacity int
	portsUsed   map[int]bool
	lastAlloted int
	mu          sync.Mutex
}

var portSuiteHelper *PortSuiteHelper

func GetPortSuiteHelper() *PortSuiteHelper {
	if portSuiteHelper == nil {

		portSuiteHelper = &PortSuiteHelper{
			startPort: config.GetConfig().Ims.Ports.Start,
			capacity: (config.GetConfig().Ims.Ports.End - config.GetConfig().Ims.Ports.Start)/2,
			portsUsed: map[int]bool{},
			lastAlloted: -1,
		}
	}
	return portSuiteHelper
}

func (psHelper *PortSuiteHelper)generatePortSuiteByID(id int) *PortSuite {
	base := psHelper.startPort + id*2
	return &PortSuite{
		PortRtp:  base,
		PortRtcp: base + 1,
	}
}

func (psHelper *PortSuiteHelper) AllotPort() (*PortSuite, error) {
	psHelper.mu.Lock()
	defer psHelper.mu.Unlock()

	if len(psHelper.portsUsed) >= psHelper.capacity {
		return nil, errors.New("insufficient capacity")
	}

	for {
		psHelper.lastAlloted = (psHelper.lastAlloted + 1) % psHelper.capacity
		if _, ok := psHelper.portsUsed[psHelper.lastAlloted]; !ok {
			break
		}
		//fmt.Println("lopp")
	}
	psHelper.portsUsed[psHelper.lastAlloted] = true
	return psHelper.generatePortSuiteByID(psHelper.lastAlloted), nil
}

func (psHelper *PortSuiteHelper) ReleasePort(portRtp int) {
	psHelper.mu.Lock()
	defer psHelper.mu.Unlock()
	delete(psHelper.portsUsed, (portRtp-psHelper.startPort)/2)
}
