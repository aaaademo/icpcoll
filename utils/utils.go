package utils

import (
	"fmt"
	"icpcoll/api/miit"
	"icpcoll/logger"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
	"golang.org/x/net/html"
)

func RemoveEmptyAndDuplicateString(list []string) []string {
	uniqueMap := make(map[string]bool)
	var result []string
	for _, str := range list {
		v := strings.TrimSpace(str)
		if v != "" && !uniqueMap[v] {
			uniqueMap[v] = true
			result = append(result, v)
		}
	}
	return result
}

func RemoveEmptyStrings(slice []string) []string {
	result := make([]string, 0)
	for _, str := range slice {
		if str != "" {
			result = append(result, str)
		}
	}
	return result
}

func RemoveEmptyAndDuplicateAndJoinStrings(list []string, sep string) string {
	list = RemoveEmptyAndDuplicateString(list)
	if len(list) == 0 {
		return ""
	}
	return strings.Join(list, sep)
}

func GenFilenameTimestamp() string {
	formattedTime := time.Now().Format("2006-01-02-15-04-05")
	return formattedTime
}

func GenTimestampOutput() string {
	formattedTime := time.Now().Format("2006/01/02 15:04:05")
	return formattedTime
}

// GenUnixTimestamp 获取当前 Unix 时间戳（以秒为单位）
func GenUnixTimestamp() int64 {
	return time.Now().Unix()
}

func HtmlHasID(n *html.Node, id string) bool {
	for _, attr := range n.Attr {
		if attr.Key == "id" && attr.Val == id {
			return true
		}
	}
	return false
}

func columnNumberToName(n int) string {
	name := ""
	for n > 0 {
		n--
		name = string(byte(n%26)+'A') + name
		n /= 26
	}
	return name
}

func SaveToTempfile(unitName string, result *miit.Result) error {
	tempDir := "temp"
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		//先创建temp目录
		_ = os.MkdirAll(tempDir, 0750)
	}
	tempDirFull, _ := filepath.Abs(tempDir)
	tempFile := filepath.Join(tempDirFull, unitName+".txt")
	// fmt.Println(tempFile)

	// 提取所有 serviceName
	serviceNames := make([]string, 0)
	for _, item := range result.Items {
		serviceNames = append(serviceNames, item.ServiceName)
	}

	// 打印所有 serviceName
	// fmt.Printf("========= %s ===========\n", unitName)
	for _, serviceName := range serviceNames {
		if serviceName == "" {
			return nil
		}
		fmt.Println(serviceName)
		file, err := os.OpenFile(tempFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err == nil {
			defer file.Close()
		}
		output := fmt.Sprintf("%s\n", serviceName)
		// fmt.Printf(output)
		_, err = file.WriteString(output)
		if err != nil {
			fmt.Println(err)
		}
	}
	// fmt.Println("========================")
	return nil
}

func SaveToExcel(data [][]any, outputFilepath string) error {
	file := excelize.NewFile()

	// 添加数据
	for i := 0; i < len(data); i++ {
		row := data[i]
		startCell, err := excelize.JoinCellName("A", i+1)
		if err != nil {
			logger.Info(err.Error())
			return err
		}
		if i == 0 {
			// 首行大写
			for j := 0; j < len(row); j++ {
				if value, ok := row[j].(string); ok {
					row[j] = strings.ToUpper(value)
				}
			}
			if err = file.SetSheetRow("Sheet1", startCell, &row); err != nil {
				logger.Info(err.Error())
				return err
			}
			continue
		}
		if err = file.SetSheetRow("Sheet1", startCell, &row); err != nil {
			return err
		}
	}

	// 表头颜色填充
	headerStyle, err := file.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#d0cece"}, Pattern: 1, Shading: 1},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
		},
	})
	if err != nil {
		return err
	}

	err = file.SetCellStyle("Sheet1", "A1", columnNumberToName(len(data[0]))+"1", headerStyle)
	if err != nil {
		return err
	}

	// 添加边框
	dataStyle, err := file.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern"},
		Alignment: &excelize.Alignment{
			Horizontal: "left",
		},
		Border: []excelize.Border{
			{
				Type:  "right",
				Color: "#000000",
				Style: 1,
			},
			{
				Type:  "left",
				Color: "#000000",
				Style: 1,
			},
			{
				Type:  "top",
				Color: "#000000",
				Style: 1,
			},
			{
				Type:  "bottom",
				Color: "#000000",
				Style: 1,
			},
		},
	})
	if err != nil {
		return err
	}
	err = file.SetCellStyle("Sheet1", "A1", columnNumberToName(len(data[0]))+strconv.Itoa(len(data)), dataStyle)
	if err != nil {
		return err
	}

	if err2 := file.SaveAs(outputFilepath); err2 != nil {
		return err2
	}
	fmt.Printf("[+] xlsx file saved -> %s \n", outputFilepath)
	return nil
}

func SaveToZip() {

}

func SaveToTxt() {

}

func GetFileContent(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

func FileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}
