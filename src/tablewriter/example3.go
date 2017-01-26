package main

import (
	"github.com/olekukonko/tablewriter"
	"os"
)

func main() {
	table, _ := tablewriter.NewCSV(os.Stdout, "test_info.csv", true)
	table.SetAlignment(tablewriter.ALIGN_LEFT) // Set Alignment
	table.Render()
}
