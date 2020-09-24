package services

import (
	"context"

	"github.com/nhsh1997/be-test/config"
	"github.com/nhsh1997/be-test/internal/stores"
	"github.com/nhsh1997/be-test/pb"
)

const Version = "1.0.0"

type service struct {
	isReady   bool
	cfg       *config.Config
	mainStore *stores.MainStore
}

func (s *service) Sum(ctx context.Context, request *pb.SumRequest) (*pb.SumResponse, error) {
	panic("implement me")
}

func New(config *config.Config,
	mainStore *stores.MainStore) pb.BeTestServer {
	return &service{
		isReady:   true,
		cfg:       config,
		mainStore: mainStore,
	}
}

func (s *service) Version(context context.Context, req *pb.VersionRequest) (*pb.VersionResponse, error) {
	return &pb.VersionResponse{
		Version: Version,
	}, nil
}
