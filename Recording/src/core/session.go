package core

import (
	"config"
	"context"
	"errors"
	"fmt"
	"net"
	"record"
	"sync"
	"time"
)

type Session struct {
	sessionid string
	left      *net.UDPAddr
	right     *net.UDPAddr
	portSuite *PortSuite
	media     string //711 amrnb amrwb

	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
	wait200 chan bool
}

func (s Session) String() string {
	str := fmt.Sprintf("left=%s:%d", s.left.IP.String(), s.left.Port)
	if s.portSuite != nil {
		str += fmt.Sprintf("\t rtp=%d rtcp=%d", s.portSuite.PortRtp, s.portSuite.PortRtcp)
	}
	if s.right != nil {
		str += fmt.Sprintf("\t right=%s:%d", s.right.IP.String(), s.right.Port)
	} else {
		str += "\t right=nil"
	}
	return str
}

func (s *Session) Start() {

	sess, _ := GetSessionManage()
	s.wg.Add(2)

	//process rtp
	go func() {
		localAddr, _ := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", sess.localIp, s.portSuite.PortRtp))
		conn, _ := net.ListenUDP("udp", localAddr)
		defer conn.Close()

		remoteAddrLeft := s.left
		remoteAddrRight := s.right

		fileLeft := record.CreateRecordFile(fmt.Sprintf("/tmp/record/%s_early_left.amr", s.sessionid), s.media)
		fileRight := record.CreateRecordFile(fmt.Sprintf("/tmp/record/%s_early_right.amr", s.sessionid), s.media)
		defer fileLeft.Close()
		defer fileRight.Close()

		t1, t2 := config.AppConf.Timeout.T1, config.AppConf.Timeout.T2
		tWait := t1
		data := make([]byte, 200)
		for {
			select {
			case <-s.wait200:
				fileLeft = record.CreateRecordFile(fmt.Sprintf("/tmp/record/%s_left.amr", s.sessionid), s.media)
				fileRight = record.CreateRecordFile(fmt.Sprintf("/tmp/record/%s_right.amr", s.sessionid), s.media)
				defer fileLeft.Close()
				defer fileRight.Close()
			case <-s.ctx.Done():
				// fmt.Println("<<rtp stoped!>> - cancel")
				s.wg.Done()
				return
			default:
				// fmt.Println("------------- rtp ------------>")
				conn.SetReadDeadline(time.Now().Add(time.Millisecond * time.Duration(t1)))
				n, remoteAddr, err := conn.ReadFromUDP(data)
				if err != nil {
					//fmt.Println("rtp:", err, " twait=", tWait)
					tWait += t1
					if tWait > t2 {
						//exit
						sess, _ := GetSessionManage()
						sess.Stop(s.sessionid)
						//fmt.Println("<<rtp stoped!>> - timeout")
						s.wg.Done()
						return
					}
					continue
				} else {
					tWait = t1
				}

				//fmt.Println(remoteAddr, "receive:", string(data[:n]))

				if remoteAddr.Port == s.left.Port {
					conn.WriteToUDP(data[:n], remoteAddrRight)
					fileLeft.Write(data[:n])
				} else {
					conn.WriteToUDP(data[:n], remoteAddrLeft)
					fileRight.Write(data[:n])
				}
			}
		}
	}()

	//process rtcp
	go func() {
		localAddr, _ := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", sess.localIp, s.portSuite.PortRtcp))
		conn, _ := net.ListenUDP("udp", localAddr)
		defer conn.Close()
		remoteAddrLeft, _ := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", s.left.IP.String(), s.left.Port+1))
		remoteAddrRight, _ := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", s.right.IP.String(), s.right.Port+1))
		fmt.Printf("start rtcp - left:%s  right:%s\n", remoteAddrLeft, remoteAddrRight)
		data := make([]byte, 200)
		for {
			select {
			case <-s.ctx.Done():
				fmt.Println("<<rtcp stoped!>> - cancel")
				s.wg.Done()
				return
			default:
				conn.SetReadDeadline(time.Now().Add(time.Second * 1))
				n, remoteAddr, err := conn.ReadFromUDP(data)
				if err != nil {
					//fmt.Println("rtcp:", err)
					continue
				}
				//fmt.Println(remoteAddr, "receive:", string(data[:n]))

				if remoteAddr.Port == (s.left.Port + 1) {
					conn.WriteToUDP(data[:n], remoteAddrRight)
					fmt.Println("from left:", s.left.IP.String(), "to right:", remoteAddrRight.IP.String())
				} else if remoteAddr.Port == (s.right.Port + 1) {
					conn.WriteToUDP(data[:n], remoteAddrLeft)
					fmt.Println("from right:", s.right.IP.String(), "to left:", remoteAddrLeft.IP.String())
				} else {
					fmt.Println("fuckkkkkkkkkkkkkk")
				}
			}
		}
	}()
}

func (s *Session) Update(leftip string, leftport int, rightip string, rightport int) {
	s.left, _ = net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", leftip, leftport))
	s.right, _ = net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", rightip, rightport))
}

func (s *Session) SplitRecord() {
	s.wait200 <- true
}

type SessionManage struct {
	store   map[string]*Session // call count
	localIp string
	mu      sync.Mutex
}

func (sess *SessionManage) GetLocalIp() string {
	return sess.localIp
}

var sessionManage *SessionManage

func InitSessionManage(localIp string) (*SessionManage, error) {
	if sessionManage == nil {
		sessionManage = &SessionManage{store: map[string]*Session{}, localIp: localIp}
		return sessionManage, nil
	} else {
		return nil, errors.New("session manage has been init before")
	}
}

func GetSessionManage() (*SessionManage, error) {
	if sessionManage == nil {
		return nil, errors.New("session has not been init")
	} else {
		return sessionManage, nil
	}
}

func (sess *SessionManage) Size() int {
	sess.mu.Lock()
	defer sess.mu.Unlock()

	return len(sess.store)
}
func (sess *SessionManage) BindLeft(callid string, leftip string, leftport int) (*PortSuite, error) {
	sess.mu.Lock()
	defer sess.mu.Unlock()

	if _, ok := sess.store[callid]; ok {
		return nil, errors.New("callid already exists")
	}

	ps, err := GetPortSuiteHelper().AllotPort()
	if err != nil {
		return nil, err
	}

	leftAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", leftip, leftport))
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	sess.store[callid] = &Session{
		sessionid: callid,
		left:      leftAddr,
		right:     nil,
		portSuite: ps,
		ctx:       ctx,
		cancel:    cancel,
		wg:        sync.WaitGroup{},
		wait200:   make(chan bool),
	}

	return ps, nil
}

func (sess *SessionManage) BindRight(callid string, rightip string, rightport int, media string) error {
	sess.mu.Lock()
	defer sess.mu.Unlock()

	s, ok := sess.store[callid]
	if !ok {
		return errors.New("callid does not exist")
	}
	rightAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", rightip, rightport))
	if err != nil {
		return err
	}

	s.right = rightAddr
	s.media = media
	return nil
}

func (sess *SessionManage) Start(callid string) error {
	sess.mu.Lock()
	defer sess.mu.Unlock()

	s, ok := sess.store[callid]
	if !ok {
		return errors.New(fmt.Sprintf("Session [callid=%s]does not exist", callid))
	}

	//fmt.Println("start:", s)
	s.Start()
	return nil
}

func (sess *SessionManage) SplitRecord(sessionid string) error {
	sess.mu.Lock()
	defer sess.mu.Unlock()

	s, ok := sess.store[sessionid]
	if !ok {
		return errors.New(fmt.Sprintf("Session [sessionid=%s]does not exist", sessionid))
	}

	s.SplitRecord()
	return nil
}

/* 此函数最长可能会阻塞1S时间 */
func (sess *SessionManage) Stop(callid string) error {
	sess.mu.Lock()
	defer sess.mu.Unlock()

	s, ok := sess.store[callid]
	if !ok {
		return errors.New(fmt.Sprintf("Session [callid=%s]does not exist", callid))
	}
	go func() {
		//关闭所有线程
		s.cancel()
		//阻塞等待RTP和RTCP线程关闭
		s.wg.Wait()

		//清理会话对象
		sess.mu.Lock()
		delete(sess.store, callid)
		defer sess.mu.Unlock()
		GetPortSuiteHelper().ReleasePort(s.portSuite.PortRtcp)
	}()
	return nil
}

func (sess *SessionManage) Preview() {
	fmt.Println("--------------------")
	fmt.Println(sess.store)
	fmt.Println("--------------------")
}
