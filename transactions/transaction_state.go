package transactions

import (
	"time"

	"github.com/mataharibiz/sange"
	"github.com/yudhiana99/ward/transactions/models"
)

func StateAction(data interface{}) {
	currentTime := time.Now().UTC()
	stateData := models.TransactionState{}

	sangeEvent := sange.EventData{
		EventType:   "state-actions",
		PublishDate: &currentTime,
		Data:        stateData,
	}

	sangeEvent.Publish(sange.GetEnv("OBSERVER_EVENT", "dmp_observer"))
}
