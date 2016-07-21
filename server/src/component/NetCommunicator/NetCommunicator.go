package NetCommunicator

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	//	"sync"
	//	"dawn"
	"time"
)

import (
	"component/safechan"
)

var (
	maxSilentTime uint64 = 60 * 3 // 连接静默的最长时间，3分钟,因为有心跳存在,所以如果不是僵尸死连接也不至于出现这种3分钟静默的情况
)

const (
	socketBufSize = 64 << 10 // 消息缓冲区64K
	headSize      = 8        // 消息头占用字节数 size(4) + magic_number(2) + msg_type(2)
	magicNumber   = 0xdaff   // 协议头魔数"dawn 拂晓"
)

// 连接建立时调用的回调函数原型
type OnNewConnFunc func(c *ConnBetweenTwoComputer)

type NetMsg struct {
	MsgType uint16 // 消息类型
	Data    []byte // 消息体数据
}

// 分配新的消息包
func NewNetMsg(msgtype uint16, data []byte) *NetMsg {
	return &NetMsg{MsgType: msgtype, Data: data}
}

//创建一个网络服务器或者客户端的描述文件
type ConnectInfo struct {
	//	netConnectorGo      net.Conn          // tcp/udp连接的类
	isServer    bool         // 是否是服务器
	ipAddr      *net.TCPAddr // 网址
	connType    string       // 网络连接类型(tcp/tcp4/tcp6/udp/udp4/udp6)
	isEnabled   bool         // 是否可用中
	inMsgLimit  int          // 接收消息管道容量上限
	outMsgLimit int          // 发送消息管道容量上限
}

//一个真实两台机器连接,以及他们之间收发消息的处理
type ConnBetweenTwoComputer struct {
	connectInfo *ConnectInfo // 网络通信类型
	isClosed    bool         // 标明该连接是否已关闭
	//	closeReason         CloseReason       // 连接关闭原因
	inPipe  safechan.AnyChan // 接收消息通道
	outPipe safechan.AnyChan // 发送消息通道
	conn    net.Conn         // tcp/udp连接
	id      uint32           // 连接id
	//	disconnCallbackFunc func(interface{}) // 连接异常断开的通知回调函数
	//	disconnCallbackArg  interface{}       // 调用模块的私有参数
	lastRecvTimeStamp int64  // 记录socket上次收包的时间戳
	aboutToClose      bool   // 标记该连接为将管道中的消息发送后即关闭
	sendbuf           []byte // 消息发送缓存
	recvbuf           []byte // 消息收取缓存
}

// 创建新的连接
func newConnBetweenTwoComputer(connectInfo *ConnectInfo, conn net.Conn) *ConnBetweenTwoComputer {
	newconn := new(ConnBetweenTwoComputer)
	newconn.connectInfo = connectInfo
	newconn.isClosed = false
	newconn.inPipe = make(safechan.AnyChan, connectInfo.inMsgLimit)
	newconn.outPipe = make(safechan.AnyChan, connectInfo.outMsgLimit)
	newconn.conn = conn
	//	logger.Debugf(newconn.commu.connType, "established new con %s", newconn)
	//	connectInfo.wgRecvConns.Add(1)
	go newconn.recv()
	//	connectInfo.isTCP ||
	if !connectInfo.isServer {
		// udp服务器暂时只能作为数据接收方，因此不启动发送协程
		//		c.wgSendConns.Add(1)
		go newconn.send()
	}
	return newconn
}

//var int Home = 1

func Hello() int {
	fmt.Println("hello dawn")
	return 100
}

func NewConnectInfo(isServer bool, connType string, addr string) *ConnectInfo {
	connect := new(ConnectInfo)
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

// 收包routine
func (self *ConnBetweenTwoComputer) recv() {
	//	defer c.commu.wgRecvConns.Done()
	self.recvbuf = make([]byte, socketBufSize)
	offset := uint32(0)
	readOffset := uint32(0)
	currentTimeStamp := time.Now().Unix() //用这个不是比较正常么
	self.lastRecvTimeStamp = currentTimeStamp

	for {
		if !self.connectInfo.isEnabled {
			// 网络通信已被关闭
			self.closeNow(true, false)
		}

		if self.isClosed || self.aboutToClose {
			return
		}

		// 若缓存已满则不再继续收包
		// 从tcp连接的数据缓存中读取数据到recvbuf,并计算出读取到的总数据长度为offset，直到读出所有数据
		// 可能读出了多个消息,每个消息前面有四字节用来描述消息长度
		// http://studygolang.com/articles/581
		if offset < socketBufSize {
			self.conn.SetReadDeadline(time.Now().Add(time.Second))
			recvLen, err := self.conn.Read(self.recvbuf[offset:])
			if err == io.EOF {
				// 对端主动关闭网络通信，通知调用模块
				//				logger.Debugf(c.commu.connType, "detected %v closed by the other end", c)
				self.closeNow(false, true)
				return
			} else if err != nil {
				// 检查是否存在系统错误
				//				ne, ok := err.(net.Error)
				//				if !ok || (!ne.Temporary() && !ne.Timeout()) {
				//					logger.Errorf(c.commu.connType,
				//						"read con %v failed %v, close invalid connection", c, err)
				//					c.closeNow(false, false)
				//					return
				//				}
			}
			//&& self.connectInfo.isTCP
			if self.connectInfo.isServer {
				if err == nil {
					self.lastRecvTimeStamp = currentTimeStamp
				} else if currentTimeStamp-self.lastRecvTimeStamp > int64(maxSilentTime) {
					// 关闭静默时间过长的连接
					//					logger.Errorf(c.commu.connType, "conn %v has been silent too long!", c)
					self.closeNow(false, false)
					return
				}
			}

			offset += uint32(recvLen)
			if offset < headSize {
				// 至少要收到8个字节才能校验消息头
				continue
			}

			// 判断是否接收到完整的消息报文，没有则继续收取
			size := binary.BigEndian.Uint32(self.recvbuf[:4])
			if size <= headSize || size >= socketBufSize {
				//				logger.Errorf(c.commu.connType,
				//					"invalid packet size %d, close invalid connection %v", size, c)
				self.closeNow(false, false)
				return
			}
			if offset < size {
				continue
			}
		}

		//处理消息缓冲区
		//把读取到的数据,转成消息结构放到读入管道中

		readOffset = 0
		for {
			size := binary.BigEndian.Uint32(self.recvbuf[readOffset : readOffset+4])
			if size >= socketBufSize {
				//				logger.Errorf(c.commu.connType,
				//					"packet size too large %d, close invalid connection, %v", size, c)
				self.closeNow(false, false)
				return
			}

			magic := binary.BigEndian.Uint16(self.recvbuf[readOffset+4 : readOffset+6])
			if magic != uint16(magicNumber) {
				//				logger.Errorf(c.commu.connType, "invalid magic number %x", magic)
				self.closeNow(false, false)
				return
			}

			msgtype := binary.BigEndian.Uint16(self.recvbuf[readOffset+6 : readOffset+headSize])

			// 向收取管道中添加新的消息体
			data := make([]byte, size-headSize)
			copy(data, self.recvbuf[readOffset+headSize:readOffset+size])
			self.inPipe.Write(NewNetMsg(msgtype, data), nil, time.Second*5)
			//			err :=
			//			if err != nil {
			//				if err == safechan.CHERR_TIMEOUT {
			//					logger.Debugf(c.commu.connType, "send msg into pipe timeout! %d, %d, %s",
			//						len(c.inPipe), cap(c.inPipe), c)
			//					continue
			//				} else if err == safechan.CHERR_CLOSED {
			//					logger.Debugf(c.commu.connType, "channel is closed! con: %s %s", c, c.closeReason)
			//					return
			//				} else {
			//					logger.Errorf(c.commu.connType, "write channel failed! con: %s %s", c, c.closeReason)
			//					c.closeNow(false, false)
			//					return
			//				}
			//			}

			// 分析缓冲区中是否还有消息可以继续处理
			readOffset += size
			if offset-readOffset > headSize {
				size := binary.BigEndian.Uint32(self.recvbuf[readOffset : readOffset+4])
				if size >= socketBufSize {
					//					logger.Errorf(c.commu.connType,
					//						"packet size too large %d, close invalid connection %v", size, c)
					//					self.closeNow(false, false)
					return
				}

				if size+readOffset <= offset {
					continue
				}
			}
			//			copy(c.recvbuf, c.recvbuf[readOffset:offset])
			offset = offset - readOffset
			break
		}
	}
}

// 从消息管道收取报文
func (self *ConnBetweenTwoComputer) readPacketsFromPipe() error {
	self.sendbuf = self.sendbuf[:0]
	readSize := 0
	sleepSeconds := 0

	// 从发送管道中获取报文
	for readSize < socketBufSize/2 {
		re, e := self.outPipe.Read(nil, time.Second*time.Duration(sleepSeconds))
		if e != nil {
			if e == safechan.CHERR_TIMEOUT {
				return nil
			} else if e == safechan.CHERR_EMPTY {
				if readSize > 0 {
					break //发送收到的消息
				} else {
					sleepSeconds = 5 //等待5s
					continue
				}
			} else if e == safechan.CHERR_CLOSED {
				//				logger.Debugf(c.commu.connType,
				//					"pipe is closed! con: %s %s", c, c.closeReason)
				return e
			} else {
				//				logger.Errorf(c.commu.connType,
				//					"read pipe failed! con: %s %s", c, c.closeReason)
				self.closeNow(false, false)
				return e
			}
		} else {
			msg := re.([]byte)
			self.sendbuf = append(self.sendbuf, msg...)
			readSize += len(msg)
			sleepSeconds = 0
		}
	}

	return nil
}

// 发包routine
func (self *ConnBetweenTwoComputer) send() {
	//defer self.connectInfo.wgSendConns.Done()
	self.sendbuf = make([]byte, 0, socketBufSize*2)

	for {
		if self.isClosed {
			return
		}

		if self.aboutToClose && len(self.outPipe) == 0 {
			self.closeNow(true, false)
		}

		// 收取报文
		e := self.readPacketsFromPipe()
		if e != nil {
			//			logger.Errorf(c.commu.connType, "read packets failed: %s", e)
			return
		}
		if len(self.sendbuf) == 0 {
			continue
		}

		// 向socket发送报文
		sendLen, err := self.conn.Write(self.sendbuf)
		if err != nil || sendLen != len(self.sendbuf) {
			if !self.isClosed {
				//				logger.Errorf(c.commu.connType,
				//					"write con %v failed: %s, close invalid connection", c, err)
				self.closeNow(false, false)
			}
			return
		}
	}
}

// 启动网络通信，开始监听/连接端口，连接建立成功后通过回调函数回传给调用模块
func (c *ConnectInfo) Start(f OnNewConnFunc) {
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

func listen(c *ConnectInfo, f OnNewConnFunc) {
	//	defer c.wgListeners.Done()
	fmt.Println(c.connType, "start listening to", c.ipAddr)
	//	logger.Debugf(, c.tcpAddrs[index]).

	l, err := net.ListenTCP(c.connType, c.ipAddr)
	if err != nil {
		fmt.Println("listen to tcp addr failed", c.ipAddr, err)
		return
	}

	var tempDelay time.Duration
	for {
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

		//		c.netConnector = con
		// 通知调用模块新连接的建立
		f(newConnBetweenTwoComputer(c, con))
	}
}

// 启动tcp服务协程
func startTCPServer(c *ConnectInfo, f OnNewConnFunc) {
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
func startTCPClient(c *ConnectInfo, f OnNewConnFunc) {
	go func() {
		//		for _, addr := range c.tcpAddrs {
		con, err := net.DialTCP(c.connType, nil, c.ipAddr)
		if err != nil {
			fmt.Println("connect to tcp addr failed", c.ipAddr, err)
			//				logger.Errorf(c.connType, "dial tcp addr %v failed %s", addr, err)
			//				continue
			return
		}
		//			c.netConnector = con
		// 通知调用模块新连接的建立
		f(newConnBetweenTwoComputer(c, con))
		//		}
	}()
}

// 关闭连接，内部函数
//	normalClose: 表示由服务主动关闭
//	isEOF: 表示通信对端已关闭
func (self *ConnBetweenTwoComputer) closeNow(normalClose bool, isEOF bool) {
	if normalClose && len(self.outPipe) > 0 && !self.aboutToClose {
		//主动关闭连接之前需将管道中的数据发送完毕
		self.aboutToClose = true
		return
	}

	err := self.conn.Close()
	if err == nil {
		//self.isClosed = true
		//if normalClose {
		//	self.closeReason = BY_SELF
		//} else {
		//	if isEOF {
		//		self.closeReason = BY_OTHER
		//	} else {
		//		self.closeReason = BY_ACCIDENT
		//	}
		//}
		close(self.inPipe)
		close(self.outPipe)
		//		logger.Debugf(c.commu.connType, "connection %s closed %s", c, c.closeReason)

		//		if c.disconnCallbackFunc != nil {
		//			c.disconnCallbackFunc(c.disconnCallbackArg)
		//		}
	}
}

// 发送消息包至消息管道，成功返回true，超时或失败则返回false;
// 参数d为0表示消息发送会阻塞直至成功返回
func (self *ConnBetweenTwoComputer) SendMsg(msg *NetMsg, d time.Duration) error {
	//	if self.aboutToClose {
	//		logger.Errorf(c.commu.connType, "connection about to close")
	//		return errors.New("connection about to close")
	//	}
	//
	//	if self.isClosed {
	//		logger.Errorf(c.commu.connType, "connection already closed")
	//		return errors.New("connection closed")
	//	}
	//
	//	if !self.commu.isTCP && self.commu.isServer {
	//		logger.Errorf(c.commu.connType, "udp server doesn't support sending message yet!")
	//		return errors.New("udp server can't send msg")
	//	}
	//
	//	if len(msg.Data) == 0 || len(msg.Data) >= socketBufSize-headSize {
	//		logger.Errorf(c.commu.connType, "invalid msg len %d", len(msg.Data))
	//		return errors.New("invalid msg format")
	//	}

	// 添加消息头
	//	http://blog.csdn.net/sunshine1314/article/details/2309655
	packet := make([]byte, len(msg.Data)+headSize)
	binary.BigEndian.PutUint32(packet[0:4], uint32(len(msg.Data)+headSize))
	binary.BigEndian.PutUint16(packet[4:6], uint16(magicNumber))
	binary.BigEndian.PutUint16(packet[6:headSize], msg.MsgType)
	copy(packet[headSize:], msg.Data)

	err := self.outPipe.Write(packet, nil, d)
	if err != nil {
		//		logger.Errorf(c.commu.connType, "write failed: %s", err)
		return err
	}
	return nil
}

// 发送信令，主要用于客户端组件和服务端组件之间的控制协议
//func (self *ConnBetweenTwoComputer) SendCmd(data []byte, d time.Duration) error {
//	return self.SendMsg(NewNetMsg(uint16(proto.MsgType_CMD), data), d)
//}

// 发送正常通信协议
func (self *ConnBetweenTwoComputer) SendNormalMsg(data []byte, d time.Duration) error {
	//	if logger.GetLevel() == log.DEBUG {
	//		var m proto.Msg
	//		e := protobuf.Unmarshal(data, &m)
	//		if e == nil && m.GetHeader().GetSHeader().GetSeq() == 0 {
	//			m.GetHeader().GetSHeader().Seq = protobuf.Uint32(uuidGen.GenID())
	//			m.GetHeader().GetSHeader().TimeStamp =
	//			protobuf.Uint64(uint64(time.Now().UnixNano()))
	//			data, _ = protobuf.Marshal(&m)
	//		}
	//	}
	return self.SendMsg(NewNetMsg(uint16(0), data), d)
}

// 从消息管道中获取消息包，成功返回消息切片和nil，超时或失败则返回nil和错误;
// 参数d为0表示在没有消息可收的情况下马上返回nil和空管道错误
func (self *ConnBetweenTwoComputer) RecvMsg(d time.Duration) (*NetMsg, error) {
	msg, err := self.inPipe.Read(nil, d)
	//	if err != nil {
	//		if err != safechan.CHERR_TIMEOUT && err != safechan.CHERR_EMPTY {
	//		logger.Errorf(c.commu.connType, "read failed: %s", err)
	//		}
	//		return nil, err
	//	}

	re, ok := msg.(*NetMsg)
	if !ok {
		//		logger.Errorf(c.commu.connType, "invalid msg type %v", msg)
		return nil, err
	}
	return re, nil
}

// func (self *ConnBetweenTwoComputer) OhMyGod() bool {
// 	return true
// }

func (self *ConnBetweenTwoComputer) IsClosed() bool {
	return self.isClosed
}

// 关闭连接
func (self *ConnBetweenTwoComputer) Close() {
	self.closeNow(true, false)
}

// // 关闭连接
// func (self *ConnBetweenTwoComputer) SayHello() {

// }

// // 检查连接的状态
// func (self *ConnBetweenTwoComputer) IsClosed() bool {
// 	return self.isClosed
// }

// func (self *ConnBetweenTwoComputer) IsClosedQ() bool {
// 	return true
// 	//return self.isClosed
// }
