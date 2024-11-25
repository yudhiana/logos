package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/mataharibiz/ward"
	"github.com/mataharibiz/ward/logging"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler func(*sarama.ConsumerMessage, sarama.ConsumerGroup, sarama.ConsumerGroupSession) error

type WriteLog func(db *mongo.Database, collection string, args ...any) (id string, err error)

type RecordError struct {
	Active     bool
	DB         *mongo.Database
	Collection string
	WriteLog   WriteLog
}

type ConsumerGroup struct {
	cancel context.CancelFunc
	ctx    context.Context

	latestRetryCount int
	errorId          string
	RecordError      RecordError

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
	cg.cancel = cancel

	retryInterval := cg.RetryConfiguration.Interval
	maxRetries := cg.RetryConfiguration.MaxRetries
	backoffFactor := cg.RetryConfiguration.BackoffFactor

	go func() {
		<-WaitForSignal()
		cancel()
	}()

	var errorId string
	for attempt := 0; maxRetries < 0 || attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			logging.NewLogger().Info(fmt.Sprintf("attempt %d to connect to Kafka...", attempt))
		}

		// set up the Kafka consumer group
		consumerGroup, errConsumerGroup := GetConsumerGroup(cg)
		if errConsumerGroup != nil {
			logging.NewLogger().Error("failed to create kafka consumer group", "error", errConsumerGroup)

			if attempt == maxRetries && maxRetries != -1 {
				return
			}

			cg.writeLog(errorId, errConsumerGroup)
			errorId = cg.GetErrorId()

			// exponential backoff for retries
			sleepDuration := retryInterval * time.Duration(attempt) * time.Duration(backoffFactor)
			if errSleepCtx := ward.SleepWithContext(ctx, sleepDuration); errSleepCtx != nil {
				logging.NewLogger().Error("failed to sleep", "error", errSleepCtx)
				return
			}

			continue
		}
		cg.latestRetryCount = attempt

		logging.NewLogger().Info("connected to Kafka successfully")
		defer consumerGroup.Close()

		errorId = "" // reset error id
		errConsumer := cg.consumerMessage(ctx, consumerGroup, handler)
		if errConsumer != nil {
			cg.writeLog(errorId, errConsumer)
			errorId = cg.GetErrorId()
		}

		if ctx.Err() == context.Canceled {
			break
		}

		// reconnect in case of failure
		logging.NewLogger().Info("lost connection to Kafka... attempting to reconnecting Kafka...")
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
			logging.NewLogger().Info("shutting down by context", "operation", "kafka-consumer-disconnect")
			if err := client.Close(); err != nil {
				logging.NewLogger().Error("failed to close kafka consumer", "error", err)
			}
			return ctx.Err()
		default:
			logging.NewLogger().Info("Listening for kafka messages")
			err := client.Consume(ctx, cg.Topics, NewKafkaHandler(handler, client))
			if err != nil {
				logging.NewLogger().Error("kafka consumer error", "error", err)
				return err
			}
		}
	}
}

func (cg *ConsumerGroup) CanceledConsumer() {
	cg.cancel()
}

func (cg *ConsumerGroup) writeLog(errorId string, errMsg error) {
	re := cg.RecordError
	if re.Active {
		if errorId == "" {
			id, errLog := re.WriteLog(re.DB, re.Collection, errMsg)
			if errLog != nil {
				logging.NewLogger().Error("failed to write log", "error", errLog)
			}
			cg.errorId = id
		}
	}
}
func (cg *ConsumerGroup) GetLatestRetryCount() int {
	return cg.latestRetryCount
}

func (cg *ConsumerGroup) GetErrorId() string {
	return cg.errorId
}
