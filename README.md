# 🎆 MD.WEB.Back

GoLang编写，GoFiber搭建 私有 Markdown 笔记应用（后台部分）,核心代码已被单独封装。

- 在线预览：[建设中](https://gitee.com)

- 同系列前端：[传送门](https://gitee.com/FM107/MD.Web)

## ⭐️ 功能一览

- 快速搭建私人服务器部署的笔记应用，也可以当其为轻型博客；
- 内置的白色主题和暗黑主题，支持动态切换；
- 支持快速在线对程序中的文档进行增删改查；
- 轻松和快速的分享你的创意；
- 完善的 MD 语法支持；
- 无需数据库支持，本应用目前应是无依赖数据库版本；
- 较强的可移植性和稳定性，可批量导入文档；

## 📦 安装

1. 直接运行
   1. 下载对应系统的包，然后在终端中执行，按照命令提示进行操作
2. 获取源码
   1. 克隆或者下载到本地
   2. 终端切换到根目录执行`go install`
   3. 执行 `go run main.go` 测试
   4. 执行 `go build` 打包，详细命令请参考 [官方文档](https://go.dev/doc/)
## 🗺 预览图

| 默认模式(读者主页)                                           | 暗黑模式                                                    |
| ------------------------------------------------------------ | ----------------------------------------------------------- |
| ![](https://pic.imgdb.cn/item/63b16e282bbf0e7994722929.jpgg) | ![](https://pic.imgdb.cn/item/63b16e952bbf0e799472b40c.jpg) |

| 写者主页                                                    | 编辑模式                                                    |
| ----------------------------------------------------------- | ----------------------------------------------------------- |
| ![](https://pic.imgdb.cn/item/63b16e952bbf0e799472b40c.jpg) | ![](https://pic.imgdb.cn/item/63b170fc2bbf0e799475cdcf.jpg) |

## ✨ 小提示

1. 非写者只能看到已公开的文档
2. 导入时要注意是否是之前自己之前导出的部分和现有的部分重合，重合的文档将以导入的文档内容为主进行覆盖
3. 目前服务器还不支持 https，密码为明文传输，一定要注意网络安全，如浏览器插件和非安全网络环境等，防止数据泄露

## 📝 开发计划

1. 将服务器升级为 https
2. 增加评论系统
3. 实现国际化
4. 主题自定义
