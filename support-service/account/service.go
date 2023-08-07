package account

import (
	context "context"
	support_service "github.com/regen-network/keystone/support-service"
)

type service struct {
	*support_service.UnimplementedAccountServiceServer
}

var _ support_service.AccountServiceServer = service{}

func (s service) Bootstrap(ctx context.Context, request *support_service.AccountServiceBootstrapRequest) (*support_service.AccountServiceBootstrapResponse, error) {
	return nil, nil
}
