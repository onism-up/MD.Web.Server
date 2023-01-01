package runCore

import (
	"MDServer/baseLoad"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2/utils"
	"os"
	"path"
	"strings"
	"time"
)

// Export 对本地全部数据进行导出
func Export(filepath string) {
	file, err := os.Create(path.Join(filepath, "md_file.json"))
	if err != nil {
		panic("数据导出失败: " + err.Error())
	} else {
		exportBase := ExportFileType{User: baseLoad.UserInfo, Base: baseLoad.MDBase}
		byteBase, _ := json.Marshal(exportBase)
		_, err = file.Write(byteBase)
		if err != nil {
			println("数据写入失败: ", err.Error())
		} else {
			println("导出成功")
		}
		defer file.Close()
	}
}

// ExportMD 对本地MD数据进行导出
func ExportMD(filepath string) {
	file, err := os.Create(path.Join(filepath, "md_file_data.json"))
	if err != nil {
		panic("数据导出失败: " + err.Error())
	} else {
		exportBase := baseLoad.MDBase
		byteBase, _ := json.Marshal(exportBase)
		_, err = file.Write(byteBase)
		if err != nil {
			panic("数据写入失败: " + err.Error())
		}
		defer file.Close()
	}
}

// Import 对外部资源进行引入
func Import(filepath string) {
	file, err := os.ReadFile(filepath)
	if err != nil {
		panic("目标文件读取失败: " + err.Error())
	}
	exportFile := new(ExportFileType)
	err = json.Unmarshal(file, exportFile)
	if err != nil {
		panic("导入时文件解析失败: " + err.Error())
	}
	err = DistributionImport(*exportFile, true)
	if err != nil {
		println("导入过程中出现的错误: " + err.Error())
	}
}

// DistributionImport 对导入的数据进行分发存储
func DistributionImport(exportFile ExportFileType, danger bool) error {
	var resMsg string
	errMd := MDImport(exportFile.Base)
	errUs := SaveUser(exportFile.User, danger)
	if errMd != nil {
		resMsg += errMd.Error()
	}
	if errUs != nil {
		resMsg += errUs.Error()
	}

	if resMsg != "" {
		return errors.New(resMsg)
	}

	return nil
}

// MDImport MD数据的导入
func MDImport(exportFile baseLoad.MDBaseType) error {
	var errMsg strings.Builder
	for _, v := range exportFile {
		err := SaveMd(v)
		if err != nil {
			errMsg.WriteString(fmt.Sprint("文档名为: ", v.Title, " ID为: ", v.Id, " 的文档导入失败\n"))
		}
	}
	relMsg := errMsg.String()
	if relMsg != "" {
		return errors.New(relMsg)
	}
	return nil
}

// Init 对资源进行重新初始化
func Init() {
	if _, err := os.ReadDir(baseLoad.DefaultDataPATH); err != nil {
		err := os.RemoveAll(baseLoad.DefaultDataPATH)
		if err != nil {
			panic("数据文件重置失败: " + err.Error())
		}
		baseLoad.MDBase = baseLoad.MDBaseType{}
	}
	UserInit()
	baseLoad.Init()

}

// UserInit 用户数据初始化
func UserInit() {
	_, err := os.ReadFile(baseLoad.DefaultUserInfoPATH)
	if err != nil {
		panic("用户文件读取失败: " + err.Error())
	}
	err = os.Remove(baseLoad.DefaultUserInfoPATH)
	if err != nil {
		panic("用户数据文件重置失败: " + err.Error())
	}

	baseLoad.UserInfo = baseLoad.UserInterface{}
}

// SaveUser 用户数据的存储
func SaveUser(userInfo baseLoad.UserInterface, interrupt bool) error {
	filepath := baseLoad.DefaultUserInfoPATH
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		errMsg := "用户数据读取失败: " + err.Error()
		if interrupt {
			panic(errMsg)
		} else {
			return errors.New(errMsg)
		}
	}
	byteFile, _ := json.Marshal(userInfo)
	_, err = file.Write(byteFile)
	if err != nil {
		errMsg := "用户数据写入失败: " + err.Error()
		if interrupt {
			panic(errMsg)
		} else {
			return errors.New(errMsg)
		}
	} else {
		baseLoad.UserInfo = userInfo //加载到缓存中
		return nil
	}
	defer file.Close()
	return nil
}

// UpdateUserToken 用户Token更新
func UpdateUserToken(token string) error {
	baseLoad.UserInfo.Token = token
	return SaveUser(baseLoad.UserInfo, false)
}

// CreteMD 创建新的MD并载入缓存
func CreteMD(title string) (string, error) {
	newMdFile := baseLoad.MDInterface{}
	newMdFile.Title = title
	newMdFile.Id = utils.UUIDv4()
	nowTime := time.Now().String()
	newMdFile.CreateTime = nowTime
	newMdFile.LastChangTime = nowTime
	return newMdFile.Id, SaveMd(newMdFile)
}

// SaveMd 保存并更新文件然后载入缓存
func SaveMd(mdFile baseLoad.MDInterface) error {
	path := path.Join(baseLoad.DefaultDataPATH, fmt.Sprint(mdFile.Id, ".json"))
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	fileByte, _ := json.Marshal(mdFile)
	_, err = file.Write(fileByte)
	if err != nil {
		return err
	}
	defer file.Close()
	baseLoad.MDBase[mdFile.Id] = mdFile
	return nil
}

type MDFileListType = []MDFileType

type MDFileType = struct {
	Id        string `json:"id"`
	Title     string `json:"title"`
	CreatTime string `json:"creat_time"`
}

func ReadMDFileList(isPub bool) MDFileListType {
	MDFileList := MDFileListType{}
	for _, v := range baseLoad.MDBase {
		if isPub && !v.IsPublic {
			continue
		}
		newMDFile := MDFileType{}
		newMDFile.Title = v.Title
		newMDFile.CreatTime = v.CreateTime
		newMDFile.Id = v.Id
		MDFileList = append(MDFileList, newMDFile)
	}
	return MDFileList
}

// DelMD 以ID为标准删除某个MD文件并同步到缓存
func DelMD(id string) error {
	_, y := baseLoad.MDBase[id]
	if !y {
		return errors.New("id is not curren")
	}
	filePATH := path.Join(baseLoad.DefaultDataPATH, fmt.Sprint(id, ".json"))
	err := os.Remove(filePATH)
	delete(baseLoad.MDBase, id)
	return err
}

type SearchResult []SearchResultItem

type SearchResultItem struct {
	Title         SearchItem `json:"title"`
	Body          SearchItem `json:"body"`
	Id            string     `json:"id"`
	CreateTime    string     `json:"create_time"`
	LastChangTime string     `json:"last_chang_time"`
}

type SearchItem struct {
	Data           []string `json:"data"`
	HighlightIndex []int    `json:"highlight_index"`
}

func Search(key string, pub bool) SearchResult {
	newSeachResult := SearchResult{}
	for _, v := range baseLoad.MDBase {
		if pub && !v.IsPublic {
			continue
		}

		title := v.Title
		body := v.Body

		resTitle := new(SearchItem)
		resBody := new(SearchItem)

		titleParsingResult := SearchParsing(title, key, resTitle)

		bodyParsingResult := false

		if len(body) > 0 {
			bodyParsingResult = SearchParsing(body, key, resBody)
		}

		if titleParsingResult || bodyParsingResult {

			newItem := SearchResultItem{}
			newItem.Title = *resTitle
			newItem.Body = *resBody
			newItem.Id = v.Id
			newItem.CreateTime = v.CreateTime
			newItem.LastChangTime = v.LastChangTime

			newSeachResult = append(newSeachResult, newItem)
		}
	}
	return newSeachResult
}

func SearchParsing(origin, key string, item *SearchItem) bool {
	if strings.Contains(origin, key) {
		res := strings.Split(origin, key)
		lenRes := len(res)
		for index, iv := range res {
			item.Data = append(item.Data, iv)
			if index != lenRes-1 {
				item.Data = append(item.Data, key)
				item.HighlightIndex = append(item.HighlightIndex, len(item.Data)-1)
			}
		}
		return true
	} else {
		item.Data = append(item.Data, origin)
		return false
	}
}
