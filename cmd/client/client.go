package client

import (
	"context"
	"git.begroup.team/platform-core/kitchen/l"
	"github.com/nhsh1997/be-test/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
)

type BeTestClient struct {
	pb.BeTestClient
}

func NewBeTestClient(address string) *BeTestClient {
	conn, err := grpc.DialContext(context.Background(), address,
		grpc.WithInsecure(),
		grpc.WithBalancerName(roundrobin.Name),
	)

	if err != nil {
		ll.Fatal("Failed to dial BeTest service", l.Error(err))
	}

	c := pb.NewBeTestClient(conn)

	return &BeTestClient{c}
}
