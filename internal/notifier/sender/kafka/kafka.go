package kafka

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/KseniiaSalmina/tikkichest-portfolio-service/internal/config"
	"log"
	"strconv"
	"sync"
)

type ProducerManager struct {
	producer      sarama.AsyncProducer
	topic         string
	finishClosing *sync.WaitGroup
}

func NewProducerManager(cfg config.Kafka) (*ProducerManager, error) {
	prod, err := sarama.NewAsyncProducer([]string{fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer manager: %w", err)
	}

	return &ProducerManager{
		producer: prod,
		topic:    cfg.Topic,
	}, nil
}

func (pm *ProducerManager) Run(ctx context.Context) {
	pm.finishClosing.Add(1)
	defer pm.finishClosing.Done()

	for {
		select {
		case err := <-pm.producer.Errors():
			log.Println(err) // TODO: logger
		case <-ctx.Done():
			return
		}
	}
}

func (pm *ProducerManager) SendNotification(id int, notification string) {
	message := sarama.ProducerMessage{
		Topic: pm.topic,
		Key:   sarama.StringEncoder(strconv.Itoa(id)),
		Value: sarama.StringEncoder(notification)}

	pm.producer.Input() <- &message
}

func (pm *ProducerManager) Shutdown() error {
	err := pm.producer.Close()
	pm.finishClosing.Wait()
	return err
}
