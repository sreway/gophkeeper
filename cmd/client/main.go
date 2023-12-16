package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/sreway/gophkeeper/cmd/client/commands"
)

var (
	listenSignals = []os.Signal{
		os.Interrupt,
	}
	exitCode int
	exitCH   chan int
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), listenSignals...)
	defer func() {
		stop()
		if exitCode != 0 {
			os.Exit(exitCode)
		}
	}()

	cmd := commands.NewRootCmd(ctx)

	exitCH = make(chan int)

	go func() {
		if err := cmd.ExecuteContext(ctx); err != nil {
			stop()
			exitCH <- 1
		}

		exitCH <- 0
	}()

	exitCode = <-exitCH
}
