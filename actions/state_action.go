package actions

import (
	"time"

	"github.com/mataharibiz/sange"
)

func StateAction(data interface{}) {
	currentTime := time.Now().UTC()
	payload := map[string]interface{}{
		"time":  currentTime,
		"state": "initial",
		"data":  data,
	}

	sangeEvent := sange.EventData{
		EventType:   "state-actions",
		PublishDate: &currentTime,
		Data:        payload,
	}

	sangeEvent.Publish(sange.GetEnv("OBSERVER_EVENT", "dmp_observer"))
}
