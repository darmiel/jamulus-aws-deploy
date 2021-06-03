// +build windows

package thin

import (
	"fmt"
	"golang.org/x/sys/windows"
	"log"
	"os"
)

func report(err error) {
	if err == nil {
		return
	}
	log.Println("Patch | WARN:", err)
}

func init() {
	fmt.Println("Applying windows terminal patch ...")
	stdout := windows.Handle(os.Stdout.Fd())
	var originalMode uint32
	report(windows.GetConsoleMode(stdout, &originalMode))
	report(windows.SetConsoleMode(stdout, originalMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING))
}
