package cmd

import (
    "flag"
)

type Args struct {
    Pid       int
    OutputDir string
    Verbose   bool
}

func ParseArgs() Args {
    var args Args

    flag.IntVar(&args.Pid, "pid", 0, "Process id of target application")
    flag.StringVar(&args.OutputDir, "output-dir", "", "Directory to save dumped dex files")
    flag.BoolVar(&args.Verbose, "verbose", false, "Enable verbose output")

    flag.Parse()

    return args
}
