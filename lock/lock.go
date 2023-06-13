package lock

import "context"

type Lock interface {
	TryLock() (bool, error)
	UnLock() (bool, error)
	TryLockCtx(ctx context.Context) (bool, error)
	UnLockCtx(ctx context.Context) (bool, error)
}
