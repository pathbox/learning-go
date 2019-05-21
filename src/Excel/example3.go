package main

import (
	"fmt"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func main() {
	xlsxFile := "/Users/pathbox/model.xlsx"
	m, _ := ReadCustomerFromExcel(xlsxFile)
	fmt.Printf("%+v", m)
}

func ReadCustomerFromExcel(xlsxFile string) (map[string][]map[string]string, error) {
	m := map[string][]map[string]string{}
	xFile, err := excelize.OpenFile(xlsxFile)
	if err != nil {
		fmt.Println(err)
		return m, err
	}

	for _, name := range xFile.GetSheetMap() {
		rows, _ := xFile.GetRows(name)
		cuSlice := make([]map[string]string, 0)
		for index, row := range rows {
			if index == 0 {
				continue
			}
			asix := fmt.Sprintf("A%d", index+1)
			_, url, _ := xFile.GetCellHyperLink(name, asix)
			row = append(row, url)
			// fmt.Println(row)

			cu := map[string]string{}
			cu["Name"] = row[0]
			cu["University"] = row[1]
			cu["CraduateYear"] = row[2]
			cu["Company"] = row[3]
			cu["Title"] = row[4]
			cu["ExperienceYear"] = row[5]
			cu["Email"] = row[6]
			cu["Phone"] = row[7]
			cu["Wechat"] = row[8]
			cu["Address"] = row[9]
			cu["HomeURL"] = row[10]
			cuSlice = append(cuSlice, cu)
		}
		m[name] = cuSlice
	}
	return m, nil
}
