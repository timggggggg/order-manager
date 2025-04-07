package outbox

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/IBM/sarama"
)

type ConsumerWorkerPool struct {
	consumer     sarama.Consumer
	topic        string
	workersCount int
	wg           sync.WaitGroup
}

func NewConsumerWorkerPool(workersCount int, brokers []string, topic string) (*ConsumerWorkerPool, error) {
	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating consumer: %v", err)
	}
	return &ConsumerWorkerPool{
		consumer:     consumer,
		topic:        topic,
		workersCount: workersCount,
	}, nil
}

func (c *ConsumerWorkerPool) Start(ctx context.Context) {
	for range c.workersCount {
		c.wg.Add(1)
		go c.run(ctx)
	}
}

func (c *ConsumerWorkerPool) Shutdown() {
	c.consumer.Close()
	c.wg.Wait()
}

func (c *ConsumerWorkerPool) run(ctx context.Context) {
	defer c.wg.Done()

	partitionConsumer, err := c.consumer.ConsumePartition(c.topic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}
	defer partitionConsumer.Close()

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			var audit any
			if err := json.Unmarshal(msg.Value, &audit); err != nil {
				fmt.Printf("Failed to unmarshal message: %v\n\n", err)
				continue
			}
			fmt.Printf("Consumed message: %v\n\n", audit)
		case err := <-partitionConsumer.Errors():
			fmt.Printf("Consumer error: %v\n\n", err)
		case <-ctx.Done():
			fmt.Print("Stopping consumer\n\n")
			return
		}
	}
}
