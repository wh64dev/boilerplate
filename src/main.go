package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"git.wh64.net/devproje/devproje-boilerplate/src/config"
	"git.wh64.net/devproje/devproje-boilerplate/src/modules"
	"git.wh64.net/devproje/devproje-boilerplate/src/modules/database"
	"git.wh64.net/devproje/devproje-boilerplate/src/modules/sample"
	"github.com/devproje/commando"
	"github.com/gin-gonic/gin"
)

func cli(args []string) {
	var command = commando.NewCommando(args)
	err := command.Execute()
	if err != nil {
		log.Println("error")
		return
	}
}

func main() {
	var args = os.Args
	if len(args) > 1 {
		cli(args[1:])
		return
	}

	cnf := config.Get
	app := gin.Default()

	modules.LOADER.Insmod(database.DatabaseModule)
	modules.LOADER.Insmod(sample.SampleServiceModule)
	modules.LOADER.Load()

	var webserver = &http.Server{
		Handler:           app,
		ReadHeaderTimeout: 5 * time.Second,
		Addr:              fmt.Sprintf("%s:%d", cnf.Host, cnf.Port),
	}

	go func() {
		fmt.Printf("WebServer bind at http://%s:%d\n", cnf.Host, cnf.Port)
		if err := webserver.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	var sig = make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig

	log.Printf("shutting down server...\n")
	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := webserver.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v\n", err)
	}

	modules.LOADER.Unload()
	log.Println("Server closed")
}
