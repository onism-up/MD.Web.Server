package baseLoad

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2/utils"
	"log"
	"os"
	"path"
)

func Init() bool {
	return InitDataBase() && InitUserInfo() && InitTokenKey()
}

func InitDataBase() bool {
	if fileDirectory, err := os.ReadDir(DefaultDataPATH); err != nil {
		if err = os.Mkdir(DefaultDataPATH, 0666); err != nil {
			panic("数据初始化失败: " + err.Error())
		} else {
			log.Println("数据初始化完成")
			return true
		}
	} else {
		for _, fileName := range fileDirectory {
			filePath := path.Join(DefaultDataPATH, fileName.Name())
			loadBase(filePath)
		}
		log.Println("数据载入成功")
		return true
	}
}

func loadBase(filePath string) {
	if file, err := os.ReadFile(filePath); err != nil {
		panic("文件读取失败: " + err.Error())
	} else {
		newBaseItem := new(MDInterface)
		if err = json.Unmarshal(file, newBaseItem); err != nil {
			panic("文件解析失败: " + err.Error())
		} else {
			MDBase[newBaseItem.Id] = *newBaseItem
		}
	}
}

func InitUserInfo() bool {
	if file, err := os.ReadFile(DefaultUserInfoPATH); err != nil {
		if file, err := os.Create(DefaultUserInfoPATH); err != nil {
			panic("用户数据初始化失败: " + err.Error())
		} else {
			userByte, _ := json.Marshal(UserInfo)
			file.Write(userByte)
			defer file.Close()
			log.Println("用户数据初始化完成")
			return true
		}
		return false
	} else {
		newUserInfo := new(UserInterface)
		if err = json.Unmarshal(file, newUserInfo); err != nil {
			panic("用户数据文件解析失败: " + err.Error())
		} else {
			log.Println("用户信息载入成功")
			UserInfo = *newUserInfo
			return true
		}
	}
}

func InitTokenKey() bool {
	v, err := os.ReadFile(TokenKeyPATH)
	if err != nil {
		file, err := os.Create(TokenKeyPATH)
		if err != nil {
			panic("token密匙创建失败: " + err.Error())
		} else {
			_, err = file.Write([]byte(utils.UUIDv4()))
			if err != nil {
				panic("token密匙写入失败: " + err.Error())
			} else {
				return true
			}
		}
	} else {
		TokenKey = v
		log.Println("token密匙载入成功")
		return true
	}
}
