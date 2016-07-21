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

	serviceInstance.isClient = false
	//注册服务
	servceFramework.Service = serviceInstance.(framework.BaseService*)

	//启动框架
	servceFramework.Run()
}
