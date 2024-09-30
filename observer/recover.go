package observer

import (
	"fmt"

	"github.com/mataharibiz/sange"
)

func Panic(data interface{}) {
	panicEvent := sange.EventData{
		EventType: "panic-observer",
		Data:      fmt.Sprintf("panic: %v", data),
	}
	panicEvent.Publish(sange.GetEnv("OBSERVER_EVENT", "dmp_observer"))
}
