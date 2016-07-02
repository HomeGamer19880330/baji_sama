package NetCommunicator

import (
//	"encoding/binary"
	"fmt"
//	"io"
	"net"
//	"sync"
//	"dawn"
	"time"
)

// 连接建立时调用的回调函数原型
type OnNewConnFunc func(c *Connect)

//创建一个网络服务器或者客户端的描述文件
type ConnectInfo struct {
//	netConnectorGo      net.Conn          // tcp/udp连接的类
	isServer            bool              // 是否是服务器
	ipAddr				*net.TCPAddr	  // 网址
	connType    		string            // 网络连接类型(tcp/tcp4/tcp6/udp/udp4/udp6)
	isEnabled 			bool              // 是否可用中
	inMsgLimit  		int            // 接收消息管道容量上限
	outMsgLimit 		int            // 发送消息管道容量上限
}

//一个真实两台机器连接,以及他们之间收发消息的处理
type ConnBetweenTwoComputer struct {
	connectInfo         *Connect            // 网络通信类型
	isClosed            bool              // 标明该连接是否已关闭
	closeReason         CloseReason       // 连接关闭原因
	inPipe              safechan.AnyChan  // 接收消息通道
	outPipe             safechan.AnyChan  // 发送消息通道
	conn                net.Conn          // tcp/udp连接
	id                  uint32            // 连接id
//	disconnCallbackFunc func(interface{}) // 连接异常断开的通知回调函数
//	disconnCallbackArg  interface{}       // 调用模块的私有参数
//	lastRecvTimeStamp   int64             // 记录socket上次收包的时间戳
//	aboutToClose        bool              // 标记该连接为将管道中的消息发送后即关闭
	sendbuf             []byte            // 消息发送缓存
	recvbuf             []byte            // 消息收取缓存
}

// 创建新的连接
func newConn(connectInfo *ConnectInfo, conn net.Conn) *ConnBetweenTwoComputer {
	newconn := new(ConnBetweenTwoComputer)
	newconn.connectInfo = connectInfo
	newconn.isClosed = false
	newconn.inPipe = make(safechan.AnyChan, c.inMsgLimit)
	newconn.outPipe = make(safechan.AnyChan, c.outMsgLimit)
	newconn.conn = conn
	logger.Debugf(newconn.commu.connType, "established new con %s", newconn)
	c.wgRecvConns.Add(1)
	go newconn.recv()
	if c.isTCP || !c.isServer {
		// udp服务器暂时只能作为数据接收方，因此不启动发送协程
		c.wgSendConns.Add(1)
		go newconn.send()
	}
	return newconn
}

//var int Home = 1

//func Hello() int {
//	fmt.Println("hello dawn")
//	return 100
//}

func NewConnect(isServer bool, connType string, addr string) *Connect {
	connect := new(Connect)
	connect.isServer = isServer
	ipAddr, error := net.ResolveTCPAddr(connType, addr)
	if error != nil {
		fmt.Println("hello world", error)

	}
	connect.ipAddr = ipAddr
//	net.ResolveTCPAddr(connType, addr)
	connect.connType = connType
	connect.isEnabled = false
	
	return connect
//	switch c.connType {
//	case "tcp", "tcp4", "tcp6":
//		c.isTCP = true
//		for _, s := range addrs {
//			d, err := net.ResolveTCPAddr(c.connType, s)
//			if err != nil {
//				logger.Errorf(c.connType, "resolve tcp addr %s failed %s", s, err)
//				return nil
//			}
//			c.tcpAddrs = append(c.tcpAddrs, d)
//		}
//	case "udp", "udp4", "udp6":
//		c.isTCP = false
//		for _, s := range addrs {
//			d, err := net.ResolveUDPAddr(c.connType, s)
//			if err != nil {
//				logger.Errorf(c.connType, "resolve udp addr %s failed %s", s, err)
//				return nil
//			}
//			c.udpAddrs = append(c.udpAddrs, d)
//		}
//	default:
//		logger.Errorf(c.connType, "invalid connType")
//		return nil
//	}

//	return c

}

//func(c *Connect) Hello() {
//	fmt.Println("Connect hello world")
//}

// 启动网络通信，开始监听/连接端口，连接建立成功后通过回调函数回传给调用模块
func(c *Connect) Start(f OnNewConnFunc) {
	if f == nil {
//		logger.Errorf(c.connType, "invalid parameter")
		return
	}

	if c.isEnabled {
//		logger.Warnf(c.connType, "communication instance already enabled")
		return
	}

	c.isEnabled = true
//	if c.isTCP {
		if c.isServer {
			startTCPServer(c, f)
		} else {
			startTCPClient(c, f)
		}
//	} else {
//		if c.isServer {
//			startUDPServer(c, f)
//		} else {
//			startUDPClient(c, f)
//		}
//	}
}

// tcp端口监听routine
const max = 3
func listen(c *Connect, f OnNewConnFunc) {
//	defer c.wgListeners.Done()
	fmt.Println(c.connType, "start listening to", c.ipAddr)
//	logger.Debugf(, c.tcpAddrs[index]).

	l, err := net.ListenTCP(c.connType, c.ipAddr)
	if err != nil {
		fmt.Println("listen to tcp addr failed", c.ipAddr, err)
		return
	}

	var tempDelay time.Duration
	for 
	{
		if !c.isEnabled {
			l.Close() // 网络通信已被关闭，退出协程
			break
		}

		con, err := l.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				time.Sleep(tempDelay)
				continue
			}
//			logger.Errorf(c.connType, "accept tcp addr %v failed, %s", c.ipAddr, err)
			return
		}
		tempDelay = 0

		// 启动保活, 9分钟(60 + 8 * 60)内没有任何通信则会自动断开连接
		con.(*net.TCPConn).SetKeepAlive(true)
		con.(*net.TCPConn).SetKeepAlivePeriod(time.Minute)
		
		c.netConnector = con
		// 通知调用模块新连接的建立
		f(c)
	}
}

// 启动tcp服务协程
func startTCPServer(c *Connect, f OnNewConnFunc) {
//	service:=":9090"
//  tcpAddr, err := net.ResolveTCPAddr("tcp4", c.ipAddr)
//  l,err := net.ListenTCP("tcp",tcpAddr)
//  conn,err := l.Accept()
	
//	go Handler(conn) //此处使用go关键字新建线程处理连接，实现并发

//	for i, _ := range c.tcpAddrs {
//		c.wgListeners.Add(1)
	go listen(c, f)
//	}
}

// 启动tcp连接协程
func startTCPClient(c *Connect, f OnNewConnFunc) {
	go func() {
//		for _, addr := range c.tcpAddrs {
			con, err := net.DialTCP(c.connType, nil, c.ipAddr)
			if err != nil {
				fmt.Println("connect to tcp addr failed", c.ipAddr, err)
//				logger.Errorf(c.connType, "dial tcp addr %v failed %s", addr, err)
//				continue
				return
			}
			c.netConnector = con
			// 通知调用模块新连接的建立
			f(c)
//		}
	}()
}


