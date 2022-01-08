package combination

type WinForm struct {
	UIAttr

	Components []UIComponent
}

func (window *WinForm) PrintUIComponent() {
	client.printContainer(window, window)
}

fucn (window *WinForm) AddUIComponent(component UIComponent) {
	window.Components = append(window.Components, component)
}
