# wails-1809-dark-title-bar-poc

A minimal [Wails v3](https://github.com/wailsapp/wails) app that reproduces a
startup **panic** on **Windows 10 1809 / build 17763** (which includes **Windows
Server 2019**), and can also be used to verify the fix.

## The bug

A window that requests a dark native title bar via:

```go
Windows: application.WindowsWindow{ Theme: application.Dark }
```

crashes on startup on Windows 10 1809 (build 17763), before the window is shown:

```
panic: runtime error: invalid memory address or nil pointer dereference
  wails/v3/pkg/application.(*windowsWebviewWindow).run
    webview_window_windows.go:580
```

Root cause: Wails loads the undocumented dark-mode `uxtheme` exports (incl.
`AllowDarkModeForWindow`, ordinal 133) only on build ≥ 18334, but calls
`AllowDarkModeForWindow` unguarded — so on 1809 the proc is `nil` and
dereferencing it panics. Upstream issue: wailsapp/wails#5792.

## Reproduce the panic

This repo builds against the **official tagged release**
(`v3.0.0-alpha2.117`), so no extra setup is needed:

```sh
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o poc.exe .
```

Copy `poc.exe` to a Windows 10 1809 / Windows Server 2019 host (with the
WebView2 runtime installed) and run it. It panics on startup and no window
appears.

## Verify the fix

The fix lives on the Wails fork
[`roachadam/wails`](https://github.com/roachadam/wails). Point this repo at it
with a local `replace`:

```sh
# clone the patched fork next to this repo
git clone https://github.com/roachadam/wails.git ../wails

# pick the branch to verify:
#   fix/windows-1809-dark-mode-panic   -> crash fix only (window opens; light title bar on 1809)
#   fix/windows-1809-dark-title-bar    -> crash fix + dark title bar on 1809
git -C ../wails checkout fix/windows-1809-dark-mode-panic

# build against the fork
go mod edit -replace github.com/wailsapp/wails/v3=../wails/v3
go mod tidy
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o poc.exe .

# revert when done
go mod edit -dropreplace github.com/wailsapp/wails/v3
go mod tidy
```

Run `poc.exe` on the same 1809 host: with the crash-fix branch the window opens
(no panic); with the dark-title-bar branch the title bar is dark when the system
is in dark mode.
