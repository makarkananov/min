package kafka

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"min/internal/core/domain"

	"github.com/IBM/sarama"
)

type EventProducer struct {
	producer sarama.SyncProducer
	topic    string
}

func NewEventProducer(brokers []string, topic string) (*EventProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	return &EventProducer{
		producer: producer,
		topic:    topic,
	}, nil
}

func (ep *EventProducer) Produce(event *domain.Event) error {
	log.Infof("Producing event to Kafka: %v", event)
	eventBytes, err := json.Marshal(event)
	if err != nil {
		log.Errorf("Failed to marshal event: %v", err)
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: ep.topic,
		Value: sarama.ByteEncoder(eventBytes),
	}

	partition, offset, err := ep.producer.SendMessage(msg)
	if err != nil {
		log.Errorf("Failed to send message: %v", err)
		return fmt.Errorf("failed to send message: %w", err)
	}

	log.Infof("Message is stored in topic(%s)/partition(%d)/offset(%d)", ep.topic, partition, offset)
	return nil
}
