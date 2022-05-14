package link

import (
	"fmt"

	"github.com/google/uuid"
)

type Link struct {
	ID        uuid.UUID
	ShortLink string
	LongLink  string
	//LiveStatus bool
}

type LinkStore interface {
	Create(l Link) (*uuid.UUID, error)
	Read(lid uuid.UUID) (*Link, error)
	Delete(lid uuid.UUID) error
	// SearchLinks(s string) (chan Link, error)
}

type Links struct {
	lkstore LinkStore
}

func NewLinks(lkstore LinkStore) *Links {
	return &Links{
		lkstore: lkstore,
	}
}

func (lks *Links) Create(lk Link) (*Link, error) {
	id, err := lks.lkstore.Create(lk)
	if err != nil {
		return nil, fmt.Errorf("create link error: %w", err)
	}
	lk.ID = *id
	return &lk, nil
}

func (lks *Links) Read(lid uuid.UUID) (*Link, error) {
	lk, err := lks.lkstore.Read(lid)
	if err != nil {
		return nil, fmt.Errorf("read link error: %w", err)
	}
	return lk, nil
}

func (lks *Links) Delete(lid uuid.UUID) (*Link, error) {
	lk, err := lks.lkstore.Read(lid)
	if err != nil {
		return nil, fmt.Errorf("create link error: %w", err)
	}

	return lk, lks.lkstore.Delete(lid)
}

// func (lks *Links) SearchLinks(s string) (chan Link, error) {
// 	chin, err := lks.lkstore.SearchLinks(s)
// 	if err!=nil{
// 		return nil, err
// 	}
// 	chout:=make(chan Link,100)
// 	go func ()  {
// 		defer close(chout)
// 		for {
// 			lk,ok:=<-chin
// 			if !ok{
// 				return
// 			}
// 			lk.LiveStatus = 0755
// 			chout <- lk
// 		}
// 	}()
// 	return chout, nil
// }
