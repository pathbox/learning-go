package combination

type UIComponent interface {
	PrintUIComponent()
	GetUIControlName()
}

type UIComponentAddtion interface {
	AddUIComponent(component UIComponent)
}

type UIAttr struct {
	Name string
}

var client = &PrintClient{}


