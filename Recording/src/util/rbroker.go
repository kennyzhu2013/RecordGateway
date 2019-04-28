/*
@Time : 2019/4/28 16:40 
@Author : kenny zhu
@File : rbroker.go
@Software: GoLand
@Others: for rabbit-msq
*/
package util

import (
	"fmt"
	"time"
	"rabbitmq"
	"github.com/micro/go-micro/broker"
	"github.com/kennyzhu/go-os/log"
)

var (
	topic = "recording.audio"
	rbroker = rabbitmq.NewBroker()
)

func InitBroker()  {

}

// for rabbit-mq
func PubMessage(body []byte)  {
	msg := &broker.Message{
		Header: map[string]string{
			"id": fmt.Sprintf("%d", i),
		},
		Body: []byte(fmt.Sprintf("%d: %s", i, time.Now().String())),
	}

	if err := rbroker.Publish(topic, msg); err != nil {
		log.Printf("[pub] failed: %v", err)
	} else {
		fmt.Println("[pub] pubbed message:", string(msg.Body))
	}

	uuid.NewUUID()
}