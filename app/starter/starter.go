package starter

import (
	"context"
	"sync"

	"github.com/LightAlykard/LittleLinksByAnthony/app/repos/link"
)

type App struct {
	lk *link.Links
}

type HTTPServer interface {
	Start(lk *link.Link)
	Stop()
}

func NewApp(lkt link.LinkStore) *App {
	a := &App{
		lk: link.NewLinks(lkt),
	}
	return a
}

func (a *App) Serve(ctx context.Context, wg *sync.WaitGroup, hs HTTPServer) {
	defer wg.Done()
	hs.Start(a.lk)
	<-ctx.Done()
	hs.Stop()
}
