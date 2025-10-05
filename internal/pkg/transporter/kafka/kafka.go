package kafka

import (
	"fmt"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/rohanchauhan02/sequence-service/internal/config"
	"github.com/rohanchauhan02/sequence-service/internal/pkg/logger"
)

var log = logger.NewLogger("KAFKA")

type KafkaClient interface {
	Publish(topic string, message []byte) error
	Close() error
}

type kafkaClient struct {
	producer sarama.SyncProducer
	brokers  []string
}

func NewKafkaClient(conf config.ImmutableConfig) (KafkaClient, error) {
	kConf := conf.GetKafkaConf()
	brokers := strings.Split(kConf.Broker, ",")

	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Return.Successes = true
	cfg.Producer.Retry.Max = 5
	cfg.Producer.Timeout = 30 * time.Second
	cfg.Version = sarama.V2_8_0_0

	// Retry connecting to Kafka
	var producer sarama.SyncProducer
	var err error
	for i := 0; i < 5; i++ {
		producer, err = sarama.NewSyncProducer(brokers, cfg)
		if err == nil {
			break
		}
		log.Warnf("Kafka producer not ready, retrying in 5s... (%v)", err)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer after retries: %w", err)
	}

	log.Infof("Kafka producer connected to %v", brokers)

	return &kafkaClient{
		producer: producer,
		brokers:  brokers,
	}, nil
}

func (c *kafkaClient) Publish(topic string, message []byte) error {
	if c.producer == nil {
		log.Error("Kafka producer is not initialized")
		return fmt.Errorf("producer not initialized")
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	}

	partition, offset, err := c.producer.SendMessage(msg)
	if err != nil {
		log.Errorf("Failed to publish message to topic %s: %v", topic, err)
		return fmt.Errorf("failed to publish message to topic %s: %w", topic, err)
	}

	log.Infof("Published message to %s [partition=%d offset=%d]", topic, partition, offset)
	return nil
}

func (c *kafkaClient) Close() error {
	if c.producer == nil {
		return nil
	}

	if err := c.producer.Close(); err != nil {
		log.Errorf("Failed to close Kafka producer: %v", err)
		return fmt.Errorf("failed to close Kafka producer: %w", err)
	}

	log.Infof("Kafka producer disconnected from %v", c.brokers)
	return nil
}
