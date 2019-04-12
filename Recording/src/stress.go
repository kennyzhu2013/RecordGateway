package main

import (
	"core"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"
)

const startport_left = 5000
const pps = 50 //每秒50个包
const sleeptime = time.Millisecond * 1000 / pps

var portsUsed = map[int]bool{} //最大容量为threads数量
var lastAlloted int = -1

var mu sync.Mutex

func selectPort(capacity int) (int, error) {
	mu.Lock()
	defer mu.Unlock()

	if len(portsUsed) >= capacity {
		return -1, errors.New("insufficient capacity")
	}

	for {
		lastAlloted = (lastAlloted + 1) % capacity
		if _, ok := portsUsed[lastAlloted]; !ok {
			break
		}
	}
	portsUsed[lastAlloted] = true
	return lastAlloted + startport_left, nil
}
func releasePort(port int) {
	mu.Lock()
	defer mu.Unlock()

	delete(portsUsed, port-startport_left)
}

func generateMsg() string {
	return fmt.Sprintf("%daaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", time.Now().UnixNano())
}

func stringTotime(s string) (time.Time) {
	sec, _ := strconv.ParseInt(s[0:10], 10, 64)
	nano, _ := strconv.ParseInt(s[10:], 10, 64)
	return time.Unix(sec, nano)
}

func SendLoop(callid string, ipLocal string, portLocal int, ipRemote string, portRemote int, c chan int) {
	laddr, _ := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ipLocal, portLocal))
	raddr, _ := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ipRemote, portRemote))

	conn, _ := net.DialUDP("udp", laddr, raddr)
	defer conn.Close()


	//早期媒体5秒
	for count := 5 * pps; count > 0; count-- {
		//conn.Write([]byte(fmt.Sprintf("%d xxxxxxxx %d", time.Now().UnixNano(), portLocal)))
		conn.Write([]byte(generateMsg()))
		time.Sleep(sleeptime)
	}

	core.Post(fmt.Sprintf("http://%s:8080/v1/endpoint/200", ipRemote), map[string]interface{}{
		"callid": callid,
	})
	//随机通话：15-45秒
	//for count := (15 + rand.Intn(30)) * pps; count > 0; count-- {
	for count := 5 * pps; count > 0; count-- {
		//conn.Write([]byte(fmt.Sprintf("%d xxxxxxxx %d", time.Now().UnixNano(), portLocal)))
		conn.Write([]byte(generateMsg()))
		time.Sleep(sleeptime)
	}

	core.Post(fmt.Sprintf("http://%s:8080/v1/endpoint/stop", ipRemote), map[string]interface{}{
		"callid": callid,
	})
	releasePort(portLocal)
	<-c
}

func ReceiveLoop(ipLocal string, portStatis int) {
	addr, _ := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ipLocal, portStatis))
	conn, _ := net.ListenUDP("udp", addr)
	defer conn.Close()

	i := 0
	var p10, p25, p50, p100, p150, p200, p300, p400, p500, pbig int

	for true {
		var buf [200]byte
		_, _, err := conn.ReadFromUDP(buf[0:])
		if err != nil {
			return
		}
		//fmt.Println(remoteAddr, "recv:", string(buf[0:n]))
		timespan := time.Now().Sub(stringTotime(string(buf[0:19])))
		miliseconds := timespan.Seconds() * 1000
		if miliseconds < 10 {
			p10++
		} else if miliseconds < 25 {
			p25++
		} else if miliseconds < 50 {
			p50++
		} else if miliseconds < 100 {
			p100++
		} else if miliseconds < 150 {
			p150++
		} else if miliseconds < 200 {
			p200++
		} else if miliseconds < 300 {
			p300++
		} else if miliseconds < 400 {
			p400++
		} else if miliseconds < 500 {
			p500++
		} else {
			pbig++
		}

		i++
		//每一万个包统计一次
		if i == 1000 {
			i = 0
			if p10 == 1000 {
				//perfect
				fmt.Print("*")
			} else if p10+p25 == 1000 {
				//good
				fmt.Print("#")
			} else {
				fmt.Print("\n")
				fmt.Println(p10, p25, p50, p100, p150, p200, p300, p400, p500, pbig)
			}
			p10, p25, p50, p100, p150, p200, p300, p400, p500, pbig = 0, 0, 0, 0, 0, 0, 0, 0, 0, 0
		}
	}
}

func main() {
	clientip := flag.String("clientip", "10.153.138.126", "clientip")
	remoteip := flag.String("serverip", "10.153.138.128", "remoteip")
	threads := flag.Int("threads", 2000, "count of concurrent threads")
	portStatis := flag.Int("portStatis", 4000, "local port")


	flag.Parse()

	go ReceiveLoop(*clientip, *portStatis)

	ch := make(chan int, *threads)

	for index := 0;index<1;{
		index++
		ch <- 0
		clientport, _ := selectPort(*threads)
		//fmt.Println("clientport:", clientport)
		callid := fmt.Sprintf("%s-%d", *clientip, index)
		status1, resp1, _ := core.Post(fmt.Sprintf("http://%s:8080/v1/endpoint/bindleft", *remoteip), map[string]interface{}{
			"callid":   callid,
			"leftip":   clientip,
			"leftport": clientport,
		})
		if status1 != 200 {
			//fmt.Println(resp1.Get("msg").MustString())
		}

		status2, _, _ := core.Post(fmt.Sprintf("http://%s:8080/v1/endpoint/bindright", *remoteip), map[string]interface{}{
			"callid":    callid,
			"rightip":   clientip,
			"rightport": *portStatis,
		})
		if status2 != 200 {
			//fmt.Println(resp2.Get("msg").MustString())
		}

		go SendLoop(callid, *clientip, clientport, *remoteip, resp1.Get("port").MustInt(), ch)
	}

	fmt.Println("Client start...")
	c := make(chan os.Signal)
	signal.Notify(c)
	<-c

	//释放所有链接

	fmt.Println("Bye!")

}
