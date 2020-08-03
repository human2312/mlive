package util

// @Time : 2019年12月27日21:51:59
// @Author : Lemyhello
// @Desc: excel工具类 非通用

import (
	"github.com/tealeg/xlsx"
	"log"
)

//excelUtil excel工具类
type ExcelUtil struct{}

var file *xlsx.File
var sheet *xlsx.Sheet

// var row *xlsx.Row
// var cell *xlsx.Cell

func (u *ExcelUtil) Reset() {
	file = nil
	sheet = nil
	// row = nil
	// cell = nil
}

func (u *ExcelUtil) Init() {
	file = xlsx.NewFile()
	sheet, _ = file.AddSheet("Sheet1")
	// row = sheet.AddRow()
}

func (u *ExcelUtil) Save(outFile string) {
	err := file.Save(outFile)
	if err != nil {
		log.Printf(err.Error())
	}
}

//NewFile 新建excel
// func (u *ExcelUtil) NewFile(outFile string, outfit []string) {
// 	row := sheet.AddRow()
// 	for _, v := range outfit {
// 		cell = row.AddCell()
// 		cell.Value = v
// 	}
// }

//AppendFile 追加excel
func (u *ExcelUtil) AppendFile(outfit [9]string) {
	row := sheet.AddRow()
	for _, v := range outfit {
		cell := row.AddCell()
		cell.Value = v
	}

}
