package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"git.wh64.net/devproje/devproje-boilerplate/config"
	"git.wh64.net/devproje/devproje-boilerplate/modules"
	"git.wh64.net/devproje/devproje-boilerplate/modules/database"
	"git.wh64.net/devproje/devproje-boilerplate/modules/sample"
	"git.wh64.net/devproje/devproje-boilerplate/routes"
	"github.com/devproje/commando"
	"github.com/gin-gonic/gin"
)

//go:embed public/*
var static embed.FS

func cli(args []string) bool {
	var command = commando.NewCommando(args)
	err := command.Execute()
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func mime(path string) string {
	switch {
	case strings.HasSuffix(path, ".ico"):
		return "image/x-icon"
	case strings.HasSuffix(path, ".txt"):
		return "text/plain"
	case strings.HasSuffix(path, ".svg"):
		return "image/svg+xml"
	case strings.HasSuffix(path, ".json"):
		return "application/json"
	default:
		return "application/octet-stream"
	}
}

func main() {
	var args = os.Args
	if len(args) > 1 {
		var ok = cli(args[1:])
		if !ok {
			return
		}
	}

	cnf := config.Get
	app := gin.Default()

	modules.LOADER.Insmod(database.DatabaseModule)
	modules.LOADER.Insmod(sample.SampleServiceModule)
	modules.LOADER.Load()

	routes.API(app)

	var public, _ = fs.Sub(static, "public")
	var assets, _ = fs.Sub(public, "assets")
	var index, _ = fs.ReadFile(public, "index.html")

	app.StaticFS("/assets", http.FS(assets))
	app.GET("/", func (ctx *gin.Context) {
		ctx.Data(200, "text/html; charset=utf-8", index)
	})

	app.NoRoute(func (ctx *gin.Context) {
		path := ctx.Request.URL.Path[1:]

		if data, err := fs.ReadFile(public, path); err == nil {
			ctx.Data(200, mime(path), data)
			return
		}

		ctx.Data(200, "text/html; charset=utf-8", index)
	})

	var webserver = &http.Server{
		Handler:           app,
		ReadHeaderTimeout: 5 * time.Second,
		Addr:              fmt.Sprintf("%s:%d", cnf.Host, cnf.Port),
	}

	go func() {
		log.Printf("WebServer bind at http://%s:%d\n", cnf.Host, cnf.Port)
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
