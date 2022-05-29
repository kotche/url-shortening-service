package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/kotche/url-shortening-service/internal/app/service"
	"github.com/kotche/url-shortening-service/internal/app/usecase"
)

type DB struct {
	conn *sql.DB
}

func NewDB(DSN string) (*DB, error) {
	conn, err := sql.Open("pgx", DSN)
	if err != nil {
		return nil, err
	}
	db := &DB{conn: conn}
	db.init()
	return db, nil
}

func (d *DB) Add(userID string, url *service.URL) error {
	ctx := context.Background()

	_, err := d.conn.ExecContext(ctx,
		"INSERT INTO public.users(user_id) VALUES ($1) ON CONFLICT (user_id) DO UPDATE SET user_id=EXCLUDED.user_id;", userID)
	if err != nil {
		return err
	}

	stmt, err := d.conn.PrepareContext(ctx,
		"INSERT INTO public.urls(short,origin,user_id) VALUES ($1,$2,$3) ON CONFLICT (origin,user_id) DO UPDATE SET origin=EXCLUDED.origin RETURNING short")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result := stmt.QueryRowContext(ctx, url.Short, url.Origin, userID)

	var output string
	result.Scan(&output)
	if output != url.Short {
		return usecase.ErrConflictURL{Err: errors.New("duplicate URL"), ShortenURL: output}
	}

	return nil
}

func (d *DB) GetByID(id string) (*service.URL, error) {
	var output sql.NullString
	row := d.conn.QueryRow("SELECT origin FROM public.urls WHERE short=$1", id)
	row.Scan(&output)
	if output.Valid && output.String != "" {
		url := service.NewURL(id, output.String)
		return url, nil
	} else {
		return nil, fmt.Errorf("key not found")
	}
}

func (d *DB) GetUserURLs(userID string) ([]*service.URL, error) {
	urls := make([]*service.URL, 0)

	rows, err := d.conn.Query("SELECT short, origin FROM public.urls WHERE user_id=$1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var url service.URL
		err = rows.Scan(&url.Short, &url.Origin)
		if err != nil {
			return nil, err
		}
		urls = append(urls, &url)
	}

	err = rows.Err()
	if err != nil {
		log.Println(err.Error())
	}

	return urls, nil
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

func (d *DB) WriteBatch(ctx context.Context, userID string, urls []*service.URL) error {

	tx, err := d.conn.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = d.conn.ExecContext(ctx, "INSERT INTO public.users(user_id) VALUES ($1) ON CONFLICT (user_id) DO UPDATE SET user_id=EXCLUDED.user_id;", userID)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx,
		"INSERT INTO public.urls(short,origin,user_id) VALUES ($1,$2,$3) ON CONFLICT (origin,user_id) DO UPDATE SET origin=EXCLUDED.origin RETURNING short")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, url := range urls {
		result := stmt.QueryRowContext(ctx, url.Short, url.Origin, userID)
		var output sql.NullString
		result.Scan(&output)
		if output.Valid {
			url.Short = output.String
		} else {
			log.Printf("URL short:%s, origin:%s no valid", url.Short, url.Origin)
		}
	}

	return tx.Commit()
}

func (d *DB) init() {
	_, err := d.conn.Exec(`CREATE TABLE IF NOT EXISTS public.users(
		    user_id VARCHAR(500) NOT NULL PRIMARY KEY
		);

		CREATE TABLE IF NOT EXISTS public.urls(
		short VARCHAR(50) NOT NULL PRIMARY KEY,
		origin VARCHAR(500) NOT NULL,
		user_id VARCHAR(500) NOT NULL,
    	CONSTRAINT uniq_origin_user_id UNIQUE (origin, user_id),
    	FOREIGN KEY (user_id) REFERENCES public.users (user_id));`)

	if err != nil {
		log.Fatal(err.Error())
	}
}
