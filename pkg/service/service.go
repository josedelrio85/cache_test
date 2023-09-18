package service

import "fmt"

type Actionable interface {
	DoSomething(input string) string
}

type actionShout struct {
	Desc string
}

func NewActionShout(input string) actionShout {
	return actionShout{
		Desc: input,
	}
}

func (a actionShout) DoSomething(input string) string {
	return fmt.Sprintf("%s %s", a.Desc, input)
}

type actionWhistle struct {
	Desc     string
	Whistels int
}

func NewActionWhistle(input string, number int) actionWhistle {
	return actionWhistle{
		Desc:     input,
		Whistels: number,
	}
}

func (a actionWhistle) DoSomething(input string) string {
	return fmt.Sprintf("%s - %d %s", a.Desc, a.Whistels, input)
}
