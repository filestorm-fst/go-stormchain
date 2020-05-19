package ExcelHelper

import (
	"github.com/tealeg/xlsx"
)

type Excel struct {
	File *xlsx.File
	Sheet *xlsx.Sheet

}

//1 NewExcel
//2 Write
//3 Save

func NewExcel() (*Excel, error) {
	excel := new(Excel)

	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		return nil, err
	}
	excel.File = file
	excel.Sheet = sheet

	return excel, nil
}

//一次写一行
func (e *Excel) Write(values []interface{}) {
	row := e.Sheet.AddRow()
	for _, cellValue := range values {
		cell := row.AddCell()
		cell.SetValue(cellValue)
	}
}

//保存
func (e *Excel) Save(path string) (error) {
	err := e.File.Save(path)
	if err != nil {
		return err
	}
	return nil
}
