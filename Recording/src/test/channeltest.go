package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
)
type Love struct {
	wait200 chan bool
}


func main() {
	love := Love{wait200:make(chan bool)}
	go func() {
		for {
			select {
			case <-love.wait200:
				println("over")
				return
			default:
				time.Sleep(time.Second)
			}
		}
	}()

	go func() {
		time.Sleep(time.Second * 5)
		close(love.wait200)
	}()



	fmt.Println("Client start...")
	c := make(chan os.Signal)
	signal.Notify(c)
	<-c
	fmt.Println("Bye!")
}
