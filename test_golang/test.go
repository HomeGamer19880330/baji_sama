package main

import (
	"fmt"
	// "os"
	"framework"
)

// func init() {

// }

// type FrameWork struct {
// 	Service ServiceInterface //服务接口
// }

// //服务接口
// type ServiceInterface interface {
// 	Init(*FrameWork) (int, error) //
// 	SetupNetwork() (int, error)
// }

type ShitService struct {
	// framework.BaseService
	// mu       sync.Mutex
	isClient bool
	//	versionCfgMap map[string]*config.GameVerInfoCfg
}

func (self *ShitService) LogService() {
	fmt.Println("do in testf")
}

// //实现Service接口:初始化函数
// func (self *BaseService) Init(mainFrameWrok_ *FrameWork) (int, error) {
// 	self.SetupNetwork()
// 	return 0, nil
// }

// func (self *BaseService) SetupNetwork() (int, error) {
// 	return 0, nil
// }

// //设置服务接口
// func (self *FrameWork) SetService(s ServiceInterface) {
// 	self.Service = s
// }

// type childService struct {
// 	// BaseService
// }

// import (
// 	"ximigame.com/framework"
// )

func main() {
	testFrameWork := new(framework.MainFrameWork)
	testFrameWork.ShitIt(1)
	//	ooshit := new(ShitService)
	//	ooshit.LogService()
	// fmt.Println("do in testf")
	//注册服务
	// testFrameWork
	// testFrameWork.SetService(ooshit)

	// //启动框架
	// fw.Run()
}
