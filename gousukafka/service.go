package gousukafka

import (
	"fmt"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/indece-official/go-gousu/gousu"
	"github.com/indece-official/go-gousu/gousu/logger"
	"github.com/namsral/flag"
)

// ServiceName is the name of the kafka service
const ServiceName = "kafka"

var (
	kafkaHost            = flag.String("kafka_host", "127.0.0.1", "Kafka broker host")
	kafkaPort            = flag.Int("kafka_port", 9092, "Kafka broker port")
	kafkaGroupID         = flag.String("kafka_group_id", "my-group", "Kafka group id")
	kafkaAutoOffsetReset = flag.String("kafka_auto_offset_reset", "earliest", "Kafka auto offset reset")
	kafkaAutoCommit      = flag.Bool("kafka_auto_commit", false, "Kafka auto commit")
)

type kafkaDoneEvent struct {
	Error   error
	Message *kafka.Message
}

// IService defines the interface of the kafka service
type IService interface {
	gousu.IService

	Subscribe(topic string) (chan *kafka.Message, error)
	Produce(topic string, value []byte) error
	Done(msg *kafka.Message, err error)
}

// Service provides a service for basic kafka client functionality
type Service struct {
	log            *logger.Log
	error          error
	running        bool
	producer       *kafka.Producer
	consumer       *kafka.Consumer
	topics         []string
	subscribers    map[string]chan *kafka.Message
	subscriberDone chan kafkaDoneEvent
}

// Verify that *Service implements IService
var _ IService = (*Service)(nil)

// Name returns the name of the kafka service from ServiceName
func (s *Service) Name() string {
	return ServiceName
}

// Start starts the KafkaService and connects the Kafka consumer
func (s *Service) Start() error {
	var err error

	config := &kafka.ConfigMap{
		"bootstrap.servers":     fmt.Sprintf("%s:%d", *kafkaHost, *kafkaPort),
		"group.id":              *kafkaGroupID,
		"broker.address.family": "v4",
		"session.timeout.ms":    6000,
		"enable.auto.commit":    *kafkaAutoCommit,
		"auto.offset.reset":     *kafkaAutoOffsetReset,
	}

	s.log.Infof("Connecting to kafka on %s:%d", *kafkaHost, *kafkaPort)

	s.producer, err = kafka.NewProducer(config)
	if err != nil {
		return s.log.ErrorfX("Can't connect to kafka on %s:%d as producer: %s", *kafkaHost, *kafkaPort, err)
	}

	s.consumer, err = kafka.NewConsumer(config)
	if err != nil {
		return s.log.ErrorfX("Can't connect to kafka on %s:%d as consumer: %s", *kafkaHost, *kafkaPort, err)
	}

	s.logConsumerBrokers()

	if len(s.topics) > 0 {
		s.consumer.SubscribeTopics(s.topics, nil)

		if !s.running {
			s.run()
		}
	}

	return nil
}

// Stop closes kafka consumer & producer
func (s *Service) Stop() error {
	err := s.consumer.Close()
	if err != nil {
		return s.log.ErrorfX("Can't close consumer: %s", err)
	}

	s.producer.Close()

	return nil
}

func (s *Service) logConsumerBrokers() {
	meta, err := s.consumer.GetMetadata(nil, false, 1000)
	if err != nil {
		s.log.Errorf("Can't fetch metadata for consumers: %s", err)
		return
	}

	brokers := []string{}

	for _, broker := range meta.Brokers {
		brokers = append(brokers, fmt.Sprintf("%d (%s:%d)", broker.ID, broker.Host, broker.Port))
	}

	s.log.Infof("Connected to consumer brokers: %s", strings.Join(brokers, ", "))
}

// Core loop for consumers
func (s *Service) run() {
	if s.running {
		return
	}

	s.running = true

	go func() {
		s.log.Infof("Started Kafka comsumer for topics %s", s.topics)

		s.error = nil

		for s.running {
			ev := s.consumer.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				s.error = nil

				topic := *e.TopicPartition.Topic
				subscriber, ok := s.subscribers[topic]
				if !ok {
					s.log.Errorf("Missing subscriber for topic '%s'", topic)
				}
				subscriber <- e

				doneEvent := <-s.subscriberDone

				s.consumer.CommitMessage(doneEvent.Message)
			case kafka.Error:
				// Errors should generally be considered
				// informational, the client will try to
				// automatically recover.
				// But in this example we choose to terminate
				// the application if all brokers are down.
				s.log.Errorf("Consumer error: %v (%v)", e.Code(), e)
				if e.Code() == kafka.ErrAllBrokersDown {
					s.running = false
					s.error = fmt.Errorf("all brokers down")
				}
			default:
				continue
			}
		}

		s.consumer.Close()
	}()
}

// Subscribe subscribes to a topic and returns a channel
//
// Important: After receiving a message on the returned channel,
//            Done(...) must be called, else the function will block
func (s *Service) Subscribe(topic string) (chan *kafka.Message, error) {
	if _, ok := s.subscribers[topic]; ok {
		return nil, fmt.Errorf("already subscribed to topic '%s' with group '%s'", topic, *kafkaGroupID)
	}

	s.topics = append(s.topics, topic)
	s.subscribers[topic] = make(chan *kafka.Message)

	s.log.Infof("Subscribed for kafka topic '%s' with group '%s'", topic, *kafkaGroupID)

	if s.consumer != nil {
		s.consumer.SubscribeTopics(s.topics, nil)

		if !s.running {
			s.run()
		}
	}

	return s.subscribers[topic], nil
}

// Produce emits a message to a topic
// The topic is automatically created (when enabled on the kafka server) if it doesn't exist
func (s *Service) Produce(topic string, value []byte) error {
	return s.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: value,
	}, nil)
}

// Done must be called after receiving a message via Subscribe(...)
func (s *Service) Done(msg *kafka.Message, err error) {
	s.subscriberDone <- kafkaDoneEvent{
		Error:   err,
		Message: msg,
	}
}

// Health provides information about if the Service is healthy
func (s *Service) Health() error {
	if s.error != nil {
		return fmt.Errorf("kafka service unhealthy: %s", s.error)
	}

	return nil
}

// NewService creates a new initialized instance of Service
func NewService(ctx gousu.IContext) gousu.IService {
	return &Service{
		subscribers:    make(map[string](chan *kafka.Message)),
		subscriberDone: make(chan kafkaDoneEvent),
		log:            logger.GetLogger("service.kafka"),
		running:        false,
	}
}

// Assert NewService matches gousu.ServiceFactory
var _ (gousu.ServiceFactory) = NewService
