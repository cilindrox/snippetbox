package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

const (
	defaultAddr        = ":8080"
	addrUsage          = "HTTP network address"
	defaultStaticDir   = "./assets/static/"
	defaultEnvironment = "dev"
)

type config struct {
	addr      string
	env       string
	staticDir string
}

type application struct {
	config config
	logger *slog.Logger
}

func main() {
	var cfg config
	flag.StringVar(&cfg.addr, "addr", defaultAddr, addrUsage)
	flag.StringVar(&cfg.addr, "a", defaultAddr, addrUsage+" (shorthand)")

	flag.StringVar(&cfg.staticDir, "static-dir", defaultStaticDir, "Path to static assets")

	flag.StringVar(&cfg.env, "environment", defaultEnvironment, "Application environment")
	flag.StringVar(&cfg.env, "env", defaultEnvironment, "Application environment (shorthand)")

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
		config: cfg,
	}

	logger.Info("starting server", "addr", cfg.addr)

	err := http.ListenAndServe(cfg.addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}
