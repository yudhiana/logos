package observer

import (
	"github.com/google/uuid"
)

func NewObserver(name string, f ObserverAction) Observer {
	return &ConcreteObserver{
		ID:     uuid.NewString(),
		Name:   name,
		Action: f,
	}
}

type ConcreteObserver struct {
	ID     string
	Name   string
	Action ObserverAction
}

func (c *ConcreteObserver) Notify(data interface{}) {
	c.Action(data)
}
