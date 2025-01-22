package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

type Consumer struct {
	Consumer sarama.Consumer
	Handlers map[string]MessageHandler
	Done     chan struct{}
}

type MessageHandler func(message *sarama.ConsumerMessage) error

func NewConsumer(brokers []string, clientID string) (*Consumer, error) {
	config := sarama.NewConfig()
	config.ClientID = clientID
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, err
	}

	c := &Consumer{Consumer: consumer, Handlers: make(map[string]MessageHandler), Done: make(chan struct{})}
	return c, nil
}

func (c *Consumer) Subscribe(topic string, handler MessageHandler) error {
	partitionConsumer, err := c.Consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case msg := <-partitionConsumer.Messages():
				if err := handler(msg); err != nil {
					log.Println(err)
				}
			case <-c.Done:
				return
			}
		}
	}()

	c.Handlers[topic] = handler
	return nil
}

func (c *Consumer) Close() error {
	return c.Consumer.Close()
}
