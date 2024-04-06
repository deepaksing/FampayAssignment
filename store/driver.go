package store

import "context"

type Driver interface {
	Migrate(ctx context.Context) error
}
