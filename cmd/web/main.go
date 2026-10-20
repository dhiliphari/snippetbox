package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

type config struct {
	addr      string
	staticDir string
}

type application struct {
	logger *slog.Logger
	cfg    *config
}

func main() {
	cfg := parseCfg()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &application{
		cfg:    cfg,
		logger: logger,
	}

	logger.Info("starting server", slog.String("addr", cfg.addr))

	err := http.ListenAndServe(cfg.addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

func parseCfg() *config {
	var cfg config

	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.Parse()

	return &cfg
}
