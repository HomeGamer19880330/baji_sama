#coding=utf-8
#!/usr/bin/python

import os, sys
import globalConfig
import util

def compileService(serviceName):
	print ">>> building %s <<<"%(serviceName)

	os.chdir(os.path.join(globalConfig.g_service_dir, serviceName, "src"))
	cmd = "go build -o " + os.path.join(globalConfig.g_service_dir, serviceName, "bin", serviceName+".exe")
	os.system(cmd)
	os.chdir(globalConfig.g_current_dir)

def buildAll():
	for eachService in globalConfig.SERVICES_LIST:
		compileService(eachService)

def cleanAll():
	for eachService in globalConfig.SERVICES_LIST:
		cleanService(eachService)

def cleanService(serviceName):
	os.chdir(os.path.join(globalConfig.g_service_dir, serviceName, "src"))
	cmd = "go clean" 
	os.system(cmd)

	util.delete_file(os.path.join(globalConfig.g_service_dir, serviceName, "bin", serviceName))

	os.chdir(globalConfig.g_current_dir)

# {
#     while read d; do
#         cd $d/src
#         go clean
#         rm -f ../bin/$d
#         cd - > /dev/null 
# 		#cd - 回到历史目录
# 		#> /dev/null 将记录保存到/dev/null文件中,这个应该是保存一下记录没啥卵用
#     done < services.txt
# }

# function sigout()
# {
#     echo
#     echo ">>> build job has been terminated by `whoami`"
#     echo
#     exit 0
# }

# trap sigout SIGINT
# case语句,类似于其他语言中的switch
# read aNum
# case $aNum in

# $+数字 一般是位置参数的用法。
# 如果运行脚本的时候带参数，那么可以在脚本里通过 $1 获取第一个参数，$2 获取第二个参数......依此类推，一共可以直接获取9个参数（称为位置参数）。$0用于获取脚本名称。
# 相应地，如果 $+数字 用在函数里，那么表示获取函数的传入参数，$0表示函数名
# case $1 in
#     clean) #--相当于case
#         clean
#         ;; #--相当于break
#     "")
#         build
#         ;;
#     *)
#         build $@
# esac

if __name__ == "__main__":
	buildAll()
