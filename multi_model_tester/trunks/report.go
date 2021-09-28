package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"

	trunks "github.com/straightdave/trunks/lib"
)

func reportCmd() command {
	fs := flag.NewFlagSet("trunks report", flag.ExitOnError)
	reporter := fs.String("reporter", "text", "Reporter [text, json, plot, hist[buckets]]")
	inputs := fs.String("inputs", "stdin", "Input files (comma separated)")
	output := fs.String("output", "stdout", "Output file")
	return command{fs, func(args []string) error {
		fs.Parse(args)
		return report(*reporter, *inputs, *output)
	}}
}

// report validates the report arguments, sets up the required resources
// and writes the report
func report(reporter, inputs, output string) error {
	if len(reporter) < 4 {
		return fmt.Errorf("bad reporter: %s", reporter)
	}

	files := strings.Split(inputs, ",")
	srcs := make([]io.Reader, len(files))
	for i, f := range files {
		in, err := file(f, false)
		if err != nil {
			return err
		}
		defer in.Close()
		srcs[i] = in
	}
	dec := trunks.NewDecoder(srcs...)

	out, err := file(output, true)
	if err != nil {
		return err
	}
	defer out.Close()

	var (
		rep    trunks.Reporter
		report trunks.Report
	)

	switch reporter[:4] {
	case "text":
		var m trunks.Metrics
		rep, report = trunks.NewTextReporter(&m), &m
	case "json":
		var m trunks.Metrics
		rep, report = trunks.NewJSONReporter(&m), &m
	case "plot":
		var rs trunks.Results
		rep, report = trunks.NewPlotReporter("Vegeta Plot", &rs), &rs
	case "hist":
		if len(reporter) < 6 {
			return fmt.Errorf("bad buckets: '%s'", reporter[4:])
		}
		var hist trunks.Histogram
		if err := hist.Buckets.UnmarshalText([]byte(reporter[4:])); err != nil {
			return err
		}
		rep, report = trunks.NewHistogramReporter(&hist), &hist
	default:
		return fmt.Errorf("unknown reporter: %q", reporter)
	}

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)

decode:
	for {
		select {
		case <-sigch:
			break decode
		default:
			var r trunks.Result
			if err = dec.Decode(&r); err != nil {
				if err == io.EOF {
					break decode
				}
				return err
			}
			report.Add(&r)
		}
	}

	if c, ok := report.(trunks.Closer); ok {
		c.Close()
	}

	return rep.Report(out)
}
