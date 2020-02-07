package service

import (
	"log"

	"github.com/adrianosela/rdtp/controller"
	"github.com/pkg/errors"
)

// Service represents an executable program
type Service struct {
	ctrl *controller.Controller
}

// NewService is the controller for the service
func NewService() *Service {
	return &Service{ctrl: controller.NewController()}
}

// Start runs the rdtp service without blocking
func (s *Service) Start() {
	go s.ctrl.Start()
}

// Run runs the rdtp service blocking execution
func (s *Service) Run() {
	if err := s.ctrl.Start(); err != nil {
		log.Println(errors.Wrap(err, "could not start controller"))
	}
}

// Stop gracefully shuts down the service
func (s *Service) Stop() {
	s.ctrl.Shutdown()
}
