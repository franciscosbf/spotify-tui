package app

import (
	"flag"
	"fmt"
	"os"

	"github.com/franciscosbf/spotify-tui/pkg/ui"
)

func die(err any) {
	fmt.Fprintf(os.Stderr, "app error: %s\n", err)
	os.Exit(1)
}

func Run() {
	var config string

	flag.StringVar(&config, "config", "<path>", "configuration file")
	flag.Parse()

	if config == "" {
		die("config parameter missing")
	}

	tui := ui.New("./configs/config.json")

	if err := tui.Start(); err != nil {
		die(err)
	}
}
