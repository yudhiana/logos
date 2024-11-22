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
	RetryConfiguration  RetryConfiguration
}

type RetryConfiguration struct {
	Interval      time.Duration
	MaxRetries    int
	BackoffFactor float64
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

	retryInterval := cg.RetryConfiguration.Interval
	maxRetries := cg.RetryConfiguration.MaxRetries
	backoffFactor := cg.RetryConfiguration.BackoffFactor

	for attempt := 0; maxRetries < 0 || attempt <= maxRetries; attempt++ {
		logging.NewLogger().Info("config", "maxRetries", maxRetries, "backoffFactor", backoffFactor, "retryInterval", retryInterval, "attempt", attempt)
		if attempt > 0 {
			logging.NewLogger().Info(fmt.Sprintf("attempt %d to connect to Kafka...", attempt))
		}

		// set up the Kafka consumer group
		client, errConsumer := GetConsumerGroup(cg)
		if errConsumer != nil {
			logging.NewLogger().Error("failed to create kafka consumer", "error", errConsumer)

			if attempt == maxRetries && maxRetries > 0 {
				logging.NewLogger().Warn("reached max retries", "maxRetries", maxRetries)
				break
			}

			// exponential backoff
			sleepDuration := retryInterval * time.Duration(attempt) * time.Duration(backoffFactor)
			time.Sleep(sleepDuration)
			continue
		}

		logging.NewLogger().Info("connected to Kafka successfully")
		defer client.Close()

		go cg.handleGracefulShutdown(ctx, client)

		if err := cg.consumerMessage(ctx, client, handler); err != nil {
			logging.NewLogger().Error("kafka consumer error", "error", err)
		}

		// reconnect in case of failure
		logging.NewLogger().Info("reconnecting to Kafka...")
		time.Sleep(retryInterval)
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

func (cg *ConsumerGroup) consumerMessage(ctx context.Context, client sarama.ConsumerGroup, handler Handler) error {
	defer Recover()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			logging.NewLogger().Info("Listening for kafka messages")
			err := client.Consume(ctx, cg.Topics, NewKafkaHandler(handler, client))
			if err != nil {
				return err
			}
		}
	}
}

func (cg *ConsumerGroup) handleGracefulShutdown(ctx context.Context, client sarama.ConsumerGroup) {
	select {
	case <-WaitForSignal():
		logging.NewLogger().Warn("shutting down by signals", "operation", "kafka-consumer-disconnect")
		if err := client.Close(); err != nil {
			logging.NewLogger().Error("failed to close kafka consumer", "error", err)
		}

	case <-ctx.Done():
		logging.NewLogger().Info("shutting down by context", "operation", "kafka-consumer-disconnect")
		if err := client.Close(); err != nil {
			logging.NewLogger().Error("failed to close kafka consumer", "error", err)
		}
	}
}
