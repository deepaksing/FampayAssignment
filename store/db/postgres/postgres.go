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
