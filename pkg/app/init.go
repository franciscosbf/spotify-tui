package app

import (
	"fmt"
	"os"

	"github.com/franciscosbf/spotify-tui/pkg/ui"
)

func Run() {
	tui := ui.New("./configs/auth.json")

	if err := tui.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "app error: %s", err)
	}
}
