package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

// KafkaProducerInterface описывает интерфейс для продюсера, что позволяет легко подменять реализацию (например, для тестирования).
type KafkaProducerInterface interface {
	// PublishMessage отправляет сообщение в Kafka с поддержкой контекста.
	PublishMessage(ctx context.Context, message interface{}) error
	// Close закрывает соединение продюсера.
	Close() error
}

// kafkaProducer реализует KafkaProducerInterface.
type kafkaProducer struct {
	producer sarama.SyncProducer
	topic    string
	logger   *zap.Logger
}

// NewKafkaProducer создаёт нового синхронного продюсера Kafka с заданными настройками.
// Использование интерфейса улучшает тестирование и заменяемость компонента.
func NewKafkaProducer(brokers []string, topic string, logger *zap.Logger) (KafkaProducerInterface, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Producer.Timeout = 5 * time.Second

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	return &kafkaProducer{
		producer: producer,
		topic:    topic,
		logger:   logger,
	}, nil
}

// PublishMessage отправляет сообщение, предварительно сериализуя его в JSON.
// Добавлен параметр контекста (context.Context) для возможности отмены операции, если это потребуется.
func (kp *kafkaProducer) PublishMessage(ctx context.Context, message interface{}) error {
	// Проверяем, не отменён ли контекст до начала операции.
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if message == nil {
		return fmt.Errorf("message is nil")
	}

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: kp.topic,
		Value: sarama.ByteEncoder(data),
	}

	// Отправка сообщения.
	partition, offset, err := kp.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	kp.logger.Info("Сообщение отправлено в Kafka",
		zap.Int32("partition", partition),
		zap.Int64("offset", offset))
	return nil
}

// Close корректно закрывает продюсера, оборачивая возможные ошибки.
func (kp *kafkaProducer) Close() error {
	if err := kp.producer.Close(); err != nil {
		return fmt.Errorf("failed to close Kafka producer: %w", err)
	}
	return nil
}
