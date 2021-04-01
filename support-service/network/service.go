package network

import (
	"context"
	support_service "github.com/regen-network/keystone/support-service"
)

type service struct {
	*support_service.UnimplementedNetworkServiceServer
	networks []*support_service.NetworkInfo
}

var _ support_service.NetworkServiceServer = service{}

func (s service) List(context.Context, *support_service.NetworkServiceListRequest) (*support_service.NetworkServiceListResponse, error) {
	return &support_service.NetworkServiceListResponse{
		Networks: s.networks,
	}, nil
}
