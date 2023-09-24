package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/leonid-tankov/OTUS_go_final_project/internal/config"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	Writer *kafka.Writer
	Broker string
}

func NewProducer(cfg *config.Config) *Producer {
	broker := fmt.Sprintf("%s:%s", cfg.Kafka.Host, cfg.Kafka.Port)
	return &Producer{
		Writer: &kafka.Writer{
			Addr:  kafka.TCP(broker),
			Topic: "banner-rotation-analytics",
		},
		Broker: broker,
	}
}

func (p *Producer) Produce(ctx context.Context, msg Message) error {
	value, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return p.Writer.WriteMessages(ctx, kafka.Message{
		Value: value,
	})
}

func (p *Producer) CreateTopic(topic string) error {
	var conn *kafka.Conn
	conn, err := kafka.Dial("tcp", p.Broker)
	if err != nil {
		conn, err = reconnect(p.Broker, 1)
		if err != nil {
			return err
		}
	}
	defer conn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}
	err = conn.CreateTopics(topicConfigs...)
	if err != nil {
		return err
	}
	return nil
}

func reconnect(broker string, retryCount int) (*kafka.Conn, error) {
	time.Sleep(time.Second * 3)
	conn, err := kafka.Dial("tcp", broker)
	if err != nil {
		if retryCount < 5 {
			return reconnect(broker, retryCount+1)
		}
		return conn, err
	}
	return conn, nil
}
