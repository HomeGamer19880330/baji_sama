package NetServer

//gorilla
//http://www.gorillatoolkit.org/
import (
	"fmt"
	"time"
	"net/http"
)

import (
	"component/NetCommunicator"
	"component/Utils/Errors"

	"github.com/gorilla/websocket"
	"component/UniqueIdGenerator"
)

const (
	readBufferSize = 64 << 10        // 读消息缓冲区64K
	writeBufferSize = 64 << 10       // 写消息缓冲区64K
}

// 服务端组件
type WebSocketServer struct {
	inMsgLimit      int                          // 接收消息管道容量上限
	outMsgLimit     int                          // 发送消息管道容量上限
	isStarted       bool
	addr            string // 服务地址
	allConnectsBetweenComputerLock       sync.RWMutex      
	allConnectsBetweenComputer         	map[uint32]*NetCommunicator.ConnBetweenTwoComputer                                // 保存连接的map
	//msgProc   framework.MsgProcessor // 消息处理方
	//isGateway bool                   // 该服务组件是否属于接入服
	serverId uint32 // 服务id
	UniqueIdGenerator *idGenerator
	upgrader *websocket.Upgrader //用来创建websocket tcp连接的东西
}

//var (
//	svrInstance *NetServer // 服务组件单件实例
//	//uuidGen     = new(uuid.UuidGen) // 客户端连接id生成器
//)

const (
	logTag = "WebSocketServer"
)

func (self *WebSocketServer) Initialize(addr string) bool {
	self.addr = addr
	self.upgrader = &websocket.Upgrader{ReadBufferSize: readBufferSize, WriteBufferSize: writeBufferSize}
	self.idGenerator = new(UniqueIdGenerator.UniqueIdGenerator)
	return true
}

// 服务端组件的字符串表示
func (self *NetServer) String() string {
	//	cfg, _ := cfgMgr.Instance().Get(s.cfgName).(*protoCfg.ServerCfg)
	//	if cfg != nil {
	//	return fmt.Sprintf("<server,%s>", *cfg.Id)
	//	}
	return fmt.Sprintf("<server>")
}

// 服务端组件是否属于网关服务
/*func (self *NetServer) IsGateWay() bool {
	return self.isGateway
}*/

// 获取服务地址
func (self *NetServer) GetAddr() string {
	//cfg := cfgMgr.Instance().Get(s.cfgName).(*protoCfg.ServerCfg)
	//return cfg.Addrs[0]
	return ""
}

// 服务组件每个连接的消息循环
func (self *WebSocketServer) connJob(connectOnTwoComputer *NetCommunicator.ConnBetweenTwoComputer) {
	//	defer communicator.Close()

	// 等待客户端组件的注册命令
	//cmd := new(proto.CmdMsg)
	//if !s.isGateway {
	//msg, recvErr := conn.RecvMsg(time.Second * 2)
	//if recvErr != nil {
	//s.logger.Errorf(logTag, "reg for con %s failed: %s", conn, recvErr)
	//return
	//}
	//if msg.MsgType != uint16(proto.MsgType_CMD) {
	//s.logger.Errorf(logTag, "invalid msg type %d on con %s", msg.MsgType, conn)
	//return
	//}
	//if e := protobuf.Unmarshal(msg.Data, cmd); e != nil {
	//s.logger.Errorf(logTag, "decode cmd on con %s failed!", conn)
	//return
	//}
	//if *cmd.Cmd != proto.CmdID_REG {
	//s.logger.Errorf(logTag, "invalid cmd id %d on con %s", *cmd.Cmd, conn)
	//return
	//}
	//if *cmd.DstServiceID != s.serverId {
	//s.logger.Errorf(logTag, "invalid server id %d, %d on con %s",
	//*cmd.DstServiceID, s.serverId, conn)
	//return
	//}
	//conn.SetId(*cmd.SrcServiceID)
	//s.logger.Debugf(logTag, "reg for conn %s succeeded!", conn)
	//} else {
	//conn.SetId(uuidGen.GenID())
	//s.logger.Debugf(logTag, "extablished new conn %s for client!", conn)
	//}

	// 通知新连接建立
	//newConnCh := s.connMgr.GetNewConnChan()
	//if newConnCh == nil {
	//s.logger.Errorf(logTag, "get new conn channel for con %s failed!", conn)
	//return
	//}
	//newConnCh <- conn

	// 注册连接异常断开回调函数
	//conn.RegisterDisconnCallback(func(arg interface{}) {
	//diconnCh := s.connMgr.GetDisconnChan()
	//if diconnCh == nil {
	//s.logger.Errorf(logTag, "get disconn channel for con %s failed!", conn)
	//return
	//}

	// 通知连接异常断开
	//diconnCh <- conn
	//}, nil)

	// 进入消息循环
	//	maxIdleTime := 15
	//	lastNotifyTime := time.Now().Unix()
	//	checkTimer := time.NewTicker(time.Second * 5)
	//	defer checkTimer.Stop()
	for {
		// value :=
		//
		if connectOnTwoComputer.IsClosed() {
			self.removeConnectsBetweenComputer(connectOnTwoComputer)
			return
		}

		//if !s.isGateway {
		//// 非接入服的服务端组件每5秒钟检查一次保活状态
		//select {
		//case <-checkTimer.C:
		//nowSeconds := time.Now().Unix()
		//if nowSeconds-lastNotifyTime > int64(maxIdleTime) {
		//s.logger.Errorf(logTag, "client idle too long on con %s", conn)
		//return
		//}
		//default:
		//}
		//}

		msg, recvErr := connectOnTwoComputer.RecvMsg(time.Second)
		if recvErr != nil {
			self.removeConnectsBetweenComputer(connectOnTwoComputer)
			continue
		}
		fmt.Printf("server recieve message", msg.Data)
		//if msg.MsgType == uint16(proto.MsgType_CMD) {
		//if s.isGateway {
		//// 接入服不可能收到信令
		//s.logger.Errorf(logTag, "invalid msg from client %s", conn)
		//return
		//}
		//if e := protobuf.Unmarshal(msg.Data, cmd); e != nil {
		//s.logger.Errorf(logTag, "decode cmd on con %s failed!", conn)
		//return
		//}
		//if *cmd.Cmd != proto.CmdID_KEEPALIVE {
		//s.logger.Errorf(logTag, "invalid cmd id %d on con %s", *cmd.Cmd, conn)
		//return
		//}
		//if *cmd.DstServiceID != s.serverId {
		//s.logger.Errorf(logTag, "invalid server id %d, %d on con %s",
		//*cmd.DstServiceID, s.serverId, conn)
		//return
		//}
		//lastNotifyTime = time.Now().Unix()
		//} else {
		//if s.isGateway {
		//// 解码消息并加添加nethead头信息
		//var m proto.Msg
		//if err := protobuf.Unmarshal(msg.Data, &m); err != nil {
		//s.logger.Errorf(logTag, "decode msg failed, %s", err)
		//return
		//}
		//header := m.GetHeader()
		//header.NHeader = s.GenNetHead(conn, nil)
		//data, err := protobuf.Marshal(&m)
		//if err != nil {
		//s.logger.Errorf(logTag, "marshal msg failed, %s", err)
		//return
		//}
		//msg.Data = data
		//}
		//
		//if e := s.msgProc.OnNewMsg(msg.Data); e != nil {
		//s.logger.Errorf(logTag, "process msg failed: %s", e)
		//continue
		//}
		//}
	}
}

// 启动服务端组件
func (self *WebSocketServer) Start() (err error) {
	if self.isStarted {
		//		s.logger.Errorf(logTag, "server already started!")
		err = Errors.New("server already started")
		return err
	}
	
	http.HandleFunc("/ws", self.HttpHandleFunc)
	
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err)
	}
	self.isStarted = true
}

// 停止服务端组件
func (self *WebSocketServer) Stop() {
	//	s.mu.Lock()
	//	defer s.mu.Unlock()

	//	s.logger.Infof(logTag, "server %s is stopped", s)
	//	self.netCommunicator.Stop()
	//	s.isStarted = false
}
func (self *WebSocketServer) addConnectsBetweenComputer(){
	self.allConnectsBetweenComputerLock.Lock()
	self.allConnectsBetweenComputer[connectBetweenComputer.uniqueId] = connectBetweenComputer
	self.allConnectsBetweenComputerLock.Unlock()
}

func (self *WebSocketServer) removeConnectsBetweenComputer(connectBetweenComputer *NetCommunicator.ConnBetweenTwoComputer){
	_, exist := self.allConnectsBetweenComputerLock[node.uuid]
	if exist {
		self.allConnectsBetweenComputerLock.Lock()
		delete(self.allConnectsBetweenComputer, connectBetweenComputer.uniqueId)
		self.allConnectsBetweenComputerLock.Unlock()
	}
}
func (self *WebSocketServer) HttpHandleFunc(w http.ResponseWriter, r *http.Request) {
	//服务器还未初始化时直接返回，不做处理
	if server.inited == false {
		return
	}

	conn, err := self.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	
	connectBetweenComputer := NetCommunicator.NewConnBetweenTwoComputer(
		inMsgLimit,
		outMsgLimit,
		conn
		self.idGenerator.GenUint32)

		
	go self.connJob(connectBetweenComputer)
//	server.newConnChan <- conn
}

