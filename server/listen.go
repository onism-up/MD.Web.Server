package server

import (
	"MDServer/parameter"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// Listen 监听本机端口，并提供服务
func Listen() {
	app := fiber.New()

	token := NewToken()

	token.Use(checkLocalhostToken) //启用token检测

	apiApp := app.Group("/api", cors.New()) //url嵌套并为此开启跨域

	apiApp.Use(token.Handler())

	apiApp.Get("/test", test) //测试服务器是否运行，此接口始终返回true

	apiApp.Post("/register", register) //注册

	apiApp.Post("/login", login) //登录

	apiApp.Post("modify_user", ChangeUserInfo) //修改

	apiApp.Get("/validation", test) //检验token是否正确，此接口始终返回true

	apiApp.Post("/file", CreateMdFile) //创建新的文件

	apiApp.Get("/file", GetDetailMD(false)) //获取文件详细信息

	apiApp.Get("/pub_file", GetDetailMD(true)) //获取公开的文件信息

	apiApp.Patch("/file", SaveMD) //保存文件

	apiApp.Delete("/file", DelMD) //删除文件

	apiApp.Get("/file_list", GetFileList) //获得文件列表

	apiApp.Get("/pub_file_list", GetPubFileList) //获取公开文件列表

	apiApp.Post("/upload", FileReceive) //全服数据覆盖式导入

	apiApp.Post("/upload_append", MDReceive) //仅文档数据覆盖式导入

	apiApp.Get("/search", GetSearch(false)) //实现搜索所有数据

	apiApp.Get("/pub_search", GetSearch(true)) //搜索仅公开数据

	apiApp.Get("/user_init", UserInit)

	err := app.Listen(fmt.Sprint(":", *parameter.Port)) //鉴定解析到的端口或者默认端口
	if err != nil {
		panic("服务器监听失败: " + err.Error())
	}
}
