# wails-1809-dark-title-bar-poc

A minimal [Wails v3](https://github.com/wailsapp/wails) app that reproduces — and
verifies the fix for — a dark title-bar problem on **Windows 10 1809 / build
17763** (which includes **Windows Server 2019**).

## The problem

On Wails v3 `alpha2.117`, a window that requests a dark native title bar via:

```go
Windows: application.WindowsWindow{ Theme: application.Dark }
```

crashes on startup on Windows 10 1809 (build 17763):

```
panic: runtime error: invalid memory address or nil pointer dereference
  wails/v3/pkg/application.(*windowsWebviewWindow).run
    webview_window_windows.go:580
```

Root cause: Wails loads the undocumented dark-mode `uxtheme` exports only on
build ≥ 18334, but calls `AllowDarkModeForWindow` unguarded — so on 1809 the
proc is `nil` and dereferencing it panics. Separately, even without the crash,
1809 needs a different mechanism than modern Windows to actually darken the
title bar (the DWM `DWMWA_USE_IMMERSIVE_DARK_MODE` attribute is 1903+ only).

## The fix

Lives on the Wails fork branch
[`roachadam/wails@fix/windows-1809-dark-title-bar`](https://github.com/roachadam/wails/tree/fix/windows-1809-dark-title-bar):

1. Nil-guard the `AllowDarkModeForWindow` calls in the window theme setup.
2. Lower the `uxtheme` proc-load gate to 17763 so the app-level dark opt-in runs
   on 1809.
3. Route the title bar per build in `SetTheme`: the DWM immersive-dark attribute
   on 1903+, and the legacy per-window property + a forced non-client repaint on
   1809.

**Result (verified on Windows Server 2019, build 17763, system in dark mode):**
the window opens without a panic and the title bar is dark.

## Build & run

This repo builds against a **local clone of the patched Wails fork** via a
relative `replace` directive (`../wails/v3`). Set it up as a sibling directory:

```sh
# clone the patched fork next to this repo
git clone https://github.com/roachadam/wails.git ../wails
git -C ../wails checkout fix/windows-1809-dark-title-bar

# cross-compile a Windows binary (Wails v3 Windows needs no CGO)
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o poc.exe .
```

Copy `poc.exe` to a Windows 10 1809 / Server 2019 host (with the WebView2
runtime installed) and run it. With the system in dark mode, the title bar
should be dark; without the fix it panics on startup instead.
