package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"contrib.go.opencensus.io/exporter/stackdriver"
	kgrpc "git.begroup.team/platform-core/kitchen/grpc"
	"git.begroup.team/platform-core/kitchen/id"
	"git.begroup.team/platform-core/kitchen/l"
	"google.golang.org/grpc"

	commonpb "git.begroup.team/platform-core/be-central-proto/common"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/nhsh1997/be-test/config"
	"github.com/nhsh1997/be-test/pb"

	"go.opencensus.io/trace"
)

var (
	ll = l.New()
)

func main() {
	// load configs
	cfg := config.Load()
	id.SetGEnv(cfg.Environment)

	// setup google cloud tracing
	if cfg.Tracing.Enabled {
		exporter, err := stackdriver.NewExporter(stackdriver.Options{})
		if err != nil {
			ll.Fatal("failed to new exporter", l.Error(err))
		}
		trace.RegisterExporter(exporter)
	}

	// grpc server
	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_prometheus.UnaryServerInterceptor,
			kgrpc.LogUnaryServerInterceptor(ll),
		)),
	)

	grpc_prometheus.EnableHandlingTimeHistogram()
	grpc_prometheus.Register(s)
	// Register Prometheus metrics handler.

	svc := registerService(cfg)
	pb.RegisterBeTestServer(s, svc)

	commonsvc := registerBaseService(cfg)
	commonpb.RegisterBaseServer(s, commonsvc)

	// handle signal
	_, ctxCancel := context.WithCancel(context.Background())
	go func() {
		osSignal := make(chan os.Signal, 1)
		signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM)
		<-osSignal
		ctxCancel()
		// Wait for maximum 15s
		go func() {
			var durationSec time.Duration = 15
			if cfg.Environment == "D" {
				durationSec = 1
			}
			timer := time.NewTimer(durationSec * time.Second)
			<-timer.C
			ll.Fatal("Force shutdown due to timeout!")
		}()
	}()

	go func() {
		gw := NewServer(cfg)
		ll.Info("HTTP server start listening", l.Int("HTTP address", cfg.HTTPAddress))
		err := gw.RunGRPCGateway()
		if err != nil {
			ll.Fatal("error listening to address", l.Int("address", cfg.HTTPAddress), l.Error(err))
			return
		}
	}()

	ll.Info("GRPC server start listening", l.Int("GRPC address", cfg.GRPCAddress))
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCAddress))
	if err != nil {
		ll.Fatal("error listening to address", l.Int("address", cfg.GRPCAddress), l.Error(err))
		return
	}

	s.Serve(listener)
}
