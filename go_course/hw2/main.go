package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	From := flag.Int("From", 0, "an int")
	To := flag.Int("To", 0, "an int")
	Workers := flag.Int("Workers", 10, "an int")
	Timeout := flag.Duration("Timeout", 5*time.Second, "duration")
	flag.Parse()

	if *From <= 0 || *To <= 0 || *From > *To {
		fmt.Fprintln(os.Stderr, "Please type the flags -From=N and -To=M.From should be < To")
		os.Exit(1)
	}

	if *Workers <= 0 {
		fmt.Fprintln(os.Stderr, "Workers should be > 0")
		os.Exit(1)
	}
	if *Timeout <= 0 {
		fmt.Fprintln(os.Stderr, "Timeout should be >0")
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(
		context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	wp := WorkerPool{
		concurrency: *Workers,
		timeout:     *Timeout,
	}
	wp.Run(ctx, *From, *To)
	if ctx.Err() != nil {
		fmt.Println()
		fmt.Println(context.Cause(ctx))
		fmt.Println("exiting")
	}
}
