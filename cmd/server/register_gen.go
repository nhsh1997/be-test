// IM AUTO GENERATED, BUT CAN BE OVERRIDDEN

package main

import (
	"github.com/nhsh1997/be-test/config"
	"github.com/nhsh1997/be-test/internal/services"
	"github.com/nhsh1997/be-test/internal/stores"
	"github.com/nhsh1997/be-test/pb"
)

func registerService(cfg *config.Config) pb.BeTestServer {

	mainStore := stores.NewMainStore()

	return services.New(cfg, mainStore)
}
