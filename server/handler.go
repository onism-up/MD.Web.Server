package server

import (
	"MDServer/baseLoad"
	"MDServer/runCore"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"io"
	"time"
)

// checkLocalhostToken 检查本地的token是否和解析出的token正确
func checkLocalhostToken(ctx *fiber.Ctx, _, originToken string) error {
	if originToken == baseLoad.UserInfo.Token {
		ctx.Next()
		return nil
	} else {
		if len(baseLoad.UserInfo.Token) == 0 {
			return errors.New("2:user reset")
		}
		return errors.New("1:token error")
	}
}

func test(ctx *fiber.Ctx) error {
	Message.Success("success", ctx)
	return nil
}

type RegisterResponseType struct {
	Name string `json:"name"`
	Pwd  string `json:"pwd"`
}

func UserDataCheck(userInfo RegisterResponseType) error {
	nameLen := len(userInfo.Name)
	pwdLen := len(userInfo.Pwd)
	if nameLen < 3 || nameLen > 10 {
		return errors.New("6:The name field must be between 3 to 10")
	}
	if pwdLen < 6 || nameLen > 18 {
		return errors.New("6:The pwd field must be between 6 to 18")
	}
	return nil
}

func register(ctx *fiber.Ctx) error {

	if checkHaveBody(ctx) {
		return FasterError(ctx, "4:need body")
	}
	newReg := new(RegisterResponseType)
	err := json.Unmarshal(ctx.Body(), newReg)

	if err != nil {
		return FasterError(ctx, "5:parsing error: "+err.Error())
	}

	err = UserDataCheck(*newReg)
	if err != nil {
		return FasterError(ctx, err.Error())
	}

	newUser := baseLoad.UserInterface{}
	newUser.Name = newReg.Name
	newUser.Pwd = newReg.Pwd
	newUser.Init = true
	token := createRandomToken()
	newUser.Token = token
	err = runCore.SaveUser(newUser, false)
	if err != nil {
		return FasterError(ctx, "12:save user error: "+err.Error())
	}
	Message.Success(token, ctx)
	return nil
}

func login(ctx *fiber.Ctx) error {

	if baseLoad.UserInfo.Init == false {
		return FasterError(ctx, "7:user is not init")
	}

	if checkHaveBody(ctx) {
		return FasterError(ctx, "4:need body")
	}
	newReg := new(RegisterResponseType)
	err := json.Unmarshal(ctx.Body(), newReg)

	if err != nil {
		return FasterError(ctx, "5:parsing error: "+err.Error())
	}

	err = UserDataCheck(*newReg) //检查用户账号密码的数据合法性
	if err != nil {
		return FasterError(ctx, err.Error())
	}

	if newReg.Name != baseLoad.UserInfo.Name || newReg.Pwd != baseLoad.UserInfo.Pwd {
		return FasterError(ctx, "8:user or pwd error")
	}

	newToken := createRandomToken()

	runCore.UpdateUserToken(newToken)

	Message.Success(newToken, ctx)

	return nil
}

func checkHaveBody(ctx *fiber.Ctx) bool {
	return len(ctx.Body()) == 0
}

func createRandomToken() string {
	token, err := MakeToken(utils.UUIDv4())
	if err != nil {
		println(baseLoad.TokenKey)
		println(err.Error())
	}
	return token
}

func ChangeUserInfo(ctx *fiber.Ctx) error {
	if checkHaveBody(ctx) {
		return FasterError(ctx, "4:need body")
	}
	newReg := new(RegisterResponseType)
	err := json.Unmarshal(ctx.Body(), newReg)

	if err != nil {
		return FasterError(ctx, "5:parsing error: "+err.Error())
	}

	err = UserDataCheck(*newReg)
	if err != nil {
		return FasterError(ctx, err.Error())
	}

	baseLoad.UserInfo.Name = newReg.Name
	baseLoad.UserInfo.Pwd = newReg.Pwd

	newToken := createRandomToken() //生成新的Token

	baseLoad.UserInfo.Token = newToken

	err = runCore.SaveUser(baseLoad.UserInfo, false)

	if err != nil {
		return FasterError(ctx, "12:save user error: "+err.Error())
	}

	Message.Success(newToken, ctx)

	return nil

}

type NewFileType struct {
	Title string `json:"title"`
}

func CreateMdFile(ctx *fiber.Ctx) error {
	if checkHaveBody(ctx) {
		return FasterError(ctx, "4:need body")
	}
	newFileTitle := new(NewFileType)
	err := json.Unmarshal(ctx.Body(), newFileTitle)

	if err != nil {
		return FasterError(ctx, "5:parsing error: "+err.Error())
	}

	lenTitle := len(newFileTitle.Title)
	if lenTitle <= 0 || lenTitle > 27 {
		return FasterError(ctx, "6:The title field length must be between 0 to 28")
	}

	fileId, err := runCore.CreteMD(newFileTitle.Title)

	if err != nil {
		return FasterError(ctx, "9:create new file error: "+err.Error())
	}

	Message.Success(fileId, ctx)

	return nil
}

func GetFileList(ctx *fiber.Ctx) error {
	Message.Success(runCore.ReadMDFileList(false), ctx)
	return nil
}

func GetPubFileList(ctx *fiber.Ctx) error {
	Message.Success(runCore.ReadMDFileList(true), ctx)
	return nil
}

func GetDetailMD(isPub bool) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Query("id", "")
		if id == "" {
			return FasterError(ctx, "10:need query id")
		}
		value, y := baseLoad.MDBase[id]

		if !y || (isPub && isPub != value.IsPublic) {
			return FasterError(ctx, "11:this id is not corresponding to the data")
		}
		Message.Success(value, ctx)
		return nil
	}
}

func SaveMD(ctx *fiber.Ctx) error {
	if checkHaveBody(ctx) {
		return FasterError(ctx, "4:need body")
	}
	MDFile := new(baseLoad.MDInterface)
	err := json.Unmarshal(ctx.Body(), MDFile)

	if err != nil {
		return FasterError(ctx, "5:parsing error: "+err.Error())
	}
	lenTitle := len(MDFile.Title)
	if lenTitle <= 0 || lenTitle > 27 {
		return FasterError(ctx, "6:The title field length must be between 0 to 28")
	}
	_, y := baseLoad.MDBase[MDFile.Id]
	if !y {
		return FasterError(ctx, "11:id is not corresponding to the data")
	}
	MDFile.LastChangTime = time.Now().String()

	err = runCore.SaveMd(*MDFile)

	if err != nil {
		return FasterError(ctx, "12:file save error: "+err.Error())
	}
	Message.Success(MDFile.Id, ctx)
	return nil
}

func DelMD(ctx *fiber.Ctx) error {
	id := ctx.Query("id", "")
	if id == "" {
		return FasterError(ctx, "10:need query id")
	}
	err := runCore.DelMD(id)
	if err != nil {
		return FasterError(ctx, fmt.Sprint("13:", err.Error()))
	}
	Message.Success("", ctx)
	return nil
}

func FileReceive(ctx *fiber.Ctx) error {
	if checkHaveBody(ctx) {
		return FasterError(ctx, "4:need body")
	}
	exportFile := new(runCore.ExportFileType)
	err := json.Unmarshal(ctx.Body(), exportFile)

	if err != nil {
		return FasterError(ctx, "5:parsing error: "+err.Error())
	}

	err = runCore.DistributionImport(*exportFile, false)

	if err != nil {
		return FasterError(ctx, "13:In the process of import for the mistakes: "+err.Error())
	}

	Message.Success("", ctx)
	return nil
}

func MDReceive(ctx *fiber.Ctx) error {
	if checkHaveBody(ctx) {
		return FasterError(ctx, "4:need body")
	}
	header, err := ctx.FormFile("file")
	fileSlice := make([]byte, 1024)
	file, err := header.Open()
	if err != nil {
		return FasterError(ctx, "15:File Read Error: "+err.Error())
	}
	var chunk []byte
	for {
		line, err := file.Read(fileSlice)
		if err != nil {
			if err == io.EOF {
				chunk = append(chunk, fileSlice[:line]...)
				break
			}
		} else {
			chunk = append(chunk, fileSlice[:line]...)
		}
	}

	exportFile := new(baseLoad.MDBaseType)
	err = json.Unmarshal(chunk, exportFile)

	if err != nil {
		return FasterError(ctx, "5:parsing error: "+err.Error())
	}
	err = runCore.MDImport(*exportFile)
	if err != nil {
		return FasterError(ctx, "13:In the process of import for the mistakes: "+err.Error())
	}

	Message.Success("", ctx)
	return nil
}

func GetSearch(isPublic bool) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		key := ctx.Query("key", "")
		if key == "" {
			return FasterError(ctx, "10:need query key")
		}
		Message.Success(runCore.Search(key, isPublic), ctx)
		return nil
	}
}

func UserInit(ctx *fiber.Ctx) error {
	if baseLoad.UserInfo.Init {
		Message.Success("", ctx)
		return nil
	} else {
		return FasterError(ctx, "")
	}
}
