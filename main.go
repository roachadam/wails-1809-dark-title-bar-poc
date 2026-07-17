// Command wailspoc is a minimal Wails v3 app used to verify the Windows Server
// 2019 (build 17763 / 1809) dark-mode title-bar fix.
//
// It creates a single window that requests a Dark native title bar. On an
// unpatched Wails this panics on startup (nil AllowDarkModeForWindow); with the
// fix the window should open, and — if the 1809 legacy path works — its title
// bar should be dark.
//
// Build a Windows binary from macOS/Linux (no CGO needed for Wails v3 Windows):
//
//	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o Ordara-DarkModePOC.exe .
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
  <p>Built against patched Wails &rarr; the window opens (no nil-pointer panic),
     and on Windows 10 1809 / Server 2019 in dark mode the <b>title bar</b> is dark.</p>
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
