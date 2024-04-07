package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/deepaksing/FampayAssignment/store"
)

type DB struct {
	db *sql.DB
}

func NewDB() (store.Driver, error) {
	// postgres DSN
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}
	var driver store.Driver = &DB{
		db: db,
	}
	return driver, nil
}

func (d *DB) Migrate(ctx context.Context) error {
	buf, err := os.ReadFile("store/db/postgres/SCHEMA.sql")
	if err != nil {
		return fmt.Errorf("failed to read latest schema file: %w", err)
	}
	stmt := string(buf)
	_, err = d.db.ExecContext(ctx, stmt)
	if err != nil {
		return err
	}
	return nil
}

func (d *DB) StoreVideosInDB(videos []*store.Video) error {
	stmt, err := d.db.Prepare("INSERT INTO videos(title, description, published_at, thumbnails) VALUES($1, $2, $3, $4)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, video := range videos {
		_, err := stmt.Exec(video.Title, video.Description, video.PublishedAt, video.Thumbnails)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *DB) GetVideo(ctx context.Context, pageNum int, pageSize int) ([]*store.Video, error) {
	rows, err := d.db.Query("SELECT id, title, description, published_at, thumbnails FROM videos ORDER BY published_at DESC LIMIT $1 OFFSET $2", pageSize, (pageNum-1)*pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []*store.Video
	for rows.Next() {
		var video store.Video
		err := rows.Scan(&video.ID, &video.Title, &video.Description, &video.PublishedAt, &video.Thumbnails)
		if err != nil {
			return nil, err
		}
		videos = append(videos, &video)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return videos, nil
}

func (d *DB) SearchVideo(ctx context.Context, query string) ([]*store.Video, error) {
	matchedVideos := []*store.Video{}

	rows, err := d.db.Query("SELECT title, description FROM videos WHERE title ILIKE '%' || $1 || '%' OR description ILIKE '%' || $1 || '%'", query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var v store.Video
		if err := rows.Scan(&v.Title, &v.Description); err != nil {
			return nil, err
		}
		matchedVideos = append(matchedVideos, &v)
	}

	return matchedVideos, nil
}
