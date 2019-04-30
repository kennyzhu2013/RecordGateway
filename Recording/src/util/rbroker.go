/*
@Time : 2019/4/28 16:40 
@Author : kenny zhu
@File : rbroker.go
@Software: GoLand
@Others: for rabbit-msq
*/
package util

import (
	"rabbitmq"
	"github.com/micro/go-micro/broker"
	"log/log"

	"github.com/pborman/uuid"
	"config"
	"fmt"
)

var (
	pubTopic = "recording.audio"
	rbroker = rabbitmq.NewBroker()
)

func InitBroker()  error {
	// init with given url and user name.
	rbroker.Init( broker.Addrs(config.AppConf.Rabbitmq.Url) )
	pubTopic = config.AppConf.Rabbitmq.Topic

	// use one connect
	if err := rbroker.Connect(); err != nil {
		log.Fatalf("InitBroker Connect error: %v", err)
		return err
	}
	return nil
}

// for rabbit-mq, support goroutine.
func PubMessage(body []byte)  {
	// use generated id for message.
	msg := &broker.Message{
		Header: map[string]string{
			"id":  "recording.audio-" + uuid.NewUUID().String(),
		},
		Body: body,
	}

	if err := rbroker.Publish(pubTopic, msg); err != nil {
		log.Errorf("[pub] failed: %v", err)
	} else {
		log.Debugf("[pub] pubbed message:%v", string(msg.Body))
	}
}

// this is example
func subCallback(p broker.Publication) error {
	fmt.Println("[sub] received message:", string(p.Message().Body), "header", p.Message().Header)
	return nil
}

// get message from message queue, need it?
// must run in goroutine, handler is callback func
func SubMessage(subTopic string, handler broker.Handler)  {
	_, err := rbroker.Subscribe(subTopic, subCallback)
	if err != nil {
		fmt.Println(err)
	}
}