package main

import (
	"context"
	"os"
	"os/signal"
	"sync"

	"github.com/LightAlykard/LittleLinksByAnthony/api/handler"
	"github.com/LightAlykard/LittleLinksByAnthony/api/server"
	"github.com/LightAlykard/LittleLinksByAnthony/app/repos/link"
	"github.com/LightAlykard/LittleLinksByAnthony/app/starter"
	"github.com/LightAlykard/LittleLinksByAnthony/db/mem/linkmemstore"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	lkst := linkmemstore.NewLinks()
	a := starter.NewApp(lkst)
	us := link.NewLinks(lkst)
	h := handler.NewRouter(us)
	srv := server.NewServer(":8000", h)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go a.Serve(ctx, wg, srv)

	<-ctx.Done()
	cancel()
	wg.Wait()
}
