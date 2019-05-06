package commands

import (
	"context"
	"os"
	"os/signal"
)

func interruptContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		cancel()
	}()

	return ctx
}
