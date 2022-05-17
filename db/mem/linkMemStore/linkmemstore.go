package linkmemstore

import (
	"context"
	"database/sql"
	"sync"

	"github.com/LightAlykard/LittleLinksByAnthony/app/repos/link"
	"github.com/google/uuid"
)

//var _ link.LinkStore = &Links{}

type Links struct {
	sync.Mutex
	m map[uuid.UUID]link.Link
}

func NewLinks() *Links {
	return &Links{
		m: make(map[uuid.UUID]link.Link),
	}
}

func (lk *Links) Create(ctx context.Context, l link.Link) (*uuid.UUID, error) {
	lk.Lock()
	defer lk.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	lid := uuid.New()
	l.ID = lid
	lk.m[l.ID] = l
	return &lid, nil
}

func (lk *Links) Read(ctx context.Context, lid uuid.UUID) (*link.Link, error) {
	lk.Lock()
	defer lk.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	u, ok := lk.m[lid]
	if ok {
		return &u, nil
	}
	return nil, sql.ErrNoRows
}

func (lk *Links) Delete(ctx context.Context, lid uuid.UUID) (*link.Link, error) {
	lk.Lock()
	defer lk.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	delete(lk.m, lid)
	return nil, nil
}
