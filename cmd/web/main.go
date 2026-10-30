package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

type config struct {
	addr      string
	env       string
	staticDir string
}

type application struct {
	logger *slog.Logger
}

var cfg config

func init() {
	const (
		defaultAddr        = ":8080"
		addrUsage          = "HTTP network address"
		defaultStaticDir   = "./assets/static/"
		defaultEnvironment = "dev"
	)

	flag.StringVar(&cfg.addr, "addr", defaultAddr, addrUsage)
	flag.StringVar(&cfg.addr, "a", defaultAddr, addrUsage+" (shorthand)")

	flag.StringVar(&cfg.staticDir, "static-dir", defaultStaticDir, "Path to static assets")

	flag.StringVar(&cfg.env, "environment", defaultEnvironment, "Application environment")
	flag.StringVar(&cfg.env, "env", defaultEnvironment, "Application environment (shorthand)")
}

func main() {
	flag.Parse()

	var handler slog.Handler
	if cfg.env == "production" {
		handler = slog.NewJSONHandler(os.Stdout, nil)
	} else if cfg.env == "debug" {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		})
	} else {
		handler = slog.NewTextHandler(os.Stdout, nil)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	app := &application{
		logger: logger,
	}

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(cfg.staticDir))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.snippetCreate)
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	logger.Info("starting server", "addr", cfg.addr)

	err := http.ListenAndServe(cfg.addr, mux)
	logger.Error(err.Error())
	os.Exit(1)
}
