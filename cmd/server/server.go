package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"
	"strings"

	commonpb "git.begroup.team/platform-core/be-central-proto/common"
	"git.begroup.team/platform-core/kitchen/l"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/nhsh1997/be-test/config"
	"github.com/nhsh1997/be-test/pb"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Server struct {
	http.Server
	cfg          *config.Config
	Addrs        []string // addresses on which the server listens for new connection.
	inShutdown   uint32   // indicates whether the server is in shutdown or not
	requestCount int32    // counter holds no. of request in progress.
}

// NewServer create new server using given arguments
func NewServer(cfg *config.Config) *Server {
	return &Server{
		cfg: cfg,
	}
}

// RunGRPCGateway will start an GRPC Gateway
func (s *Server) RunGRPCGateway() (err error) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: true, EmitDefaults: true}))
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err = pb.RegisterBeTestHandlerFromEndpoint(ctx, mux, fmt.Sprintf(":%d", s.cfg.GRPCAddress), opts)
	if err != nil {
		return err
	}

	err = commonpb.RegisterBaseHandlerFromEndpoint(ctx, mux, fmt.Sprintf(":%d", s.cfg.GRPCAddress), opts)
	if err != nil {
		return err
	}

	muxHttp := http.NewServeMux()
	muxHttp.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		promhttp.Handler().ServeHTTP(w, r)
	})

	muxHttp.HandleFunc("/log/level", func(w http.ResponseWriter, r *http.Request) {
		l.ServeHTTP(w, r)
	})

	muxHttp.HandleFunc("/debug/pprof/", pprof.Index)
	muxHttp.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	muxHttp.HandleFunc("/debug/pprof/profile", pprof.Profile)
	muxHttp.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	muxHttp.HandleFunc("/debug/pprof/trace", pprof.Trace)

	muxHttp.Handle("/", forwardAccessToken(mux))

	return http.ListenAndServe(fmt.Sprintf(":%d", s.cfg.HTTPAddress), muxHttp)
}

func forwardAccessToken(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		md := make(metadata.MD)
		for k := range r.Header {
			k2 := strings.ToLower(k)
			md[k2] = []string{r.Header.Get(k)}
		}
		ctx := metadata.NewIncomingContext(r.Context(), md)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
}
