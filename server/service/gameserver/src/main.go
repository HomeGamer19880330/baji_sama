package main

import (
	"fmt"
	"os"
)

import (
	"framework"
)

func main() {

	//获取框架实例
	servceFramework := framework.Instance()
	if servceFramework == nil {
		fmt.Println("get instance failed")
		os.Exit(-1)
	}

	serviceInstance.IsClient = false
	//注册服务
	servceFramework.SetService(serviceInstance)
	// Service = serviceInstance
	//.(*framework.ServiceInterface)
	//ServiceInterface    TestGameService(framework.BaseService)
	//启动框架
	servceFramework.Run()
}
