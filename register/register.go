package register

import "context"

type Service struct {
	ID   string
	Name string
	IP   string
	Port int32
}

type Register interface {
	Register(ctx context.Context, service *Service) error
	UnRegister(ctx context.Context, service *Service) error
}
