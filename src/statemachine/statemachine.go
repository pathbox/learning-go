package statemachine

type Handler func(interface{})(string, interface{})

type Machine struct {
  Handlers map[string]Handler
  StartState string
  EndStates map[string]bool
}

func (machine *Machine) AddState(handlerName string, handlerFn Handler){
  machine.Handlers[handlerName] = handlerFn
}

func (machine *Machine) AddEndState(endState string){
  machine.EndStates[endState] = true
}

func (machine *Machine) Execute(cargo interface{}){
  if handler, present := machine.Handlers[machine.StartState]; present{
    for {
      nextState, nextCargo := handler(cargo)
      _, finished := machine.EndStates[nextState]
      if finished{
        break
      }else {
        handler, present = machine.Handlers[nextState]
        cargo = nextCargo
      }
    }
  }
}
