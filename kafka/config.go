package kafka

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/mataharibiz/ward"
	"github.com/mataharibiz/ward/logging"

	"github.com/IBM/sarama"
)

type KafkaType string

const (
	Producer KafkaType = "producer"
	Consumer KafkaType = "consumer"
)

func GetConsumerConfig() *ConsumerGroup {
	return &ConsumerGroup{
		Hosts: strings.Split(os.Getenv("KAFKA_HOST"), ","),
	}
}

func GetProducerConfig() *ProducerGroup {
	return &ProducerGroup{
		Hosts: strings.Split(os.Getenv("KAFKA_HOST"), ","),
	}
}

func GetKafkaConfig(useAs KafkaType) *sarama.Config {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Version = sarama.V3_6_0_0

	switch useAs {
	case Producer:
		saramaConfig.Producer.Return.Successes = true
		saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
		saramaConfig.Producer.Retry.Max = 10
		saramaConfig.Producer.Partitioner = sarama.NewRoundRobinPartitioner

	case Consumer:
		saramaConfig.Consumer.Return.Errors = true
		saramaConfig.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
		saramaConfig.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	}

	return saramaConfig
}

func WaitForSignal() <-chan os.Signal {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	return done
}

func Recover() {
	if err := recover(); err != nil {
		logging.NewLogger().Error("recovered from panic", "error", err)
		fmt.Println(ward.GetStackTrace())
	}
}
