package kafka

import (
	"context"
	"errors"

	"github.com/IBM/sarama"
	"github.com/mataharibiz/ward/logging"
)

type Handler func(*sarama.ConsumerMessage, sarama.ConsumerGroup, sarama.ConsumerGroupSession) error

type ConsumerGroup struct {
	ManualConfiguration *sarama.Config
	AssignmentType      KafkaConsumerAssignmentType
	AutoCommit          bool
	GroupID             string
	Topics              []string
	Hosts               []string
}

// NewKafkaConsumerGroup sets up a new Kafka consumer group with the given groupID and hosts,
// and starts consuming from the specified topics using the provided handler.
// The handler is called for each message received from the topics.
// The function will only return when either the consumer group is closed or the context is canceled.
// Closing the consumer group or canceling the context will cause the function to return.
// The function is designed to be run in a goroutine, so it will not block the calling goroutine.
func NewKafkaConsumerGroup(cg *ConsumerGroup, handler Handler) {

	// Set up the Kafka consumer group
	client, errConsumer := GetConsumerGroup(cg)
	if errConsumer != nil {
		logging.NewLogger().Error("kafka consumer error", "error", errConsumer)
		return
	}
	defer client.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			defer Recover()

			if err := client.Consume(ctx, cg.Topics, NewKafkaHandler(handler, client)); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}

				logging.NewLogger().Error("kafka consumer error", "error", err)
				return
			}

			if ctx.Err() != nil {
				return
			}
		}
	}()

	active := true
	for active {
		logging.NewLogger().Info("Listening for kafka messages")
		select {
		case <-WaitForSignal():
			logging.NewLogger().Warn("shutting down", "operation", "kafka-consumer-disconnect")
			if err := client.Close(); err != nil {
				logging.NewLogger().Error("failed to close kafka consumer", "error", err)
			}
			active = false
		case <-ctx.Done():
			logging.NewLogger().Info("shutting down", "operation", "kafka-consumer-disconnect")
			if err := client.Close(); err != nil {
				logging.NewLogger().Error("failed to close kafka consumer", "error", err)
			}
			active = false
		}

	}
}

// GetConsumerGroup creates a new Kafka consumer group with the given groupID and hosts.
// It sets up the consumer group with the config returned by GetKafkaConfig.
// The returned consumer group is ready to start consuming from Kafka.
func GetConsumerGroup(cg *ConsumerGroup) (sarama.ConsumerGroup, error) {
	consumerConfig := GetKafkaConfig(Consumer)
	switch cg.AssignmentType {
	case ConsumerGroupAssignmentStrategyRoundRobin:
		consumerConfig.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
		consumerConfig.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	case ConsumerGroupAssignmentStrategySticky:
		consumerConfig.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategySticky()}
		consumerConfig.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategySticky()
	case ConsumerGroupAssignmentStrategyRange:
		consumerConfig.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}
		consumerConfig.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRange()
	}

	if cg.AutoCommit {
		consumerConfig.Consumer.Offsets.AutoCommit.Enable = true
	}

	if cg.ManualConfiguration != nil {
		consumerConfig = cg.ManualConfiguration
	}

	consumer, errConsumer := sarama.NewConsumerGroup(cg.Hosts, cg.GroupID, consumerConfig)
	if errConsumer != nil {
		return nil, errConsumer
	}
	return consumer, nil
}
