package kafka

import "fmt"

type Kafka struct {
	Producer *Producer
	Consumer *Consumer
}

func NewKafka(brokers []string, clientID string) (*Kafka, error) {
	producer, err := NewProducer(brokers, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %v", err)
	}

	consumer, err := NewConsumer(brokers, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %v", err)
	}

	return &Kafka{Producer: producer, Consumer: consumer}, nil
}

func (k *Kafka) Close() error {
	if err := k.Producer.Close(); err != nil {
		return fmt.Errorf("failed to close producer: %v", err)
	}

	if err := k.Consumer.Close(); err != nil {
		return fmt.Errorf("failed to close consumer: %v", err)
	}

	return nil
}
