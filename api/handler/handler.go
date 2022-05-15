package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/LightAlykard/LittleLinksByAnthony/app/repos/link"
	"github.com/google/uuid"
)

type Router struct {
	*http.ServeMux
	lk *link.Links
}

type Link struct {
	ID        uuid.UUID `json:"id"`
	ShortLink string    `json:"shortlink"`
	LongLink  string    `json:"longlink"`
	Data      string    `json:"data"`
}

func NewRouter(lk *link.Links) *Router {
	r := &Router{
		ServeMux: http.NewServeMux(),
		lk:       lk,
	}
	r.HandleFunc("/create", r.AuthMiddleware(http.HandlerFunc(r.CreateLink)).ServeHTTP)
	r.HandleFunc("/read", r.AuthMiddleware(http.HandlerFunc(r.CreateLink)).ServeHTTP)
	r.HandleFunc("/delete", r.AuthMiddleware(http.HandlerFunc(r.CreateLink)).ServeHTTP)
	return r
}

func (rt *Router) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if u, p, ok := r.BasicAuth(); !ok || !(u == "admin" && p == "admin") {
				http.Error(w, "unautorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		},
	)
}

func (rt *Router) CreateLink(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	l := Link{}
	if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
		http.Error(w, "bad request", http.StatusUnauthorized)
		return
	}

	bl := &link.Link{
		ShortLink: l.ShortLink,
		LongLink:  l.LongLink,
		Data:      l.Data,
	}
	nbl, err := rt.lk.Create(r.Context(), bl)
	if err != nil {
		http.Error(w, "error when creating", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(
		Link{
			ID:        nbl.ID,
			ShortLink: nbl.ShortLink,
			LongLink:  nbl.LongLink,
			Data:      nbl.Data,
		},
	)

}

func (rt *Router) ReadLink(w http.ResponseWriter, r *http.Request) {
	slid := r.URL.Query().Get("lid")
	if slid == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	lid, err := uuid.Parse(slid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if (lid == uuid.UUID{}) {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	nbl, err := rt.lk.Read(r.Context(), lid)
	if err != nil {
		if errors.Is(err.sql.ErrNoRows) {
			http.Error(w, "not found", http.StatusNotFound)
		} else {
			http.Error(w, "error when reading", http.StatusInternalServerError)
		}
		return
	}

	_ = json.NewEncoder(w).Encode(
		Link{
			ID:        nbl.ID,
			ShortLink: nbl.ShortLink,
			LongLink:  nbl.LongLink,
			Data:      nbl.Data,
		},
	)
}

func (rt *Router) DeleteLink(w http.ResponseWriter, r *http.Request) {
	slid := r.URL.Query().Get("lid")
	if slid == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	lid, err := uuid.Parse(slid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if (lid == uuid.UUID{}) {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	nbl, err := rt.lk.Delete(r.Context(), lid)
	if err != nil {
		if errors.Is(err.sql.ErrNoRows) {
			http.Error(w, "not found", http.StatusNotFound)
		} else {
			http.Error(w, "error when reading", http.StatusInternalServerError)
		}
		return
	}

	_ = json.NewEncoder(w).Encode(
		Link{
			ID:        nbl.ID,
			ShortLink: nbl.ShortLink,
			LongLink:  nbl.LongLink,
			Data:      nbl.Data,
		},
	)
}
