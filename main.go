// Command wailspoc is a minimal Wails v3 app that reproduces a startup panic on
// Windows 10 1809 (build 17763), including Windows Server 2019.
//
// It creates a single window that requests a Dark native title bar. Built
// against the official tagged Wails release, running it on 1809 panics on
// startup with a nil pointer dereference (nil AllowDarkModeForWindow), before
// the window is shown. Built against the patched fork branch, the window opens
// instead — and with the dark-title-bar branch its title bar is dark when the
// system is in dark mode. See README.md.
//
// Build a Windows binary from macOS/Linux (Wails v3 Windows needs no CGO):
//
//	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o poc.exe .
package main

import (
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

const page = `<!doctype html>
<html>
<head>
<meta charset="utf-8">
<style>
  html, body { margin: 0; height: 100%; font-family: "Segoe UI", system-ui, sans-serif; }
  body {
    background: #1b2636; color: #e6edf3;
    display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 14px;
  }
  h1 { font-weight: 600; margin: 0; }
  p  { margin: 0; opacity: .9; }
  code { background: #0d1117; padding: 2px 8px; border-radius: 6px; }
</style>
</head>
<body>
  <h1>Wails 1809 dark title-bar POC</h1>
  <p>Requested native theme: <code>application.Dark</code></p>
  <p>You are seeing this window, so Wails did <b>not</b> panic &mdash; you built
     against the patched fork, or you are on a build &ge; 18334.</p>
  <p>On Windows 10 1809 / Server 2019 with the official release, this panics
     before the window appears.</p>
</body>
</html>`

func main() {
	app := application.New(application.Options{
		Name:        "DarkModePOC",
		Description: "Windows Server 2019 dark-mode title-bar POC",
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Ordara Dark-Mode POC",
		Width:  900,
		Height: 600,
		Windows: application.WindowsWindow{
			Theme: application.Dark,
		},
		HTML: page,
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
