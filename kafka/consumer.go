package kafka

import (
	"context"
	"fmt"
	"time"

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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	retryInterval := 5 * time.Second
	maxRetries := -1
	backoffFactor := 2.0

	for attempt := 1; maxRetries == -1 || attempt <= maxRetries; attempt++ {
		logging.NewLogger().Info(fmt.Sprintf("Attempt %d to connect to Kafka...", attempt))

		// set up the Kafka consumer group
		client, errConsumer := GetConsumerGroup(cg)
		if errConsumer != nil {
			logging.NewLogger().Error("failed to create kafka consumer", "error", errConsumer)

			if maxRetries != -1 && attempt == maxRetries {
				return
			}

			// exponential backoff
			sleepDuration := retryInterval * time.Duration(attempt) * time.Duration(backoffFactor)
			time.Sleep(sleepDuration)
			continue
		}

		logging.NewLogger().Info("Connected to Kafka successfully")
		defer client.Close()

		go cg.consumerMessage(ctx, client, handler)

		go cg.handleGracefulShutdown(ctx, client)

		// reconnect in case of failure
		logging.NewLogger().Info("Reconnecting to Kafka...")
		time.Sleep(retryInterval)
		return
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

func (cg *ConsumerGroup) consumerMessage(ctx context.Context, client sarama.ConsumerGroup, handler Handler) {
	defer Recover()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			logging.NewLogger().Info("Listening for kafka messages")
			err := client.Consume(ctx, cg.Topics, NewKafkaHandler(handler, client))
			if err != nil {
				logging.NewLogger().Error("kafka consumer error", "error", err)
			}
		}
	}
}

func (cg *ConsumerGroup) handleGracefulShutdown(ctx context.Context, client sarama.ConsumerGroup) {
	select {
	case <-WaitForSignal():
		logging.NewLogger().Warn("shutting down", "operation", "kafka-consumer-disconnect")
		if err := client.Close(); err != nil {
			logging.NewLogger().Error("failed to close kafka consumer", "error", err)
		}

	case <-ctx.Done():
		logging.NewLogger().Info("shutting down", "operation", "kafka-consumer-disconnect")
		if err := client.Close(); err != nil {
			logging.NewLogger().Error("failed to close kafka consumer", "error", err)
		}
	}
}
