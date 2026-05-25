package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/334456777/weflow-api/cmd"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	cmd.Execute(ctx)
}
