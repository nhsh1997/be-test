//  AUTO-GENERATED: DO NOT EDIT

package services

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	commonpb "git.begroup.team/platform-core/be-central-proto/common"
	"git.begroup.team/platform-core/kitchen/errors"
	"github.com/nhsh1997/be-test/config"
)

type ReadinessCheck func() bool

func DefaultReadinessCheck() bool {
	return true
}

type baseService struct {
	isReady        bool
	config         *config.Config
	readinessCheck ReadinessCheck
}

func NewBase(config *config.Config) commonpb.BaseServer {
	return &baseService{
		isReady:        true,
		config:         config,
		readinessCheck: DefaultReadinessCheck,
	}
}

func (s *baseService) Liveness(context context.Context, req *commonpb.LivenessRequest) (*commonpb.LivenessResponse, error) {
	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-osSignal:
		return nil, errors.Error(errors.Unavailable, "Server is shutting down")
	default:
		return &commonpb.LivenessResponse{Message: "OK"}, nil
	}
}
func (s *baseService) ToggleReadiness(context context.Context, req *commonpb.ToggleReadinessRequest) (*commonpb.ToggleReadinessResponse, error) {
	s.isReady = !s.isReady
	return &commonpb.ToggleReadinessResponse{Message: "OK"}, nil
}
func (s *baseService) Readiness(context context.Context, req *commonpb.ReadinessRequest) (*commonpb.ReadinessResponse, error) {
	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-osSignal:
		return nil, errors.Error(errors.Unavailable, "Server is shutting down")
	default:

		// err := s.mainStore.RelationalDatabaseCheck()

		if s.readinessCheck() == false {
			return nil, errors.Error(errors.Unavailable, "Server is not available, status: mainStore error")
		}

		if s.isReady {
			return &commonpb.ReadinessResponse{Message: "OK"}, nil
		}

		return nil, errors.Error(errors.Unavailable, "Server isn't ready, status: toggle off")
	}
}
