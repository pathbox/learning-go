package main

import (
	"os"

	table "github.com/jedib0t/go-pretty/table"
)

func main() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "First Name", "Last Name", "Salary"})
	t.AppendRows([]table.Row{
		{1, "Arya", "Stark", 3000},
		{20, "Jon", "Snow", 2000, "You know nothing, Jon Snow!"},
	})
	t.AppendRow([]interface{}{300, "Tyrion", "Lannister", 5000})
	t.AppendFooter(table.Row{"", "", "Total", 10000})

	// t.SetStyle(table.StyleLight)
	t.SetStyle(table.StyleColoredBright)
	// t.SetStyle(table.Style{
	// 	Name: "myNewStyle",
	// 	Box: table.BoxStyle{
	// 		BottomLeft:       "\\",
	// 		BottomRight:      "/",
	// 		BottomSeparator:  "v",
	// 		Left:             "[",
	// 		LeftSeparator:    "{",
	// 		MiddleHorizontal: "-",
	// 		MiddleSeparator:  "+",
	// 		MiddleVertical:   "|",
	// 		PaddingLeft:      "<",
	// 		PaddingRight:     ">",
	// 		Right:            "]",
	// 		RightSeparator:   "}",
	// 		TopLeft:          "(",
	// 		TopRight:         ")",
	// 		TopSeparator:     "^",
	// 		UnfinishedRow:    " ~~~",
	// 	},
	// 	Color: table.ColorOptions{

	// 		Footer:       text.Colors{text.BgCyan, text.FgBlack},
	// 		Header:       text.Colors{text.BgHiCyan, text.FgBlack},
	// 		Row:          text.Colors{text.BgHiWhite, text.FgBlack},
	// 		RowAlternate: text.Colors{text.BgWhite, text.FgBlack},
	// 	},
	// 	Format: table.FormatOptions{
	// 		Footer: text.FormatUpper,
	// 		Header: text.FormatUpper,
	// 		Row:    text.FormatDefault,
	// 	},
	// 	Options: table.Options{
	// 		DrawBorder:      true,
	// 		SeparateColumns: true,
	// 		SeparateFooter:  true,
	// 		SeparateHeader:  true,
	// 		SeparateRows:    false,
	// 	},
	// })
	t.Render()
	t.RenderHTML()
	t.RenderMarkdown()
}
