# 此刻 Moment

一款轻量、美观的桌面悬浮时钟应用，使用 Go + [Fyne](https://fyne.io) 构建。

![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go&logoColor=white)
![Fyne](https://img.shields.io/badge/Fyne-v2.7-5C2D91)
![Platform](https://img.shields.io/badge/Platform-Windows-0078D6?logo=windows)

## 功能特性

- 🕐 实时时钟悬浮窗，显示日期、星期、时分秒
- 🎨 双主题切换：日历白（Apple Calendar 风格）/ 暗夜黑
- 📌 窗口置顶 / 普通层级切换
- 🔒 锁定窗口位置，防止误拖动
- 🖱️ 悬浮窗支持拖动定位，位置自动记忆
- 📋 悬浮窗右键菜单 + 系统托盘菜单，操作一致
- 🚫 单实例运行，重复启动自动唤起已有窗口
- 🪟 无边框窗口，Windows 11 原生圆角支持
- ⌨️ 按 Esc 快速退出

## 截图

| 日历白 | 暗夜黑 |
|:---:|:---:|
| 白色背景 + 红色日期头 | 深色背景 + 暗红日期头 |

## 快速开始

### 环境要求

- Go 1.25+
- GCC（推荐 [TDM-GCC](https://jmeubank.github.io/tdm-gcc/) 或 w64devkit）
- Windows 10/11

### 编译

```bash
# 使用 build.bat
build.bat

# 或手动编译
set CGO_ENABLED=1 && go build -ldflags="-s -w -H windowsgui" -o bin\moment.exe .\cmd\moment
```

### 运行

```bash
bin\moment.exe
```

## 使用说明

| 操作 | 方式 |
|---|---|
| 移动窗口 | 鼠标左键拖动悬浮窗 |
| 打开菜单 | 右键点击悬浮窗，或右键系统托盘图标 |
| 切换主题 | 菜单 → 主题 → 日历白 / 暗夜黑 |
| 窗口置顶 | 菜单 → 窗口 → 置顶 |
| 锁定位置 | 菜单 → 窗口 → 锁定位置 |
| 退出 | 菜单 → 退出，或按 Esc |

## 项目结构

```
Moment/
├── cmd/moment/        # 应用入口、窗口拖动、主控逻辑
├── core/              # 配置管理、窗口管理、单实例、平台适配
├── ui/                # 时钟渲染、菜单
├── assets/            # 应用图标
├── build.bat          # 编译脚本
└── FyneApp.toml       # Fyne 应用元数据
```

## 配置文件

配置自动保存在 `%APPDATA%\Moment\config.json`，包含：

```json
{
    "window_level": 0,
    "position_x": 100,
    "position_y": 100,
    "locked": false,
    "theme": 0
}
```

## 许可证

MIT

---
