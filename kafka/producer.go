package kafka

import (
	"encoding/json"
	"sync"

	"github.com/IBM/sarama"
	"github.com/mataharibiz/ward/logging"
)

var kafkaProducerOnce sync.Once
var kafkaShutdownOnce sync.Once
var kafkaAsyncProducer sarama.AsyncProducer

type ProducerGroup struct {
	Idempotent          bool
	ManualConfiguration *sarama.Config
	Hosts               []string
	AssignmentType      KafkaProducerPartitionType
}

// NewKafkaProducer initializes and returns a new Kafka ProducerGroup.
// It sets up the asynchronous Kafka producer and starts goroutines to handle
// message success/error events and listen for interrupt signals for graceful shutdown.
func NewKafkaProducer(pg *ProducerGroup) *ProducerGroup {
	// Set up the asynchronous Kafka producer
	pg.setUpAsyncProducer()

	if kafkaAsyncProducer == nil {
		return pg
	}

	// Start a goroutine to handle success and error events
	go pg.startSuccessErrorHandler()

	// Start listening for interrupt signals to shut down gracefully
	go pg.listenForInterrupts()

	return pg
}

// PublishMessages encodes the provided message and sends it to Kafka asynchronously.
// The message is encoded using the encodeMessage method and then passed to produceMessages
// to be published to the specified Kafka topic.
func (pg *ProducerGroup) PublishMessages(topic string, key, message any) {
	if message != nil && kafkaAsyncProducer != nil {
		value := pg.encodeMessage(message)
		key := pg.encodeMessage(topic)
		logging.NewLogger().Info("message", "key", key, "value", value)
		go pg.produceMessages(topic, key, value)
	}
}

// Helper method to encode the message into sarama.Encoder
func (pg *ProducerGroup) encodeMessage(message any) sarama.Encoder {
	if message == nil {
		return nil
	}

	switch msg := message.(type) {
	case string:
		return sarama.StringEncoder(msg)
	case []byte:
		return sarama.ByteEncoder(msg)
	case []rune:
		return sarama.StringEncoder(string(msg))
	case map[string]string, map[string]any, map[any]any:
		byteMsg, _ := json.Marshal(msg)
		return sarama.ByteEncoder(byteMsg)
	default:
		byteMsg, _ := json.Marshal(msg)
		return sarama.ByteEncoder(byteMsg)
	}
}

// setUpAsyncProducer sets up the Kafka asynchronous producer for the ProducerGroup.
func (pg *ProducerGroup) setUpAsyncProducer() {
	kafkaProducerOnce.Do(func() {
		producerConfig := GetKafkaConfig(Producer)
		switch pg.AssignmentType {
		case ProducerAssignmentRandomPartition:
			producerConfig.Producer.Partitioner = sarama.NewRandomPartitioner
		case ProducerAssignmentRoundRobinPartition:
			producerConfig.Producer.Partitioner = sarama.NewRoundRobinPartitioner
		case ProducerAssignmentHashPartition:
			producerConfig.Producer.Partitioner = sarama.NewHashPartitioner
		case ProducerAssignmentManualPartition:
			producerConfig.Producer.Partitioner = sarama.NewManualPartitioner
		}

		if pg.Idempotent {
			producerConfig.Producer.Idempotent = true
		}

		if pg.ManualConfiguration != nil {
			producerConfig = pg.ManualConfiguration
		}

		var err error
		kafkaAsyncProducer, err = sarama.NewAsyncProducer(pg.Hosts, producerConfig)
		if err != nil {
			logging.NewLogger().Error("failed to create kafka producer", "error", err)
			return
		}
	})
}

// startSuccessErrorHandler listens to success and error channels continuously
func (pg *ProducerGroup) startSuccessErrorHandler() {
	for {
		select {
		case msg, ok := <-kafkaAsyncProducer.Successes():
			if !ok {
				return
			}
			logging.NewLogger().Info("Message sent successfully", "partition", msg.Partition, "offset", msg.Offset)

		case err, ok := <-kafkaAsyncProducer.Errors():
			if !ok {
				return
			}
			logging.NewLogger().Error("Failed to send message", "error", err)
		}
	}
}

// produceMessages sends the provided encoded message to the first Kafka topic in the ProducerGroup.
// The message is sent asynchronously, and its success or failure is handled by the startSuccessErrorHandler.
func (pg *ProducerGroup) produceMessages(topic string, key, value sarama.Encoder) {
	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: value,
		Key:   key,
	}

	// Send message to Kafka; success/error will be handled in startSuccessErrorHandler
	kafkaAsyncProducer.Input() <- message
}

// listenForInterrupts sets up signal handling for graceful shutdown
func (pg *ProducerGroup) listenForInterrupts() {
	sigterm := WaitForSignal()

	go func() {
		<-sigterm
		pg.CloseProducerGracefully()
	}()
}

// CloseProducerGracefully handles graceful closure of the Kafka producer
func (pg *ProducerGroup) CloseProducerGracefully() {
	kafkaShutdownOnce.Do(func() { // Ensure graceful shutdown occurs only once
		logging.NewLogger().Warn("Shutting down Kafka producer")

		// Close the async producer to release resources
		if kafkaAsyncProducer != nil {
			if err := kafkaAsyncProducer.Close(); err != nil {
				logging.NewLogger().Error("Error closing Kafka producer", "error", err)
			} else {
				logging.NewLogger().Warn("Kafka producer closed gracefully")
			}
		}
	})
}
