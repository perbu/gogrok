package main

import (
	"context"
	"fmt"
	"github.com/perbu/gogrok/analytics"
	"github.com/perbu/gogrok/render"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"runtime"

	"github.com/dustin/go-humanize"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	err := run(ctx, os.Stdout, os.Environ(), os.Args)
	if err != nil {
		fmt.Println("run() returned an error: ", err.Error())
		os.Exit(1)
	}
}

func run(ctx context.Context, output io.Writer, env, args []string) error {
	r, err := analytics.New("code")
	if err != nil {
		return fmt.Errorf("analytics.New: %w", err)
	}
	err = r.Parse()
	if err != nil {
		return fmt.Errorf("r.Parse: %w", err)
	}
	// log memory usage so far:
	logMemoryUsage()

	s, err := render.New(r)
	if err != nil {
		return fmt.Errorf("render.New: %w", err)
	}
	err = s.Start(ctx)
	if err != nil {
		return fmt.Errorf("s.Start: %w", err)
	}

	return nil
}

func logMemoryUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	slog.Info("Heap Stats",
		"HeapAlloc", humanize.Bytes(m.HeapAlloc),
		"HeapSys", humanize.Bytes(m.HeapSys),
		"HeapIdle", humanize.Bytes(m.HeapIdle),
		"HeapInuse", humanize.Bytes(m.HeapInuse),
		"HeapReleased", humanize.Bytes(m.HeapReleased),
		"HeapObjects", humanize.Comma(int64(m.HeapObjects)),
	)
}
