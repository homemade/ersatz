package main

import (
	"fmt"
	"io"
	"os"

	docopt "github.com/docopt/docopt.go"
)

const (
	version = "v0.0.1"
	usage   = `usage:
	ersatz start <port> <definitions_dir>
	ersatz -h | --help
	ersatz --version

	options:
		-h --help  		show this screen
		--version  		show version
`
)

func main() {
	os.Exit(entryPoint(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

func entryPoint(cliArgs []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) int {

	args, err := docopt.Parse(usage, cliArgs, true, version, true)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	if args["start"].(bool) {

		stop := make(chan interface{}, 1)

		startApp := NewServerApp(args["<port>"].(string), args["<definitions_dir>"].(string))

		if err := startApp.Setup(); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}

		fmt.Println("[ERSATZ] Listening on port " + startApp.Port)

		startApp.Run(stop)
	}

	return 0
}
