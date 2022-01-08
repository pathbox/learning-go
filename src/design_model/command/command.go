package command

import "fmt"

type Command interface {
	Execute()
}

type CreateCommand struct {
	reveiver Receiver 
}

func (command CreateCommand) Execute() {
	command.reveiver.Action()
}

func Receiver interface {
	Action()
}

type CreateReceiver struct{}

func (receiver CreateReceiver) Action() {
	fmt.Println("execute create")
}

type Invoker struct {
	command Command 
}

func (invoker Invoker) Call {
	invoker.command.Execute()
}