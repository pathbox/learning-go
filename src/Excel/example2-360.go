package main

import (
	"fmt"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func main() {
	f, err := excelize.OpenFile("/Users/pathbox/model.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, name := range f.GetSheetMap() {
		rows, _ := f.GetRows(name)
		for index, row := range rows {
			asix := fmt.Sprintf("A%d", index+1)
			_, url, _ := f.GetCellHyperLink(name, asix)
			row = append(row, url)
			fmt.Println(row)
		}
	}
}
