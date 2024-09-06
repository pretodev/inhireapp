package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/pretodev/inhireapp/config/db"
	"github.com/pretodev/inhireapp/config/env"
	"github.com/pretodev/inhireapp/internal/inhire"
	"github.com/pretodev/inhireapp/pkg/browser"
	"github.com/sethvargo/go-envconfig"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("expected command: update|find|version")
		os.Exit(1)
	}

	cmdName := os.Args[1]
	if cmdName == "version" {
		fmt.Println("inhire web crawler: v0.0.1")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		fmt.Println("shutdown...")
		cancel()
	}()

	var cfg env.Config
	if err := envconfig.Process(ctx, &cfg); err != nil {
		log.Fatalf("failed to load env variables: %v", err)
	}

	db, err := db.OpenConn(ctx, cfg)
	defer func() {
		if errClose := db.Close(); errClose != nil {
			log.Fatalf("failed close database: %v", errClose)
		}
	}()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	instore := inhire.NewStore(db)
	insrv := inhire.NewService(instore)

	commands := map[string]func(ctx context.Context) error{
		"download-jobs": func(ctx context.Context) error {
			browserCtx, cancel := browser.WithBrowserContext(ctx)
			defer cancel()
			return insrv.UpdateJobInfos(browserCtx)
		},
		"download-links": func(ctx context.Context) error {
			links, err := insrv.UpdateJobLinks(ctx)
			if err != nil {
				return err
			}
			log.Printf("founded %d jobs linux", len(links))
			return nil
		},
		"jobs": func(ctx context.Context) error {
			jobs, err := instore.CachedJobs(ctx)
			if err != nil {
				return err
			}
			for idx, j := range jobs {
				fmt.Printf("%d - %s: %s\n", idx, j.PositionName, j.PageURL)
			}
			return nil
		},
	}

	cmd, exists := commands[cmdName]
	if !exists {
		log.Printf("not found command: %s\n", cmdName)
		os.Exit(1)
	}

	if err := cmd(ctx); err != nil {
		log.Fatalf("Erro ao executar o comando %s: %v\n", cmdName, err)
		os.Exit(2)
	}
}
