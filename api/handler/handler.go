package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/LightAlykard/LittleLinksByAnthony/app/repos/link"
	"github.com/google/uuid"
)

type Handlers struct {
	// *http.ServeMux
	lk *link.Links
}

type Link struct {
	ID        uuid.UUID `json:"id"`
	ShortLink string    `json:"shortlink"`
	LongLink  string    `json:"longlink"`
	Data      string    `json:"data"`
}

func NewHandlers(lk *link.Links) *Handlers {
	r := &Handlers{
		// ServeMux: http.NewServeMux(),
		lk: lk,
	}
	// r.HandleFunc("/create", r.AuthMiddleware(http.HandlerFunc(r.CreateLink)).ServeHTTP)
	// r.HandleFunc("/read", r.AuthMiddleware(http.HandlerFunc(r.ReadLink)).ServeHTTP)
	// r.HandleFunc("/delete", r.AuthMiddleware(http.HandlerFunc(r.DeleteLink)).ServeHTTP)
	return r
}

func (rt *Handlers) CreateLink(ctx context.Context, l Link) (Link, error) {

	bl := link.Link{
		ShortLink: l.ShortLink,
		LongLink:  l.LongLink,
		Data:      l.Data,
	}

	nbl, err := rt.lk.Create(ctx, bl)
	if err != nil {
		// http.Error(w, "error when creating", http.StatusInternalServerError)
		return Link{}, fmt.Errorf("error when creating: %w", err)
	}

	// w.WriteHeader(http.StatusCreated)

	return Link{
		ID:        nbl.ID,
		ShortLink: nbl.ShortLink,
		LongLink:  nbl.LongLink,
		Data:      nbl.Data,
	}, nil
}

var ErrLinkNotFound = errors.New("link not found")

func (rt *Handlers) ReadLink(ctx context.Context, lid uuid.UUID) (Link, error) {

	if (lid == uuid.UUID{}) {
		return Link{}, fmt.Errorf("bad request: lid is empty")
	}

	nbl, err := rt.lk.Read(ctx, lid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Link{}, ErrLinkNotFound
		}
		return Link{}, fmt.Errorf("error when reading: %w", err)
	}

	return Link{
		ID:        nbl.ID,
		ShortLink: nbl.ShortLink,
		LongLink:  nbl.LongLink,
		Data:      nbl.Data,
	}, nil
}

func (rt *Handlers) DeleteLink(ctx context.Context, lid uuid.UUID) (Link, error) {

	if (lid == uuid.UUID{}) {
		return Link{}, fmt.Errorf("bad request: lid is empty")
	}

	nbl, err := rt.lk.Delete(ctx, lid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Link{}, ErrLinkNotFound
		}
		return Link{}, fmt.Errorf("error when reading: %w", err)
	}

	return Link{
		ID:        nbl.ID,
		ShortLink: nbl.ShortLink,
		LongLink:  nbl.LongLink,
		Data:      nbl.Data,
	}, nil
}
