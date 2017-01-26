package main

import (
  "os"

  "github.com/olekukonko/tablewriter"
)

table, _ := tablewriter.NewCSV(os.Stdout, "test.csv", true)
table.SetRowLine(true)         // Enable row line

// Change table lines
table.SetCenterSeparator("*")
table.SetColumnSeparator("‡")
table.SetRowSeparator("-")

table.SetAlignment(tablewriter.ALIGN_LEFT)
table.Render()
