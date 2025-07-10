package kafka

import (
	"github.com/IBM/sarama"
)

func NewKafkaHandler(handler Handler, client sarama.ConsumerGroup) *KafkaHandler {
	return &KafkaHandler{
		handler: handler,
		client:  client,
	}
}

type KafkaHandler struct {
	handler Handler
	client  sarama.ConsumerGroup
}

func (h *KafkaHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *KafkaHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *KafkaHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		h.handler(msg, h.client, session)
	}
	return nil
}
