package ward

import (
	"log"

	"github.com/google/uuid"
)

func NewObserver() *ConcreteObserver {
	return &ConcreteObserver{
		ID: uuid.NewString(),
	}
}

type ConcreteObserver struct {
	ID string
}

func (c *ConcreteObserver) Update(data interface{}) {
	log.Printf("Observer ID %s received data: %v", c.ID, data)
}
