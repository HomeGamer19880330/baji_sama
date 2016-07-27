package main

// import (
// "fmt"
//	"github.com/golang/protobuf/proto"
// "sync"
//	"ximigame.com/types/proto/config"
//	"ximigame.com/types/proto/external"
// )

import (
	"framework"
)

//const (
//	configRelativePath      = "config"
//	gameHallVersionInfoName = "gameHallVersion.cfg"
//)

type TestGameService struct {
	framework.BaseService
	// mu       sync.Mutex
	// isClient bool
	//	versionCfgMap map[string]*config.GameVerInfoCfg
}

var (
	serviceInstance = new(TestGameService) //创建服务实例
)

// esop-209
//重载注册配置函数
//func (s *TestGameService) RegisterCfg() (int, error) {
//	_, e := s.BaseService.RegisterCfg()
//	if e != nil {
//		s.Log.Errorf("hall_service", "register cfg failed, %s", e)
//		panic(e.Error())
//		return -1, e
//	}
//
//	s.versionCfgMap = make(map[string]*config.GameVerInfoCfg)
//	cfgName := "version." + "gameHallVersion"
//	res, e := s.Cfg.Register(
//		cfgName,
//		s.Proc.DefaultCfgDir+"/"+gameHallVersionInfoName,
//		true,
//		func() proto.Message { return new(config.GameHallVersionInfoCfg) },
//		nil,
//		func(cfgName string, msg proto.Message) (int, error) {
//			gameHallVersionInfo, ok := msg.(*config.GameHallVersionInfoCfg)
//			s.Log.Infof("VersionService", "VersionService GameHallVersionInfoCfg %v ", *gameHallVersionInfo)
//			if !ok {
//				panic("load GameHallVersionInfoCfg err")
//			} else {
//				s.mu.Lock()
//				defer s.mu.Unlock()
//				s.versionCfgMap = make(map[string]*config.GameVerInfoCfg)
//				for _, v := range gameHallVersionInfo.GetGameVerInfo() {
//					key := fmt.Sprintf("%d_%d", *v.GameId, *v.ChannelID)
//
//					_, ok := s.versionCfgMap[key]
//					if ok {
//						panic("gameid plus channelid is not unique")
//					}
//					s.versionCfgMap[key] = v
//				}
//			}
//			return 0, nil
//		})
//	if e != nil {
//		return res, e
//	}
//
//	return 0, nil
//}

//重载消息处理handle注册函数
//func (s *VersionService) RegisterMsgHandle() {
//	s.RegOneMsgHandle(uint32(external.MsgID_VER_CHECK_REQ), new(VerCheckHandle))
//}
