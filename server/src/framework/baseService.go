package framework

//基础服务实现
//package base

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"time"
)

import (
//	"github.com/golang/protobuf/proto"
//	"ximigame.com/component/cfg"
//	"ximigame.com/component/log"
//	"ximigame.com/component/net"
//	"ximigame.com/component/net/cfgsync"
	"component/NetCommunicator/NetClient"
	"component/NetCommunicator/NetServer"
//	"ximigame.com/component/oplog"
//	"ximigame.com/component/process"
//	"ximigame.com/component/reportdata"
//	"ximigame.com/component/routinepool"
//	"ximigame.com/component/timer"
//	"framework"
//	msg "ximigame.com/types/proto"
//	"ximigame.com/types/proto/config"
//	"ximigame.com/utils"
//	"ximigame.com/utils/errors"
)

//服务接口
type ServiceInterface interface {
	//	MsgProcessor
	Init(*FrameWork) (int, error)                                //初始化
	//	RegisterCfg() (int, error)                                   //注册配置
	//	SetLogLevel()                                                //设置日志等级
	SetupNetwork() (int, error)                                  //启动网络
	//	ProcessHttpCmd(h *process.HttpContext)                       //处理http命令
	//	ProcessTimer(tn *timer.TimeoutNotify)                        //处理定时器超时
	//	ProcessMsg(buff []byte) error                                //处理消息
	//	OnReload()                                                   //重载
	//	OnExit()                                                     //退出
	//	OnNetDisconn(conn *net.Conn)                                 //网络连接异常断开
	//	MainLoop()                                                   //主循环
	//	RegisterMsgHandle()                                          //注册所有消息处理
	//	RegOneMsgHandle(msgId uint32, handle MsgHandle) (int, error) //注册一个消息处理
}

//基础服务
type BaseService struct {
//	Proc          *process.Process          //进程管理组件
//	Log           *log.Logger               //运行日志组件
//	Cfg           *cfg.CfgMng               //本地配置管理组件
//	CfgSync       *cfgsync.CfgSync          //远程配置同步组件
//	UseLocalCfg   bool                      //使用本地配置
//	TM            *timer.TimerManager       //定时器管理组件
	serverInstance           *NetServer.NetServer            //网络服务组件
//	CM            *ConnMng                  //服务连接管理器
//	ReloadChan    chan int                  //用于通讯的各种管道, 接收重载信号的管道
	ExitChan      chan int                  //接收退出信号的管道
//	HttpChan      chan *process.HttpContext //接收http命令的管道
//	TNChan        chan *timer.TimeoutNotify //接收定时器超时通知的管道
//	MP            MsgMap                    //消息映射
	mainFrameWork *framework.FrameWork      //记录框架实例
//	ServiceID     uint32                    //保存服务id
//	cliCfgMap     map[string]uint32         //记录客户端组件的配置名和服务类型的映射关系
	clientMap     map[uint32]*client.Client //所有客户端组件，key为服务类型
//	OL            *oplog.OpLogger           //运营日志组件
//	TaskPool      *routinepool.RoutinePool  //消息处理协程池
//	StatusWrapper *reportdata.StatusWrapper //上报数据到状态服
//	Perf          *PerfMon                  //性能监控
}

const (
	LogTag              = "Service" //日志tag
	DefaultTNChanSize   = 10000     //默认定时器管道大小
	DefaultTaskPoolSize = 1000      //默认协程池大小
)

//实现Service接口:初始化函数
func (self *BaseService) Init(mainFrameWrok_ *framework.FrameWork) (int, error) {
		//创建各个管道
//		s.ReloadChan = make(chan int, 1)
		s.ExitChan = make(chan int, 1)
//		s.HttpChan = make(chan *process.HttpContext)
//		s.TNChan = make(chan *timer.TimeoutNotify, DefaultTNChanSize)
//
//		s.MP = NewMsgMap()
	self.mainFrameWork = mainFrameWrok_
//
//		//初始化客户端组件表
//		s.cliCfgMap = make(map[string]uint32)
//		s.CP = make(map[uint32]*client.Client)
//
//		//初始化进程管理组件
//		s.Proc = process.Instance()
//		_, e := s.Proc.Initialize("", "", "", "", s)
//		if e != nil {
//		panic(e.Error())
//		}
//
//		//初始化日志组件
//		s.Log = log.NewFileLogger(s.Proc.Name, s.Proc.DefaultLogDir, "")
//		if s.Log == nil {
//		panic("init logger failed")
//		}
//
//		defer log.FlushAll()
//
//		//设置日志等级
//		s.FW.Service.SetLogLevel()
//
//		//记录服务启动事件
//		s.Log.Infof(LogTag, "%s starting...", s.Proc.Name)
//
//		//初始化配置管理组件
//		s.Cfg = cfg.Instance()
//		s.Cfg.Initialize("", s.Log)
//
//		//远程配置同步组件初始化
//		s.CfgSync = cfgsync.Instance()
//		s.CfgSync.Initialize(s.Proc, s.Log, s.UseLocalCfg)
//		s.CfgSync.Start(s)
//
//		//注册配置
//		_, e = s.FW.Service.RegisterCfg()
//		if e != nil {
//		s.Log.Errorf(LogTag, "%s register cfg failed: %s", s.Proc.Name, e.Error())
//		return -1, e
//		}
//
//		//初始装载所有配置
//		_, e = s.Cfg.InitLoadAll()
//		if e != nil {
//		s.Log.Errorf(LogTag, "%s load cfg failed: %s", s.Proc.Name, e.Error())
//		return -1, e
//		}
//
//		//初始化定时器管理组件
//		s.TM = timer.Instance()
//		s.TM.Initialize(s.TNChan, s.Log)
//
//		//打印服务基本信息
//		fmt.Print(s.Proc)
//
//		//打印所有配置
//		s.Cfg.PrintAll()
//
//		//服务后台化
//		_, e = s.Proc.Daemonize()
//		if e != nil {
//		s.Log.Errorf(LogTag, "%s daemonize failed: %s", s.Proc.Name, e.Error())
//		return -1, e
//		}
//
//		//注册消息处理
//		s.FW.Service.RegisterMsgHandle()
//
//		//启动任务处理协程池
//		s.TaskPool = routinepool.NewPool(DefaultTaskPoolSize).Start()
//
//		//初始化性能监控
//		s.Perf = NewPerfMon(s)

	//启动网络组件
	errorCode, e = self.SetupNetwork()
	if e != nil {
//		s.Log.Errorf(LogTag, "%s set up network failed: %s", s.Proc.Name, e.Error())
		return -1, e
	}

	//记录服务启动完成日志
//	s.Log.Infof(LogTag, "%s started", s.Proc.Name)

	return 0, nil
}
//
////实现Service接口:注册所有配置
//func (s *BaseService) RegisterCfg() (int, error) {
//return 0, nil
//}
//
////实现远端配置加载回调接口
//func (s *BaseService) OnNewCfg(cfgName string, cfg proto.Message) (int, error) {
//switch cfg.(type) {
//case *config.ServerCfg:
//svrCfg := cfg.(*config.ServerCfg)
//s.Proc.HttpUrl = *svrCfg.HttpAddr
//serverID, err := utils.ConvertServerIDString2Number(*svrCfg.Id)
//if err != nil {
//panic(err)
//}
//s.ServiceID = serverID
//case *config.ClientCfg:
////判断配置是否合法，主要判断服务id是否合法且类型是否重复
//cliCfg := cfg.(*config.ClientCfg)
//if len(cliCfg.ServerEntries) == 0 {
//s.Log.Errorf(LogTag, "empty server entries!")
//return -1, errors.New("empty server entries!")
//}
//fields, err := utils.ConvertServerIDString2Fields(*cliCfg.ServerEntries[0].Id)
//if err != nil {
//s.Log.Errorf(LogTag, "invalid server id, %s", *cliCfg.ServerEntries[0].Id)
//return -1, errors.New("invalid server id, %s", *cliCfg.ServerEntries[0].Id)
//}
//svrType := fields[0]
//idMap := map[string]bool{}
//idMap[*cliCfg.ServerEntries[0].Id] = true
//for i, entry := range cliCfg.ServerEntries {
//if i == 0 {
//continue
//}
//fields, err = utils.ConvertServerIDString2Fields(*entry.Id)
//if err != nil {
//s.Log.Errorf(LogTag, "invalid server id, %s", *entry.Id)
//return -1, errors.New("invalid server id, %s", *entry.Id)
//}
//if fields[0] != svrType {
//s.Log.Errorf(LogTag, "inconsistent server type %s!", *entry.Id)
//return -1, errors.New("inconsistent server type %s!", *entry.Id)
//}
//if _, exists := idMap[*entry.Id]; exists {
//s.Log.Errorf(LogTag, "duplicate server id %s!", *entry.Id)
//return -1, errors.New("duplicate server id %s!", *entry.Id)
//}
//}
//s.cliCfgMap[cfgName] = svrType
//if com, exists := s.CP[svrType]; exists {
//com.ReloadCfg() // 通知各个客户端组件重加载配置
//}
//default:
//s.Log.Errorf(LogTag, "unknown cfg type, %s", reflect.TypeOf(cfg))
//return -1, errors.New("unknown cfg type, %s", reflect.TypeOf(cfg))
//}
//
//return 0, nil
//}
//
////实现Service接口:设置日志等级
//func (s *BaseService) SetLogLevel() {
//s.Log.SetLevel(log.ERROR)
//net.SetDebugLevel(log.ERROR)
//}
//
//实现ServiceInterface接口:启动网络
func (self *BaseService) SetupNetwork() (int, error) {
	//创建并初始化连接管理器
//	s.CM = new(ConnMng)
//	if e := s.CM.Init(s.FW.Service, s.Log); e != nil {
//		return -1, e
//	}

	//启动服务端组件
	self.serverInstance = NetServer.Instance()
	if self.serverInstance == nil {
		return -1, errors.New("get server instance failed")
	}
	if !self.serverInstance.Initialize() {
//		s.CM,
//	strings.TrimSuffix(cfgsync.SvrDeployFile, ".cfg"), s.Log, s)
		return -1, errors.New("server init failed")
	}
	e := self.serverInstance.Start()
	if e != nil {
		return -1, e
	}

	//启动客户端组件
	//原著中所有客户端SERVICE都有一个CFG同步组件CfgSync,启动后连接配置中心,从配置中心下载最新的连接配置放在deploy路径下,
	//deploy路径下的Addr: "127.0.0.1:15000" SvrType: 3 记录配置中心的IP地址和端口号,以及此服务的配置类型
//	serverId := self.Svr.GetServerId()
//	for cfgName, svrType := range s.cliCfgMap {
	self.clientMap[1] = client.NewClient(cfgName, s.Log, serverId, s)
	self.clientMap[1].Start()
//	}
//
////初始化运营日志组件
//logCli, exists := s.CP[uint32(msg.SvrType_LOG)]
//if exists {
//s.OL = oplog.Instance()
//s.OL.Initialize(logCli, s.Log)
//}
////状态服务组件
//{
//statusKeeperCli, exists := s.CP[uint32(msg.SvrType_StatusKeeper)]
//if exists {
//s.StatusWrapper = reportdata.Instance()
//s.StatusWrapper.Initialize(statusKeeperCli, s.Log)
//}
//}

return 0, nil
}
//
////实现Service接口:新消息到来
//func (s *BaseService) OnNewMsg(buff []byte) error {
//return s.FW.Service.ProcessMsg(buff)
//}
//
////实现Service接口:处理消息
//func (s *BaseService) ProcessMsg(buff []byte) error {
//pdata := new(PerfData)
//start := time.Now()
//err := s.TaskPool.DispatchTask(func() {
//waitChan := make(chan int, 1)
//pdata.taskWaitTime = time.Now().Sub(start)
//start = time.Now()
//defer func() {
//pdata.msgProcTime = time.Now().Sub(start)
//s.Perf.perfChan <- pdata
//waitChan <- 0
//}()
//
////解码消息包
//var m msg.Msg
//if e := proto.Unmarshal(buff, &m); e != nil {
//s.Log.Errorf("PROTO", "decode buffer length %d failed", len(buff))
//return
//}
//
//go func() {
//select {
//case <-waitChan:
//case <-time.After(time.Minute):
//buf := make([]byte, 1024)
//size := runtime.Stack(buf, false)
//s.Log.Errorf("PROTO",
//"Task busy too long!\nProcessing msg: %s\n%s\n",
//&m, string(buf[:size]))
//}
//}()
//
////处理消息
//_, e := s.MP.Process(&m, s)
//if e != nil {
//var msgId uint32
//if m.GetHeader() != nil || m.GetHeader().GetSHeader() != nil {
//msgId = m.GetHeader().GetSHeader().GetMsgID()
//}
//s.Log.Errorf("PROTO", "process msg %d failed:%s", msgId, e.Error())
//}
//})
//
//return err
//}
//
////实现Service接口:处理http命令
//func (s *BaseService) ProcessHttpCmd(h *process.HttpContext) {
//if h == nil {
//return
//}
//
//s.Log.Debugf("HTTP", "got http cmd %s", h.Req.URL)
//
//switch h.Req.FormValue("loglevel") {
//case "debug":
//s.Log.SetLevel(log.DEBUG)
//net.SetDebugLevel(log.DEBUG)
//h.Res.Write([]byte(fmt.Sprintf("loglevel is %v now!\n", log.DEBUG)))
//return
//case "info":
//s.Log.SetLevel(log.INFO)
//net.SetDebugLevel(log.INFO)
//h.Res.Write([]byte(fmt.Sprintf("loglevel is %v now!\n", log.INFO)))
//return
//case "warn":
//s.Log.SetLevel(log.WARN)
//net.SetDebugLevel(log.WARN)
//h.Res.Write([]byte(fmt.Sprintf("loglevel is %v now!\n", log.WARN)))
//return
//case "error":
//s.Log.SetLevel(log.ERROR)
//net.SetDebugLevel(log.ERROR)
//h.Res.Write([]byte(fmt.Sprintf("loglevel is %v now!\n", log.ERROR)))
//return
//}
//
//switch h.Req.FormValue("stack") {
//case "main":
//buf := make([]byte, 8*1024)
//size := runtime.Stack(buf, false)
//h.Res.Write(buf[:size])
//return
//case "all":
//buf := make([]byte, 10*1024*1024)
//size := runtime.Stack(buf, true)
//h.Res.Write(buf[:size])
//return
//}
//
//switch h.Req.FormValue("trace") {
//case "start":
//if s.Perf.StartTrace() == nil {
//h.Res.Write([]byte(fmt.Sprintf("trace is started now!\n")))
//} else {
//h.Res.Write([]byte(fmt.Sprintf("trace is started already!\n")))
//}
//case "stop":
//s.Perf.StopTrace()
//h.Res.Write([]byte(fmt.Sprintf("trace is stopped now!\n")))
//}
//}
//
////实现Service接口：处理定时器超时通知
//func (s *BaseService) ProcessTimer(tn *timer.TimeoutNotify) {
//s.Log.Debugf("TIMEOUT", "got time notify %#v", tn)
//}
//
////实现Service接口:其他每循环要执行的update
//func (s *BaseService) Update() {
////s.Log.Debugf("DEBUG","Update")
//log.FlushAll()
//}
//
////实现Service接口:重载
//func (s *BaseService) OnReload() {
//s.Log.Infof("CFG", "reload all begin")
//s.Cfg.ReLoadAll()
//}
//
////实现Service接口:退出
//func (s *BaseService) OnExit() {
//s.Log.Infof(LogTag, "%s exit", s.Proc.Name)
//log.FlushAll()
//}
//
////实现Service接口：网络连接异常断开
//func (s *BaseService) OnNetDisconn(conn *net.Conn) {
//
//}

//实现Service接口:主循环函数
func (self *BaseService) MainLoop() {
	for {
		select {
//			case <-s.ReloadChan: //重载配置
//				s.FW.Service.OnReload()
			case <-s.ExitChan: //进程退出
				s.FW.Service.OnExit()
				return
//			case http := <-s.HttpChan: //http命令
//				s.FW.Service.ProcessHttpCmd(http)
//				s.HttpChan <- nil
//			case tn := <-s.TNChan: //定时器超时
//				s.FW.Service.ProcessTimer(tn)
		}
	}
}

////实现Service接口:注册消息处理
//func (s *BaseService) RegisterMsgHandle() {
//}
//
////实现Service接口:注册一个消息处理
////参数msgId:消息id
////参数handle:处理消息的handle，传入的是指针
//func (s *BaseService) RegOneMsgHandle(msgId uint32, handle framework.MsgHandle) (int, error) {
//return s.MP.Reg(msgId, handle)
//}
//
////发送回复消息
//func (s *BaseService) SendResponseToClient(header *msg.MsgHeader, msgId uint32, body proto.Message) (int, error) {
//m := new(msg.Msg)
//m.Header = new(msg.MsgHeader)
//
//var connID uint32
//nh := header.GetNHeader()
//if nh != nil {
//msg.ConstructResHeader(header, m.Header, msgId, true)
//if s.Svr.IsGateWay() {
//connID = nh.GetSessionID()
//} else {
//connID = nh.GetServiceID()
//}
//s.Log.Debugf(LogTag, "net head is not nil,connID = %v", connID)
//} else {
//msg.ConstructResHeader(header, m.Header, msgId, false)
//connID = m.GetHeader().GetSHeader().GetDstServiceID()
//s.Log.Debugf(LogTag, "net head is nil,connID = %v", connID)
//}
//
////编码消息体
//var e error
//m.Body, e = proto.Marshal(body)
//if e != nil {
//return -1, e
//}
//
////编码消息
//var buff []byte
//buff, e = proto.Marshal(m)
//if e != nil {
//return -1, e
//}
//
////找到连接并发送
//conn := s.CM.GetConn(connID)
//if conn == nil {
//return -1, errors.New("can't find connection by service id %d", connID)
//}
//
//e = conn.SendNormalMsg(buff, 0)
//if e != nil {
//return -1, e
//}
//
//return 0, nil
//}
//
////获取服务id
//func (s *BaseService) GetServiceID() uint32 {
//return s.ServiceID
//}
//
////实现ProcessCommander接口:响应信号重载命令
//func (s *BaseService) SignalReload() {
//s.ReloadChan <- 0
//}
//
////实现ProcessCommander接口:响应信号退出命令
//func (s *BaseService) SignalQuit() {
//s.ExitChan <- 0
//}
//
////实现ProcessCommander接口:响应信号退出命令
//func (s *BaseService) GetHttpChan() chan *process.HttpContext {
//return s.HttpChan
//}
