package transactions

import (
	"time"

	"github.com/mataharibiz/ward"
	"github.com/mataharibiz/ward/rmq"
	"github.com/mataharibiz/ward/transactions/models"
)

func StateAction(data interface{}) {
	currentTime := time.Now().UTC()
	stateData := models.TransactionState{}

	event := rmq.EventData{
		EventType:   "state-actions",
		PublishDate: &currentTime,
		Data:        stateData,
	}

	event.Publish(ward.GetEnv("OBSERVER_EVENT", "dmp_observer"))
}
