package main

import (
	"fmt"

	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	pages := tview.NewPages()

	form := tview.NewForm()
	form.AddDropDown("称谓", []string{"先生", "女士", "博士", "老师", "师傅"}, 0, nil).
		AddInputField("姓名", "", 20, nil, nil).
		AddCheckbox("年龄 18+", false, nil).
		AddPasswordField("密码", "", 10, '*', nil).
		AddButton("保存", func() {
			_, title := form.GetFormItem(0).(*tview.DropDown).GetCurrentOption()
			userName := form.GetFormItem(1).(*tview.InputField).GetText()
			cw := form.GetFormItem(2).(*tview.Checkbox).GetLabel()

			alert(pages, "alert-dialog", fmt.Sprintf("保存成功，%s %s %s！", userName, title, cw))
		}).
		AddButton("退出", func() {
			app.Stop()
		})
	form.SetBorder(true).SetTitle("输入一些内容").SetTitleAlign(tview.AlignCenter)
	pages.AddPage("base", form, true, true)

	if err := app.SetRoot(pages, true).Run(); err != nil {
		panic(err)
	}
}

func alert(pages *tview.Pages, id string, message string) *tview.Pages {
	return pages.AddPage(
		id,
		tview.NewModal().
			SetText(message).
			AddButtons([]string{"确定"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				pages.HidePage(id).RemovePage(id)
			}),
		false,
		true,
	)
}
