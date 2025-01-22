package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

type Producer struct {
	Producer sarama.SyncProducer
}

func NewProducer(brokers []string, clientID string) (*Producer, error) {
	config := sarama.NewConfig()
	config.ClientID = clientID
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &Producer{Producer: producer}, nil
}

func (p *Producer) SendMessage(topic, message string) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := p.Producer.SendMessage(msg)
	if err != nil {
		return err
	}

	log.Printf("Message sent to topic %s [partition: %d, offset: %d]", topic, partition, offset)
	return nil
}

func (p *Producer) Close() error {
	return p.Producer.Close()
}
