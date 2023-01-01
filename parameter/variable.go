package parameter

var Port *int          //暴露的接口
var DefaultPort = 8171 //默认的服务器接口
var Export *string     //是否是导出模式
var Import *string     //是否导入数据据
var Server *bool       //是否是服务器模式
var Help *bool         //是否需要帮助
var Init *bool         //重新进行初始化
var UserReset *bool    //账号重置但是数据保存
var ExportMD *string   //暴露出数据的位置

var helpMessage = `
-------------------------帮助---------------------------
-port <int>:   				指定主机对外暴露的接口，默认为8171
-export <string>: 			导出模式，会导出到目标文件服务器的全部数据,如不指定位置不会开启此模式
-import <string>: 			导入模式，需要指定导入文件路径，会导入目标文件中的所有数据，如不指定位置不会开启此模式
-server [boolean]: 			服务器模式，将会启动服务器
-init [boolean]:  			初始化，清除一切缓存数据
-export_md <string>			导出模式，只导出MD文档的数据，如不指定位置不会开启此模式
-user_reset <boolean>			重置用户数据
提示: 只有-port参数可以与其他参数并存，如果是其他参数则只能存在一个
-------------------------------------------------------
`
