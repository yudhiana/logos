package rmq

import (
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// EventData struct for event data
type EventData struct {
	EventType   string      `json:"event_type,omitempty"`
	Data        interface{} `json:"data,omitempty"`
	PublishDate *time.Time  `json:"date,omitempty"`
}

// FlatData return flat format of this event data
func (data *EventData) FlatData() []byte {
	if data.PublishDate == nil {
		now := time.Now().UTC()
		data.PublishDate = &now
	}

	result, _ := json.Marshal(data)
	return result
}

// Publish event to exchangeName message broker. This function using internal go-routine.
func (data *EventData) Publish(exchangeName string) {
	go func() {
		defer func() {
			if v := recover(); v != nil {
				log.Println("Publisher got panic!! Please don't worry :) ", v)
			}
		}()

		connURL := GetRabbitURL()
		conn, err := amqp.Dial(connURL)
		if err != nil {
			log.Println("Failed to connect to Rabbit MQ : ", err)
			return
		}
		defer conn.Close()

		// create channel
		ch, err := conn.Channel()
		if err != nil {
			log.Println("Failed to create Rabbit MQ channel : ", err)
			return
		}
		defer ch.Close()

		// declare fanout exchange
		err = ch.ExchangeDeclare(
			exchangeName, // name
			"fanout",     //type
			true,         //durable
			false,        //auto-deleted
			false,        // internal
			false,        //no-wait
			nil,          //arguments
		)
		if err != nil {
			log.Println("Failed to declare an exchange : ", err)
			return
		}

		// publish event
		err = ch.Publish(
			exchangeName, //exchange
			"",           // routing key
			false,        //mandatory
			false,        // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        data.FlatData(),
			},
		)

		if err != nil {
			log.Println("Failed to publish a message : ", err)
			return
		}

	}()
}
