package goususmtp

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"github.com/go-mail/mail"
	"github.com/indece-official/go-gousu/v2/gousu"
	"github.com/indece-official/go-gousu/v2/gousu/broadcaster"
	"github.com/indece-official/go-gousu/v2/gousu/logger"
	"github.com/namsral/flag"
)

// ServiceName defines the name of smtp service used for dependency injection
const ServiceName = "smtp"

var (
	smtpHost     = flag.String("smtp_host", "127.0.0.1", "")
	smtpPort     = flag.Int("smtp_port", 587, "")
	smtpUser     = flag.String("smtp_user", "", "")
	smtpPassword = flag.String("smtp_password", "", "")
	smtpFrom     = flag.String("smtp_from", "", "")
)

// EmailAttachement defines the base model of an email attachemet
type EmailAttachement struct {
	Filename string
	Mimetype string
	Embedded bool
	Content  []byte
}

// Email defines the base model of an email
type Email struct {
	// From is the name of the sender
	// If empty the config flag 'smtp_from' is used
	From         string
	To           string
	Subject      string
	BodyPlain    string
	BodyHTML     string
	Attachements []EmailAttachement
}

// IService defines the interface of the smtp service
type IService interface {
	gousu.IService

	SendEmail(m *Email) error
}

// Service provides an smtp sender running in a separate thread
type Service struct {
	log             *logger.Log
	dialer          *mail.Dialer
	closer          *mail.SendCloser
	stopBroadcaster *broadcaster.Bool
	runningFuncs    sync.WaitGroup
	error           error
	lastSend        *time.Time
	mutexCloser     sync.Mutex
}

var _ IService = (*Service)(nil)

// Name returns the name of the smtp service from ServiceName
func (s *Service) Name() string {
	return ServiceName
}

func (s *Service) autoclose() {
	s.mutexCloser.Lock()
	defer s.mutexCloser.Unlock()

	if s.closer == nil {
		return
	}

	if s.lastSend == nil || time.Since(*s.lastSend) < 40*time.Second {
		return
	}

	err := (*s.closer).Close()
	if err != nil {
		s.log.Warnf("Can't close smtp connection: %s", err)
	}

	s.closer = nil
}

// Start starts the SMTP-Sender in a separate thread
func (s *Service) Start() error {
	s.stopBroadcaster = broadcaster.NewBool(false)

	s.dialer = mail.NewDialer(*smtpHost, *smtpPort, *smtpUser, *smtpPassword)
	s.dialer.Timeout = 35 * time.Second
	s.dialer.RetryFailure = true

	s.log.Infof("SMTP-Service started, ready to send emails")

	s.runningFuncs.Add(1)
	go func() {
		defer s.runningFuncs.Done()

		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		stop, subStop := s.stopBroadcaster.Subscribe()
		defer subStop.Unsubscribe()

		for {
			select {
			case <-stop:
				return
			// Close the connection to the SMTP server if no email was sent in
			// the last 30 seconds.
			case <-ticker.C:
				s.autoclose()
			}
		}
	}()

	return nil
}

// Stop stops the SMTP-Sender thread
func (s *Service) Stop() error {
	s.stopBroadcaster.Next(true)

	s.runningFuncs.Wait()

	s.mutexCloser.Lock()
	defer s.mutexCloser.Unlock()

	if s.closer != nil {
		(*s.closer).Close()
	}

	return nil
}

// Health checks if the MailService is healthy
func (s *Service) Health() error {
	if s.error != nil {
		return s.error
	}

	return nil
}

// SendEmail sents a mail via SMTP
func (s *Service) SendEmail(m *Email) error {
	var err error

	msg := mail.NewMessage()

	from := m.From
	if from == "" {
		from = *smtpFrom
	}

	msg.SetHeader("From", from)
	msg.SetHeader("To", m.To)
	msg.SetHeader("Subject", m.Subject)

	if m.BodyPlain != "" {
		msg.SetBody("text/plain", m.BodyPlain)
	}

	if m.BodyHTML != "" {
		msg.SetBody("text/html", m.BodyHTML)
	}

	if m.Attachements != nil {
		for i := range m.Attachements {
			attachement := m.Attachements[i]
			reader := bytes.NewReader(attachement.Content)

			if attachement.Embedded {
				msg.EmbedReader(attachement.Filename, reader)
			} else {
				msg.AttachReader(attachement.Filename, reader)
			}
		}
	}

	s.mutexCloser.Lock()
	defer s.mutexCloser.Unlock()

	if s.closer == nil {
		closer, err := s.dialer.Dial()
		if err != nil {
			s.error = nil
			return err
		}

		s.error = nil
		now := time.Now()
		s.lastSend = &now
		s.closer = &closer
	}

	err = mail.Send(*s.closer, msg)
	if err != nil {
		if s.closer != nil {
			(*s.closer).Close()
			s.closer = nil
		}

		return err
	}

	return nil
}

// NewService if the ServiceFactory for an initialized Service
func NewService(ctx gousu.IContext) gousu.IService {
	return &Service{
		log: logger.GetLogger(fmt.Sprintf("service.%s", ServiceName)),
	}
}

// Assert NewService fullfills gousu.ServiceFactory
var _ (gousu.ServiceFactory) = NewService
