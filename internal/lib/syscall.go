package lib

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

// create channel for exit signal
// wait for signal and canlcel global context
func SysCallProcess(
	ctx context.Context,
	cancel context.CancelFunc,
	l *slog.Logger,
	ff ...func(),
) {
	defer cancel()
	exit := make(chan os.Signal, 1)
	signal.Notify(exit,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
		syscall.SIGQUIT,
	)
	select {
	case s := <-exit:
		l.Info("received signal: ", "syscal", s.Signal)
	case <-ctx.Done():
		l.Info("base context terminated")
	}

	l.Info("prepare for shutdown")
	if len(ff) > 0 {
		l.Info("performing pre-shutdown tasks")
		for _, f := range ff {
			f()
		}
	}
	l.Info("shutdown procedure complete")
}
