package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	log "github.com/sirupsen/logrus"
	"min/internal/core/domain"
	"min/internal/core/port"
)

type Consumer struct {
	consumerGroup sarama.ConsumerGroup
	service       port.StatisticsService
}

func NewKafkaConsumer(brokers []string, groupID string, service port.StatisticsService) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer group: %w", err)
	}

	return &Consumer{
		consumerGroup: consumerGroup,
		service:       service,
	}, nil
}

func (kc *Consumer) Start(ctx context.Context, topics []string) error {
	log.Infof("starting Kafka consumer for topics: %v", topics)
	handler := &consumerGroupHandler{service: kc.service}

	for {
		if err := kc.consumerGroup.Consume(ctx, topics, handler); err != nil {
			return fmt.Errorf("failed to consume messages: %w", err)
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

type consumerGroupHandler struct {
	service port.StatisticsService
}

func (h *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (h *consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Infof("received kafka message: %s", string(msg.Value))
		var event domain.Event
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Errorf("failed to unmarshal message: %v", err)
			continue
		}

		if err := h.service.AddEvent(sess.Context(), event); err != nil {
			log.Errorf("failed to add event: %v", err)
			continue
		}

		sess.MarkMessage(msg, "")
	}
	return nil
}
