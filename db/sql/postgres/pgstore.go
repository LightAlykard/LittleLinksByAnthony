package postgres

import (
	"context"
	"database/sql"

	"time"

	"github.com/LightAlykard/LittleLinksByAnthony/app/repos/link"
	"github.com/google/uuid"
)

var _ link.LinkStore = &Links{}

type DBPgLink struct {
	ID        uuid.UUID  `db:"id"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
	ShortLink string     `db:"shortlink"`
	LongLink  string     `db:"longlink"`
	Data      string     `db:"data"`
}

type Links struct {
	db *sql.DB
}

func NewUsers(dsn string) (*Links, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS public.users (
		id uuid NOT NULL,
		created_at timestamptz NOT NULL,
		updated_at timestamptz NOT NULL,
		deleted_at timestamptz NULL,
		shortlink varchar NOT NULL,
		longlink varchar NOT NULL,
		"data" varchar NULL,
		CONSTRAINT users_pk PRIMARY KEY (id)
	);`)

	if err != nil {
		db.Close()
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	ls := &Links{
		db: db,
	}
	return ls, nil
}

func (us *Links) Close() {
	us.db.Close()
}

func (ls *Links) Create(ctx context.Context, l link.Link) (*uuid.UUID, error) {
	dbu := &DBPgLink{
		ID:        l.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ShortLink: l.ShortLink,
		LongLink:  l.LongLink,
		Data:      l.Data,
	}

	_, err := ls.db.ExecContext(ctx, `INSERT INTO links 
	(id, created_at, updated_at, deleted_at, shortlink, longlink, data)
	values ($1, $2, $3, $4, $5, $6, $7)`,
		dbu.ID,
		dbu.CreatedAt,
		dbu.UpdatedAt,
		nil,
		dbu.ShortLink,
		dbu.LongLink,
		dbu.Data,
	)
	if err != nil {
		return nil, err
	}

	return &l.ID, nil
}

func (ls *Links) Delete(ctx context.Context, uid uuid.UUID) error {
	_, err := ls.db.ExecContext(ctx, `UPDATE links SET deleted_at = $2 WHERE id = $1`,
		uid, time.Now(),
	)
	return err
}

func (ls *Links) Read(ctx context.Context, uid uuid.UUID) (*link.Link, error) {
	dbu := &DBPgLink{}
	rows, err := ls.db.QueryContext(ctx, `SELECT id, created_at, updated_at, deleted_at, shortlink, longlink, data 
	FROM links WHERE id = $1`, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(
			&dbu.ID,
			&dbu.CreatedAt,
			&dbu.UpdatedAt,
			&dbu.DeletedAt,
			&dbu.ShortLink,
			&dbu.LongLink,
			&dbu.Data,
		); err != nil {
			return nil, err
		}
	}

	return &link.Link{
		ID:        dbu.ID,
		ShortLink: dbu.ShortLink,
		LongLink:  dbu.LongLink,
		Data:      dbu.Data,
	}, nil
}
