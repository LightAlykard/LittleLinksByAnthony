package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/LightAlykard/LittleLinksByAnthony/api/handler"
	"github.com/LightAlykard/LittleLinksByAnthony/api/routergin"
	"github.com/LightAlykard/LittleLinksByAnthony/api/server"
	"github.com/LightAlykard/LittleLinksByAnthony/app/repos/link"
	"github.com/LightAlykard/LittleLinksByAnthony/app/starter"
	"github.com/LightAlykard/LittleLinksByAnthony/db/mem/linkmemstore"
	"github.com/LightAlykard/LittleLinksByAnthony/db/sql/pgstore"
)

func main() {
	if tz := os.Getenv("TZ"); tz != "" {
		var err error
		time.Local, err = time.LoadLocation(tz)
		if err != nil {
			log.Printf("error loading location '%s': %v\n", tz, err)
		}
	}

	// output current time zone
	tnow := time.Now()
	tz, _ := tnow.Zone()
	log.Printf("Local time zone %s. Service started at %s", tz,
		tnow.Format("2006-01-02T15:04:05.000 MST"))

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	var lst link.LinkStore
	stu := os.Getenv("REGUSER_STORE")

	switch stu {
	case "mem":
		lst = linkmemstore.NewLinks()
	case "pg":
		dsn := os.Getenv("DATABASE_URL")
		pgst, err := pgstore.NewUsers(dsn)
		if err != nil {
			log.Fatal(err)
		}
		defer pgst.Close()
		lst = pgst
	default:
		log.Fatal("unknown REGUSER_STORE = ", stu)
	}

	//lkst := linkmemstore.NewLinks()
	a := starter.NewApp(lst)
	ls := link.NewLinks(lst)
	h := handler.NewHandlers(ls)
	rh := routergin.NewRouterGin(h)
	//srv := server.NewServer(":8080", rh)
	srv := server.NewServer(":"+os.Getenv("PORT"), rh)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go a.Serve(ctx, wg, srv)

	<-ctx.Done()
	cancel()
	wg.Wait()
}
