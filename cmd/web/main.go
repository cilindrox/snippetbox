package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/cilindrox/snippetbox/internal/models"

	"github.com/go-playground/form/v4"

	_ "github.com/go-sql-driver/mysql"
)

const (
	defaultAddr        = ":8080"
	addrUsage          = "HTTP network address"
	defaultStaticDir   = "./assets/static/"
	defaultEnvironment = "dev"
)

type config struct {
	addr      string
	dsn       string
	env       string
	staticDir string
}

type application struct {
	config         config
	logger         *slog.Logger
	snippets       *models.SnippetModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	var cfg config
	flag.StringVar(&cfg.addr, "addr", defaultAddr, addrUsage)
	flag.StringVar(&cfg.addr, "a", defaultAddr, addrUsage+" (shorthand)")

	flag.StringVar(&cfg.staticDir, "static-dir", defaultStaticDir, "Path to static assets")

	flag.StringVar(&cfg.env, "environment", defaultEnvironment, "Application environment")
	flag.StringVar(&cfg.env, "env", defaultEnvironment, "Application environment (shorthand)")

	flag.StringVar(&cfg.dsn, "dsn", "web:pass@/snippetbox?parseTime=true", "MySQL datasource name")

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

	db, err := openDB(cfg.dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.HashTokenInStore = true

	app := &application{
		config:         cfg,
		logger:         logger,
		snippets:       &models.SnippetModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	srv := &http.Server{
		Addr:                cfg.addr,
		Handler:             app.routes(),
		ErrorLog:            slog.NewLogLogger(logger.Handler(), slog.LevelError),
		MaxHeaderBytes:      524288, // 0.5Mi
		MaxHeaderValueCount: 100,    // return 413 past this threshold
		IdleTimeout:         time.Minute,
		ReadTimeout:         5 * time.Second,
		WriteTimeout:        10 * time.Second,
	}

	logger.Info("starting server", "addr", srv.Addr)

	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
