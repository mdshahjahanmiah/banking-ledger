package broker

import (
	"context"
	"encoding/json"
	"time"

	"github.com/mdshahjahanmiah/banking-ledger/model"
	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(brokerURL, topic string) Producer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{brokerURL},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
	return &KafkaProducer{writer: writer}
}

func (p *KafkaProducer) PublishTransaction(txn model.Transaction) error {
	msg, err := json.Marshal(txn)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(txn.ID),
		Value: msg,
		Time:  time.Now(),
	})
}
