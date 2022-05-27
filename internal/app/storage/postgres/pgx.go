package postgres

import (
	"database/sql"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/kotche/url-shortening-service/internal/app/service"
)

type DB struct {
	conn *sql.DB
}

func NewDB(DSN string) (*DB, error) {
	db, err := sql.Open("pgx", DSN)
	if err != nil {
		return nil, err

	}
	return &DB{
		conn: db,
	}, nil
}

func (d *DB) Add(userID string, url *service.URL) error {
	return nil
}

func (d *DB) GetByID(id string) (*service.URL, error) {
	return nil, nil
}

func (d *DB) GetUserURLs(userID string) ([]*service.URL, error) {
	return nil, nil
}

func (d *DB) Close() error {
	if err := d.conn.Close(); err != nil {
		return err
	}
	return nil
}

func (d *DB) Ping() error {
	if err := d.conn.Ping(); err != nil {
		return err
	}
	return nil
}
