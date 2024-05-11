package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/href/online-watch/online"
	cli "github.com/jawher/mow.cli"
)

func App() *cli.Cli {
	app := cli.App(
		"online-watch",
		"Monitors targets for packetloss using TCP and ICMP(v6)")

	app.Spec = strings.Join([]string{
		"[-p|--port=<port>]",
		"[--no-tcp]",
		"[--no-icmp]",
		"[--verbose]",
		"[--interval=<ms>]",
		"[--timeout=<ms>]",
		"TARGETS...",
	}, " ")

	var w online.Watch

	app.StringsArgPtr(&w.Targets, "TARGETS", nil, "[label=]<IPv4|IPv6>[:port]")
	app.IntOptPtr(&w.Port, "p port", 22, "Disable TCP checks")
	app.BoolOptPtr(&w.NoTCP, "no-tcp", false, "Disable TCP checks")
	app.BoolOptPtr(&w.NoICMP, "no-icmp", false, "Disable ICMP checks")
	app.BoolOptPtr(&w.Verbose, "verbose", false, "Log all results")
	app.IntOptPtr(&w.Interval, "interval", 250, "Time between checks (ms)")
	app.IntOptPtr(&w.Timeout, "timeout", 250, "TCP/ICMP timeout (ms)")

	results := make(chan online.Result)

	app.Action = func() {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()

		if err := w.Run(ctx, results); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

		printer := online.Printer{
			Output:              os.Stdout,
			Verbose:             w.Verbose,
			ConsecutiveFailures: 1,
		}

		func() {
			for {
				select {
				case result := <-results:
					printer.Print(result)
				case <-ctx.Done():
					return
				}
			}
		}()

		fmt.Fprintf(os.Stdout, "\n")
		printer.Summary()
	}

	return app
}

func main() {
	if err := App().Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
