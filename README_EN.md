# Moment (此刻)

A lightweight, elegant desktop floating clock application built with Go and [Fyne](https://fyne.io).

![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go&logoColor=white)
![Fyne](https://img.shields.io/badge/Fyne-v2.7-5C2D91)
![Platform](https://img.shields.io/badge/Platform-Windows-0078D6?logo=windows)

## Features

- 🕐 Real-time floating clock displaying date, weekday, and HH:MM:SS
- 🎨 Dual themes: Calendar White (Apple Calendar style) / Dark Night
- 📌 Always-on-top / normal window level toggle
- 🔒 Lock window position to prevent accidental dragging
- 🖱️ Drag to reposition with automatic position persistence
- 📋 Right-click context menu on the clock + system tray menu
- 🚫 Single instance enforcement — re-launching brings the existing window to front
- 🪟 Borderless window with native rounded corners on Windows 11
- ⌨️ Press Esc to quit

## Quick Start

### Prerequisites

- Go 1.25+
- GCC ([TDM-GCC](https://jmeubank.github.io/tdm-gcc/) or w64devkit recommended)
- Windows 10/11

### Build

```bash
# Using the build script
build.bat

# Or build manually
set CGO_ENABLED=1 && go build -ldflags="-s -w -H windowsgui" -o bin\moment.exe .\cmd\moment
```

### Run

```bash
bin\moment.exe
```

## Usage

| Action | How |
|---|---|
| Move window | Left-click and drag the floating clock |
| Open menu | Right-click the clock, or right-click the system tray icon |
| Switch theme | Menu → Theme → Calendar White / Dark Night |
| Always on top | Menu → Window → Always on Top |
| Lock position | Menu → Window → Lock Position |
| Quit | Menu → Quit, or press Esc |

## Project Structure

```
Moment/
├── cmd/moment/        # App entry point, drag handling, main controller
├── core/              # Config, window management, single instance, platform code
├── ui/                # Clock renderer, menus
├── assets/            # App icon
├── build.bat          # Build script
└── FyneApp.toml       # Fyne app metadata
```

## Configuration

Settings are automatically saved to `%APPDATA%\Moment\config.json`:

```json
{
    "window_level": 0,
    "position_x": 100,
    "position_y": 100,
    "locked": false,
    "theme": 0
}
```

| Field | Values |
|---|---|
| `window_level` | `0` = always on top, `1` = normal |
| `locked` | `true` = position locked |
| `theme` | `0` = Calendar White, `1` = Dark Night |

## License

MIT
