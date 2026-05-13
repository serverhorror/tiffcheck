package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime/debug"

	"github.com/serverhorror/tiffcheck/pkg/tiffcheck"
)

var check = tiffcheck.Check

const usageText = `Usage: tiffcheck [--help] [--version] <filename>`

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, out io.Writer, errOut io.Writer) int {
	flags := flag.NewFlagSet("tiffcheck", flag.ContinueOnError)
	flags.SetOutput(io.Discard)

	help := flags.Bool("help", false, "show help")
	version := flags.Bool("version", false, "show version")

	if err := flags.Parse(args); err != nil {
		fmt.Fprintf(errOut, "Error: %v\n", err)
		fmt.Fprintf(errOut, "%s\n", usageText)
		return 2
	}

	if *help {
		fmt.Fprintf(out, "%s\n", usageText)
		return 0
	}
	if *version {
		printVersion(out)
		return 0
	}

	remaining := flags.Args()
	if len(remaining) != 1 {
		fmt.Fprintf(errOut, "%s\n", usageText)
		return 2
	}

	file, err := os.Open(remaining[0])
	if err != nil {
		fmt.Fprintf(errOut, "Error: %v\n", err)
		return 1
	}
	defer file.Close()

	report, err := check(file)
	if err != nil {
		fmt.Fprintf(errOut, "Error: %v\n", err)
		return 1
	}

	printReport(out, report)
	return 0
}

func printVersion(out io.Writer) {
	revision := "unknown"
	timestamp := "unknown"
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" && setting.Value != "" {
				revision = setting.Value
			}
			if setting.Key == "vcs.time" && setting.Value != "" {
				timestamp = setting.Value
			}
		}
	}

	fmt.Fprintf(out, "tiffcheck\nrevision: %s\ntime: %s\n", revision, timestamp)
}

func printReport(out io.Writer, report tiffcheck.Report) {
	fmt.Fprintf(out, "--- TIFF Structural Analysis ---\n")
	if report.Tiled {
		fmt.Fprintf(out, "Status:  TILED (Indexed for partial loading)\n")
		fmt.Fprintf(out, "Tile Dim: %v pixels wide\n", report.TileWidth)
	} else {
		fmt.Fprintf(out, "Status:  STRIPPED (Standard linear loading)\n")
		fmt.Fprintf(out, "Note:    Not optimized for partial random access.\n")
	}

	if report.GeoTIFF {
		fmt.Fprintf(out, "Metadata: GeoTIFF tags detected\n")
	}
}
