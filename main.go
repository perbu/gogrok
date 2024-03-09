package main

import (
	"context"
	"fmt"
	"github.com/perbu/gogrok/analytics"
	"io"
	"os"
	"os/signal"
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
	return nil
}
