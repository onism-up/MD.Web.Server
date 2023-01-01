package parameter

import (
	"MDServer/baseLoad"
	"MDServer/runCore"
	"flag"
	"log"
	"time"
)

func InitPar() bool {
	Export = flag.String("export", "", "The export parameters expect a string data type")
	Import = flag.String("import", "", "The import parameters expect a string data type")
	Port = flag.Int("port", DefaultPort, "The port parameters expect a int data type")
	Server = flag.Bool("server", false, "The server parameters expect a boolean data type")
	Help = flag.Bool("help", false, "The help parameters expect a boolean data type")
	Init = flag.Bool("init", false, "The init parameters expect a boolean data type")
	ExportMD = flag.String("export_md", "", "The export_md parameters expect a string data type")
	UserReset = flag.Bool("user_reset", false, "The user-reset parameters expect a boolean data type")
	flag.Parse()

	if checkPar() == false {
		log.Println("没有输入任何命令，或您输入的命令不合法，寻求帮助请使用参数 -help")
		return false
	}

	baseLoad.Init()

	if len(*Export) > 0 { //判断是否为导出模式，如果是则取消服务器的启动
		runCore.Export(*Export)
		log.Println("数据已全部导出")
		return false
	} else if *Server { //服务器启动
		log.Println("服务器模式将会在3秒后启动")
		time.Sleep(3 * time.Second)
		return true
	} else if *Help { //展示帮助
		println(helpMessage)
		return false
	} else if *Init { //初始化
		runCore.Init()
		log.Println("系统已进行初始化")
		return false
	} else if len(*Import) > 0 { //是否是导入模式
		runCore.Import(*Import)
		log.Println("数据已导入")
		return false
	} else if *UserReset {
		runCore.UserInit()
		log.Println("用户信息已重置")
		return false
	} else if len(*ExportMD) > 0 {
		runCore.ExportMD(*ExportMD)
		log.Println("文档数据已导出")
	}
	return false
}

func checkPar() bool {
	productParNum := 0
	boolArr := []*bool{Server, Help, Init, UserReset}

	for _, value := range boolArr {
		if *value {
			productParNum++
		}
	}

	if len(*Export) > 0 {
		productParNum++
	}

	if len(*Import) > 0 {
		productParNum++
	}

	if len(*ExportMD) > 0 {
		productParNum++
	}

	if productParNum != 1 {
		return false
	}

	return true
}
