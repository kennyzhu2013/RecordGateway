package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"
)

type Environment struct {
	ctx context.Context
	cancel context.CancelFunc
	wg sync.WaitGroup
}

func work1(env Environment) {
	count := 1
	for {
		count++
		time.Sleep(time.Second)

		select {
		case <-env.ctx.Done():
			env.wg.Done()
			return
		default:
			fmt.Println("work1")
		}
	}
}




func work2(env Environment) {
	count := 1
	for {
		count++
		time.Sleep(time.Second)

		select {
		case <-env.ctx.Done():
			env.wg.Done()
			return
		default:
			fmt.Println("work2")
			if count > 5 {
				env.cancel()
			}
		}
	}
}



func main() {
	ctx, cancel := context.WithCancel(context.Background())
	env := Environment{ctx, cancel, sync.WaitGroup{}}
	env.wg.Add(2)
	go work1(env)
	go work2(env)
	//time.Sleep(time.Second * 5)
	//cancel()



	fmt.Println("Client start...")
	c := make(chan os.Signal)
	signal.Notify(c)
	<-c
	fmt.Println("Bye!")

}
