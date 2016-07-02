package NetServer

// server组件实现的功能如下:
//	1. 通过标准配置格式加载、重加载服务器id、监听端口和ip地址；
//	2. 启动tcp服务并接收来自客户端的连接；
//package server

import (
	"fmt"
	"sync"
	"time"
)

import (
//protobuf "github.com/golang/protobuf/proto"
//cfgMgr "ximigame.com/component/cfg"
//"ximigame.com/component/log"
	"server/component/NetCommunicator"
//"ximigame.com/component/uuid"
//"ximigame.com/framework"
//"ximigame.com/types/proto"
//protoCfg "ximigame.com/types/proto/config"
//"ximigame.com/utils"
//"ximigame.com/utils/errors"
)

// 服务端组件
type NetServer struct {
	netCommunicator      *NetCommunicator.NetCommunicator             // 网络组件
	//logger    *log.Logger            // 日志组件
	//connMgr   ConnManager            // 连接管理器
	//mu        sync.Mutex             // 保护锁
	//cfgName   string                 // 全局配置名
	//isStarted bool                   // 是否已启动
	//msgProc   framework.MsgProcessor // 消息处理方
	//isGateway bool                   // 该服务组件是否属于接入服
	serverId  uint32                 // 服务id
}

var (
	svrInstance *NetServer             // 服务组件单件实例
	uuidGen     = new(uuid.UuidGen) // 客户端连接id生成器
)

const (
	logTag = "svr_com"
)

// 获取服务端组件单件实例
func Instance() *NetServer {
	if svrInstance == nil {
		svrInstance = new(NetServer)
	}
	return svrInstance
}

// 初始化服务端组件
//	connMgr: 连接管理器
//	cfgName: 配置标识符
//	logger: 日志管理器
func (s *NetServer) Initialize() bool
//connMgr ConnManager, cfgName string,
//logger *log.Logger, proc framework.MsgProcessor)
//bool
{
//	net.SetLogDir(logger.GetDir())
//	s.logger = logger
//	s.connMgr = connMgr
//	s.cfgName = cfgName
//	s.msgProc = proc
	return true
}

// 服务端组件的字符串表示
func (self *NetServer) String() string
{
//	cfg, _ := cfgMgr.Instance().Get(s.cfgName).(*protoCfg.ServerCfg)
//	if cfg != nil {
//	return fmt.Sprintf("<server,%s>", *cfg.Id)
//	}
	return fmt.Sprintf("<server>")
}

// 服务端组件是否属于网关服务
func (self *NetServer) IsGateWay() bool {
	return self.isGateway
}

// 获取服务地址
func (s *Server) GetAddr() string {
//cfg := cfgMgr.Instance().Get(s.cfgName).(*protoCfg.ServerCfg)
//return cfg.Addrs[0]
	return ""
}

// 服务组件每个连接的消息循环
func (self *NetServer) connJob(communicator *NetCommunicator.NetCommunicator) {
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
		if communicator.IsClosed() {
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

		msg, recvErr := communicator.RecvMsg(time.Second)
		if recvErr != nil {
			continue
		}

if msg.MsgType == uint16(proto.MsgType_CMD) {
if s.isGateway {
// 接入服不可能收到信令
s.logger.Errorf(logTag, "invalid msg from client %s", conn)
return
}
if e := protobuf.Unmarshal(msg.Data, cmd); e != nil {
s.logger.Errorf(logTag, "decode cmd on con %s failed!", conn)
return
}
if *cmd.Cmd != proto.CmdID_KEEPALIVE {
s.logger.Errorf(logTag, "invalid cmd id %d on con %s", *cmd.Cmd, conn)
return
}
if *cmd.DstServiceID != s.serverId {
s.logger.Errorf(logTag, "invalid server id %d, %d on con %s",
*cmd.DstServiceID, s.serverId, conn)
return
}
lastNotifyTime = time.Now().Unix()
} else {
if s.isGateway {
// 解码消息并加添加nethead头信息
var m proto.Msg
if err := protobuf.Unmarshal(msg.Data, &m); err != nil {
s.logger.Errorf(logTag, "decode msg failed, %s", err)
return
}
header := m.GetHeader()
header.NHeader = s.GenNetHead(conn, nil)
data, err := protobuf.Marshal(&m)
if err != nil {
s.logger.Errorf(logTag, "marshal msg failed, %s", err)
return
}
msg.Data = data
}

if e := s.msgProc.OnNewMsg(msg.Data); e != nil {
s.logger.Errorf(logTag, "process msg failed: %s", e)
continue
}
}
}
}

// 启动服务端组件
func (self *NetServer) Start() (err error) {
	//s.mu.Lock()
	//defer s.mu.Unlock()
	if self.isStarted {
//		s.logger.Errorf(logTag, "server already started!")
		err = errors.New("server already started")
		return err
	}

//	// 读取配置
//	cfg, ok := cfgMgr.Instance().Get(s.cfgName).(*protoCfg.ServerCfg)
//	if !ok {
//	s.logger.Errorf(logTag, "get cfg failed")
//	err = errors.New("get cfg failed")
//	return
//	}
//	serverId, convertErr := utils.ConvertServerIDString2Number(*cfg.Id)
//	if convertErr != nil {
//	s.logger.Errorf(logTag, "invalid server id %s, %s", *cfg.Id, convertErr)
//	err = errors.New("invalid cfg")
//	return
//	}
	self.serverId = 10065
//	if cfg.IsGateway == nil {
//		s.logger.Errorf(logTag, "invalid cfg for server")
//		err = errors.New("invalid cfg")
//		return
//	}
//	s.isGateway = *cfg.IsGateway

	self.netCommunicator = NetCommunicator.NewConnect(true, "tcp", "127.0.0.1:15000")
//		int(*cfg.RecvLimit), int(*cfg.SendLimit))
	if self.netCommunicator == nil {
//		s.logger.Errorf(logTag, "allocate commu failed")
		err = errors.New("allocate commu failed")
		return
	}

	// 启动服务组件
//	s.logger.Infof(logTag, "server %s is started", s)
	self.netCommunicator.Start(func(conn *net.Conn) {
		go self.connJob(conn)
	})

	self.isStarted = true
	return nil
}

// 停止服务端组件
func (s *Server) Stop() {
//	s.mu.Lock()
//	defer s.mu.Unlock()

//	s.logger.Infof(logTag, "server %s is stopped", s)
//	self.netCommunicator.Stop()
//	s.isStarted = false
}

// 重启服务端组件
//	ch: 服务端消息接收管道
//func (s *Server) Restart() error {
//s.logger.Debugf(logTag, "server %s is restarted", s)
//s.Stop()
//return s.Start()
//}
//
//// 获取服务id
//func (s *Server) GetServerId() uint32 {
//return s.serverId
//}
//
//// 生成net header
//func (s *Server) GenNetHead(conn *net.Conn, cmd *proto.NetCmdID) *proto.NetHeader {
//ip, port, _ := conn.RemoteAddr()
//ip = ip.To4()
//var ipaddr uint32
//ipaddr = uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
//header := new(proto.NetHeader)
//header.ServiceID = protobuf.Uint32(s.serverId)
//header.SessionID = protobuf.Uint32(conn.GetId())
//header.ClientIP = protobuf.Uint32(ipaddr)
//header.ClientPort = protobuf.Uint32(uint32(port))
//header.Cmd = cmd
//return header
//}
