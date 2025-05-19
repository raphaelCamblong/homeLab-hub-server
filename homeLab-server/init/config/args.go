package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type Args struct {
	ConfigPath string
}

func ParseArgs() *Args {
	args := &Args{}

	flag.StringVar(&args.ConfigPath, "c", "", "path to config file (default: ./config.toml)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nFlags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s                        # Use default config.toml\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -c config.dev.toml     # Use specific config file\n", os.Args[0])
	}

	flag.Parse()

	if args.ConfigPath == "" {
		wd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error getting working directory: %v\n", err)
			os.Exit(1)
		}
		args.ConfigPath = filepath.Join(wd, "config.toml")
	}

	if _, err := os.Stat(args.ConfigPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "❌ Config file not found at %s\n", args.ConfigPath)
		os.Exit(1)
	}

	return args
}
