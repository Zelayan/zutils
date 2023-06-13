package lock

import "context"

type Lock interface {
	TryLock() error
	UnLock() error
	TryLockCtx(ctx context.Context) error
	UnLockCtx(ctx context.Context) error
	Lock(ctx context.Context) error
}
