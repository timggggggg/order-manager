package kafka

import (
	"fmt"

	"github.com/IBM/sarama"
)

type KafkaProducer struct {
	producer sarama.SyncProducer
}

func NewKafkaProducer(brokers []string) (*KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("error creating producer: %v", err)
	}

	return &KafkaProducer{
		producer: producer,
	}, nil
}

func (kp *KafkaProducer) SendMessage(topic string, payload []byte) error {
	_, _, err := kp.producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(payload),
	})

	return err
}

func (kp *KafkaProducer) Close() error {
	err := kp.producer.Close()
	return err
}
