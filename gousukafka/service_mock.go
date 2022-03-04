package gousukafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/indece-official/go-gousu/gousu"
)

// MockService for simply mocking IService
type MockService struct {
	gousu.MockService

	DoneFunc            func(msg *kafka.Message, err error)
	SubscribeFunc       func(topic string) (chan *kafka.Message, error)
	ProduceFunc         func(topic string, value []byte) error
	DoneFuncCalled      int
	SubscribeFuncCalled int
	ProduceFuncCalled   int
}

// MockService implements IService
var _ (IService) = (*MockService)(nil)

// Done calls DoneFunc and increases DoneFuncCalled
func (s *MockService) Done(msg *kafka.Message, err error) {
	s.DoneFuncCalled++

	s.DoneFunc(msg, err)
}

// Subscribe calls SubscribeFunc and increases SubscribeFuncCalled
func (s *MockService) Subscribe(topic string) (chan *kafka.Message, error) {
	s.SubscribeFuncCalled++

	return s.SubscribeFunc(topic)
}

// Produce calls ProduceFunc and increases ProduceFuncCalled
func (s *MockService) Produce(topic string, value []byte) error {
	s.ProduceFuncCalled++

	return s.ProduceFunc(topic, value)
}

// NewMockService creates a new initialized instance of MockService
func NewMockService() *MockService {
	return &MockService{
		MockService: gousu.MockService{
			NameFunc: func() string {
				return ServiceName
			},
		},

		DoneFunc: func(msg *kafka.Message, err error) {},
		SubscribeFunc: func(topic string) (chan *kafka.Message, error) {
			return nil, nil
		},
		ProduceFunc: func(topic string, value []byte) error {
			return nil
		},
		DoneFuncCalled:      0,
		SubscribeFuncCalled: 0,
		ProduceFuncCalled:   0,
	}
}
