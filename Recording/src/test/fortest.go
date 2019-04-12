package main

import (
	"fmt"
	"time"
)

type Service struct {
	doneRtp  chan bool
	doneRtcp chan bool
	port     int
}

func (s *Service) StopRtp() {
	close(s.doneRtp)
}

func (s *Service) StopRtcp() {
	close(s.doneRtcp)
}

func (s *Service) StartRtp() {
	go func() {
		for {
			select {
			case <-s.doneRtp:
				fmt.Println("rtp over")
				close(s.doneRtcp)
				return
			default:
				time.Sleep(time.Second)
				fmt.Println("rtp send")
			}

		}
	}()
}

func (s *Service) StartRtcp() {
	go func() {
		for {
			select {
			case <-s.doneRtcp:
				fmt.Println("rtcp over")
				return
			default:
				time.Sleep(time.Second)
				fmt.Println("rtcp send")
			}

		}

	}()
}

func NewService() *Service {
	return &Service{make(chan bool), make(chan bool), 500}
}
func main() {
	//s := NewService()
	//s.StartRtp()
	//s.StartRtcp()
	//
	//fmt.Println("Client start...")
	//
	//time.Sleep(time.Second * 2)
	//s.StopRtp()
	//
	////time.Sleep(time.Second * 2)
	////s.StopRtcp()
	//
	//c := make(chan os.Signal)
	//signal.Notify(c)
	//<-c
	//
	//time.Sleep(time.Second * 5)

	for{
		time.Sleep(time.Second)
		println("xx")
	}
}

//params := map[string]interface{}{
//	"callid":   "callid-1",
//	"leftip":   "10.153.90.14",
//	"leftport": 5277,
//}
//_, resp, _ := core.Post("http://10.153.90.11:8080/v1/endpoint/apply", params)
//fmt.Println(resp.Get("test").MustInt())
