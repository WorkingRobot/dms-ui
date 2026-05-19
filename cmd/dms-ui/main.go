// Command dms-ui is a lightweight web admin panel for docker-mailserver.
// It wraps the container's `setup` CLI (plus doveadm/supervisorctl) over a
// mounted Docker socket and serves an embedded SvelteKit UI. It has no
// built-in auth and is designed to run behind Traefik + Authentik.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/WorkingRobot/dms-ui/internal/api"
	"github.com/WorkingRobot/dms-ui/internal/dms"
	"github.com/WorkingRobot/dms-ui/internal/web"
)

func main() {
	cfg, err := api.LoadConfig()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	client, err := dms.New(cfg.MailContainer)
	if err != nil {
		log.Fatalf("dms client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	if err := client.Ping(ctx); err != nil {
		log.Printf("warning: cannot reach mail container %q yet: %v", cfg.MailContainer, err)
	} else {
		log.Printf("connected to mail container %q", cfg.MailContainer)
	}
	cancel()

	srv := &http.Server{
		Addr:              cfg.ListenAddr,
		Handler:           api.NewServer(cfg, client, web.FS()),
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("dms-ui listening on %s (container=%s)",
			cfg.ListenAddr, cfg.MailContainer)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("shutting down")
	shutCtx, shutCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutCancel()
	_ = srv.Shutdown(shutCtx)
}
