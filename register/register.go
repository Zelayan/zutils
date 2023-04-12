package register

import "context"

type Service struct {
	ID      string
	Name    string
	IP      string
	Port    int32
	Version string
	// Endpoints are endpoint addresses of the service instance.
	// schema:
	//   http://127.0.0.1:8000?isSecure=false
	//   grpc://127.0.0.1:9000?isSecure=false
	EndPoint []string
}

type Register interface {
	Register(ctx context.Context, service *Service) error
	UnRegister(ctx context.Context, service *Service) error
}
