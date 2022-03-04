package gousustomp

import (
	"flag"
	"fmt"

	"github.com/go-stomp/stomp/v3"
	"github.com/indece-official/go-gousu/gousu"
	"github.com/indece-official/go-gousu/gousu/logger"
)

const ServiceName = "stomp"

var (
	stompHost     = flag.String("stomp_host", "localhost", "")
	stompPort     = flag.Int("stomp_port", 61613, "")
	stompUsername = flag.String("stomp_username", "", "")
	stompPassword = flag.String("stomp_password", "", "")
)

type IService interface {
	gousu.IService

	SubscribeToQueue(queue string) (chan *stomp.Message, *stomp.Subscription, error)
	Ack(msg *stomp.Message) error
	PublishToQueue(queue string, contentType string, body []byte) error
}

type Service struct {
	log  *logger.Log
	conn *stomp.Conn
}

var _ (IService) = (*Service)(nil)

func (s *Service) Name() string {
	return ServiceName
}

func (s *Service) Start() error {
	var err error

	host := fmt.Sprintf("%s:%d", *stompHost, *stompPort)

	options := []func(*stomp.Conn) error{}

	if *stompUsername != "" && *stompPassword != "" {
		options = append(options, stomp.ConnOpt.Login(*stompUsername, *stompPassword))
	}

	s.conn, err = stomp.Dial("tcp", host, options...)
	if err != nil {
		return fmt.Errorf("can't connect to stomp server on %s: %s", host, err)
	}

	return nil
}

func (s *Service) Health() error {
	return nil
}

func (s *Service) Stop() error {
	return s.conn.Disconnect()
}

func (s *Service) SubscribeToQueue(queue string) (chan *stomp.Message, *stomp.Subscription, error) {
	sub, err := s.conn.Subscribe(queue, stomp.AckClientIndividual)
	if err != nil {
		return nil, nil, fmt.Errorf("error subscribing to queue '%s': %s", queue, err)
	}

	return sub.C, sub, nil
}

func (s *Service) Ack(msg *stomp.Message) error {
	err := s.conn.Ack(msg)
	if err != nil {
		return fmt.Errorf("can't acknowledge message: %s", err)
	}

	return nil
}

func (s *Service) PublishToQueue(queue string, contentType string, body []byte) error {
	err := s.conn.Send(queue, contentType, body)
	if err != nil {
		return fmt.Errorf("error sending to queue '%s': %s", queue, err)
	}

	return nil
}

func NewService(ctx gousu.IContext) gousu.IService {
	return &Service{
		log: logger.GetLogger(fmt.Sprintf("service.%s", ServiceName)),
	}
}

var _ (gousu.ServiceFactory) = NewService
