package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"icpcoll/api/miit"
	"icpcoll/logger"
	"icpcoll/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func TestICP_GetImage() {
	client := miit.NewClient()
	err := client.SetTokenFromRemote()
	if err != nil {
		logger.Info(err)
		return
	}
	// logger.Info(client.Token, client.RefreshToken, client.ExpireIn)
	result, err := client.PageNum(1).PageSize(1).ServiceType("1").Query("baidu.com")
	if err != nil {
		logger.Info(err)
		return
	}
	res, _ := json.Marshal(result)
	logger.Info(string(res))
	fmt.Println(string(res))

}

func BatchQuery(unitList []string, serviceTypeList []string) {
	var client *miit.ICP
	for i, unitName := range unitList {
		for j, serviceType := range serviceTypeList {
			result, err := client.PageNum(1).PageSize(1).ServiceType(serviceType).Query(unitName)
			if err != nil {
				logger.Info(err.Error())
			}
			if i == len(unitList)-1 && j == len(serviceTypeList)-1 {
				logger.Info(result.Items)
			} else {
				logger.Info(result.Items)
			}
		}
	}
}

// 定义一个自定义错误类型
type CustomError struct {
	msg string // 错误消息
}

// 实现 error 接口
func (ce *CustomError) Error() string {
	return ce.msg
}

func Export(unitName string, serviceType string) error {

	dataDir := "output"
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		//先创建output目录
		_ = os.MkdirAll(dataDir, 0750)
	}
	filename := fmt.Sprintf("%s_%d.xlsx", unitName, utils.GenUnixTimestamp())
	outputAbsFilepath := filepath.Join(dataDir, filename)

	// var client *miit.ICP
	client := miit.NewClient()
	err := client.SetTokenFromRemote()
	if err != nil {
		logger.Info(err.Error())
		return err
	}

	items := make([]*miit.Item, 0)
	result, err := client.PageNum(1).PageSize(1).ServiceType(serviceType).Query(unitName)
	if err != nil {
		logger.Info(err.Error())
		return err
	}
	rt, _ := json.Marshal(result)
	logger.Info(string(rt))
	fmt.Printf("[-]Total: %d\n", result.Total)
	if result.Total == 0 {
		customErr := &CustomError{msg: "[x] 未查询到相关信息!"}
		return customErr
	}
	result, err = client.PageNum(1).PageSize(result.Total).ServiceType(serviceType).Query(unitName)
	// result, err = client.PageNum(1).PageSize(10).ServiceType("1").Query(unitName)
	// 保存所有 域名信息到 temp目录，以 unitName 命名
	if err != nil {
		logger.Info(err.Error())
		return err
	}
	res, _ := json.Marshal(result)
	logger.Info(string(res))
	// fmt.Println(string(res))

	// 保存至 temp 目录
	utils.SaveToTempfile(unitName, result)

	items = append(items, result.Items...)
	var data [][]any
	headers := append([]any{"序号"}, []any{"企业名称", "备案内容",
		"备案号", "备案类型", "备案法人", "单位性质", "审核日期"}...)
	data = append(data, headers)
	for index, item := range items {
		var tmpItem = []any{
			index + 1,
			item.UnitName,
			item.ServiceName,
			item.ServiceLicence,
			item.ServiceType,
			item.LeaderName,
			item.NatureName,
			item.UpdateRecordTime}
		data = append(data, tmpItem)
	}
	if err := utils.SaveToExcel(data, outputAbsFilepath); err != nil {
		logger.Info(err.Error())
		// r.downloadLog.UpdateStatus(fileID, constraint.Statuses.Error, err.Error())
		return err
	}
	return nil
}

func main() {

	var (
		apimode      string // api接口来源
		serviceType  string // 1,6,7,8 网站，app，小程序，快应用
		unitName     string // 域名或单位名称
		unitFileName string
	)

	// TestICP_GetImage()

	flag.StringVar(&apimode, "a", "miit", "api接口来源")
	flag.StringVar(&serviceType, "s", "1", "1,6,7,8 网站，app，小程序，快应用")
	flag.StringVar(&unitName, "n", "", "域名或单位名称")
	flag.StringVar(&unitFileName, "f", "", "单位名称文件")
	flag.Parse()

	if unitName == "" && unitFileName == "" {
		fmt.Println("Error: At least one of -n or -f must be provided.")
		flag.Usage()
		return
	} else if unitName != "" {
		err := Export(unitName, serviceType)
		if err != nil {
			logger.Info(err)
			fmt.Println(err)
		}
	} else if unitFileName != "" {
		content, err := ioutil.ReadFile(unitFileName)
		if err != nil {
			fmt.Printf("Error reading file '%s': %v\n", unitFileName, err)
			return
		}
		lines := string(content)
		for _, line := range strings.Split(lines, "\n") {
			trimmedLine := strings.TrimSpace(line) // 去除行首尾的空白字符
			if trimmedLine != "" {
				fmt.Printf("查询信息: %s ", trimmedLine)
				err := Export(trimmedLine, serviceType)
				if err != nil {
					logger.Info(err)
					fmt.Println(err)
				}
			}
		}
	} else {
		return
	}
	// err := Export("北京百度网讯科技有限公司")
	// if err != nil {
	// 	logger.Info(err)
	// 	fmt.Println(err)
	// }
}
