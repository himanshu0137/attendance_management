package main

import (
	"context"
	"flag"
	"log"
	"strings"

	_ "hunrx.com/attendance_management/modules/attendance"
	_ "hunrx.com/attendance_management/modules/organization"
	_ "hunrx.com/attendance_management/modules/payroll"
	"hunrx.com/attendance_management/pkg/app"
	"hunrx.com/attendance_management/pkg/config"
	"hunrx.com/attendance_management/pkg/db"
	"hunrx.com/attendance_management/pkg/events"
)

// Modules is the comma-separated list of modules to activate. Empty runs all.
var modules string

func main() {
	flag.StringVar(&modules, "modules", "", "comma-separated list of modules to run; empty runs all")
	flag.Parse()

	cfg := config.Load()
	if len(modules) > 0 {
		cfg.EnabledModules = split(modules)
	}

	dbPool, err := db.OpenDB(cfg.DBDSN)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer dbPool.Pool.Close()

	ctx := context.Background()
	core := &app.CoreContext{
		Engine:  nil,
		Pool:    dbPool.Pool,
		Events:  events.NewInMemoryBus(),
		Modules: nil,
	}
	if err := app.Setup(ctx, cfg, dbPool.Pool, core); err != nil {
		log.Fatalf("setup: %v", err)
	}

	srv, err := app.New(core)
	if err != nil {
		log.Fatalf("server: %v", err)
	}
	if err := srv.Run(cfg.Host, cfg.Port); err != nil {
		log.Fatalf("server: %v", err)
	}
}

func split(s string) []string {
	var out []string
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}
