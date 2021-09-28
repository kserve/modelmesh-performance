package main

import (
	"flag"
	"fmt"
	"io"
	"strings"

	trunks "github.com/straightdave/trunks/lib"
)

func dumpCmd() command {
	fs := flag.NewFlagSet("trunks dump", flag.ExitOnError)
	dumper := fs.String("dumper", "json", "Dumper [json, csv]")
	inputs := fs.String("inputs", "stdin", "Input files (comma separated)")
	output := fs.String("output", "stdout", "Output file")
	return command{fs, func(args []string) error {
		fs.Parse(args)
		return dump(*dumper, *inputs, *output)
	}}
}

func dump(dumper, inputs, output string) error {
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

	var enc trunks.Encoder
	switch dumper {
	case "csv":
		enc = trunks.NewCSVEncoder(out)
	case "json":
		enc = trunks.NewJSONEncoder(out)
	default:
		return fmt.Errorf("unsupported dumper: %s", dumper)
	}

	for {
		var r trunks.Result
		if err = dec.Decode(&r); err != nil {
			if err == io.EOF {
				break
			}
			return err
		} else if err = enc.Encode(&r); err != nil {
			return err
		}
	}

	return nil
}
