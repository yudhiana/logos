package actions

import (
	"log"
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
	log.Println("sangeEvent: ", sangeEvent)
	sangeEvent.PublishDefault()
}
