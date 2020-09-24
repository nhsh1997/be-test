package main

import (
	commonpb "git.begroup.team/platform-core/be-central-proto/common"
	"github.com/nhsh1997/be-test/config"
	"github.com/nhsh1997/be-test/internal/services"
)

func registerBaseService(cfg *config.Config) commonpb.BaseServer {
	return services.NewBase(cfg)
}
