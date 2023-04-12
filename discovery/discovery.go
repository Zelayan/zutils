package discovery

import (
	"context"
	"github.com/Zelayan/zutils/register"
)

type Discovery interface {
	GetService(ctx context.Context, serviceName string) ([]*register.Service, error)
}
