package store

import "context"

type Driver interface {
	Migrate(ctx context.Context) error
	StoreVideosInDB(videos []*Video) error
	GetVideo(ctx context.Context, pageNum int, pageSize int) ([]*Video, error)
	SearchVideo(ctx context.Context, query string) ([]*Video, error)
}
