package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/dhiliphari/snippetbox/internal/models"
	"github.com/go-playground/form/v4"

	_ "github.com/go-sql-driver/mysql"
)

type config struct {
	dsn       string
	addr      string
	staticDir string
}

type application struct {
	logger         *slog.Logger
	cfg            *config
	snippets       *models.SnippetModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	// Load app configuration
	cfg := parseCfg()

	// Setup structured logging
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Connect to database
	db, err := openDB(cfg.dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	// Create template cache
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Form decoder
	formDecoder := form.NewDecoder()

	// Setup session state manager
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	app := &application{
		cfg:            cfg,
		logger:         logger,
		snippets:       &models.SnippetModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	logger.Info("starting server", slog.String("addr", cfg.addr))

	srv := &http.Server{
		Addr:     cfg.addr,
		Handler:  app.routes(),
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	logger.Error(err.Error())
	os.Exit(1)
}

func parseCfg() *config {
	var cfg config

	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.StringVar(&cfg.dsn, "dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	return &cfg
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
