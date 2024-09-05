package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pretodev/inhireapp/config/db"
	"github.com/pretodev/inhireapp/config/env"
	"github.com/pretodev/inhireapp/internal/inhire"
	"github.com/pretodev/inhireapp/pkg/browser"
	"github.com/sethvargo/go-envconfig"
)

func main() {
	downjobscdm := flag.NewFlagSet("download-jobs", flag.ExitOnError)
	downlinkscdm := flag.NewFlagSet("download-links", flag.ExitOnError)
	jobscdm := flag.NewFlagSet("jobs", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("Comando esperado: update ou find")
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		fmt.Println("\nSinal de interrupção recebido, encerrando...")
		cancel()
	}()

	var cfg env.Config
	if err := envconfig.Process(ctx, &cfg); err != nil {
		log.Fatalf("failed to load env variables: %v", err)
	}

	dblocal, err := db.OpenConn(ctx, cfg)
	if err != nil {
		log.Fatalf("failted to connect database: %v", err)
	}

	instore := inhire.NewStore(dblocal)
	insrv := inhire.NewService(instore)

	commands := map[string]func(ctx context.Context) error{
		"download-jobs": func(ctx context.Context) error {
			downjobscdm.Parse(os.Args[2:])
			browserCtx, cancel := browser.WithBrowserContext(ctx)
			defer cancel()
			return insrv.UpdateJobInfos(browserCtx)
		},
		"download-links": func(ctx context.Context) error {
			downlinkscdm.Parse(os.Args[2:])
			links, err := insrv.UpdateJobLinks(ctx)
			if err != nil {
				return err
			}
			log.Printf("founded %d jobs linux", len(links))
			return nil
		},
		"jobs": func(ctx context.Context) error {
			jobscdm.Parse(os.Args[2:])
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

	cmdName := os.Args[1]
	cmd, exists := commands[cmdName]
	if !exists {
		log.Printf("Comando não reconhecido: %s\n", cmdName)
		os.Exit(1)
	}

	if err := cmd(ctx); err != nil {
		log.Fatalf("Erro ao executar o comando %s: %v\n", cmdName, err)
		os.Exit(2)
	}

	fmt.Println("Encerrando...")
	time.Sleep(1 * time.Second)
}
