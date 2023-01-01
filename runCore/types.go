package runCore

import "MDServer/baseLoad"

// ExportFileType 导出和导入的文件类型定义
type ExportFileType struct {
	User baseLoad.UserInterface `json:"user"`
	Base baseLoad.MDBaseType    `json:"base"`
}
