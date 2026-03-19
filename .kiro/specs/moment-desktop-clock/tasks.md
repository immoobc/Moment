# Implementation Plan: 此刻 (Moment)

## Overview

基于 Go + Fyne v2 实现跨平台桌面悬浮时钟工具。按照从核心数据层到 UI 层的顺序逐步构建，每个阶段都包含对应的测试任务。

## Tasks

- [x] 1. 项目初始化与基础结构
  - 初始化 Go module (`go mod init moment`)
  - 安装 Fyne v2 依赖 (`fyne.io/fyne/v2`)
  - 安装 gopter 测试依赖 (`github.com/leanovate/gopter`)
  - 创建目录结构: `cmd/moment/`, `core/`, `ui/`
  - 创建 `cmd/moment/main.go` 入口文件，初始化 Fyne App 并显示空窗口
  - _Requirements: 7.1, 7.2_

- [-] 2. 实现配置存储模块
  - [x] 2.1 实现 Config 数据模型和 ConfigStore
    - 在 `core/config.go` 中定义 Config 结构体、DefaultConfig()、ConfigStore 类型
    - 实现 Load()、Save()、Get()、Update() 方法
    - 使用 `encoding/json` 进行序列化
    - 配置文件路径: Windows `%APPDATA%/Moment/config.json`，macOS `~/Library/Application Support/Moment/config.json`
    - _Requirements: 8.1, 8.2, 8.3_

  - [ ]* 2.2 编写 Config 序列化 round-trip 属性测试
    - **Property 8: Config serialization round-trip**
    - **Validates: Requirements 8.1, 8.2**

  - [ ]* 2.3 编写 Config 错误处理属性测试
    - **Property 9: Corrupted config falls back to defaults**
    - **Validates: Requirements 8.3**

- [x] 3. Checkpoint - 确保配置模块测试通过
  - Ensure all tests pass, ask the user if questions arise.

- [-] 4. 实现主题管理模块
  - [x] 4.1 实现 ThemeManager 和自定义 MomentTheme
    - 在 `core/theme.go` 中定义 ThemeMode 枚举、ThemeManager 结构体
    - 实现 MomentTheme（实现 fyne.Theme 接口），支持 Light/Dark 配色
    - 实现 SetMode()、GetMode()、CurrentTheme() 方法
    - 实现系统主题检测（跟随系统模式）
    - _Requirements: 3.1, 3.2, 3.3, 3.4_

  - [ ]* 4.2 编写主题模式状态一致性属性测试
    - **Property 2: Theme mode state consistency**
    - **Validates: Requirements 3.1, 3.2**

- [-] 5. 实现窗口管理模块
  - [x] 5.1 实现 WindowManager
    - 在 `core/window.go` 中定义 WindowLevel 枚举、WindowManager 结构体
    - 实现 SetLevel()、SetLocked()、GetPosition() 方法
    - 实现平台特定的窗口置顶逻辑（Windows: `user32.dll SetWindowPos`，macOS: `NSWindow.setLevel`）
    - _Requirements: 4.1, 4.2, 4.3, 5.1, 5.2, 5.3_

  - [ ]* 5.2 编写锁定状态一致性属性测试
    - **Property 3: Lock state consistency**
    - **Validates: Requirements 5.1, 5.2**

  - [ ]* 5.3 编写位置持久化属性测试
    - **Property 4: Position persistence while unlocked**
    - **Validates: Requirements 5.3**

- [x] 6. Checkpoint - 确保核心模块测试通过
  - Ensure all tests pass, ask the user if questions arise.

- [x] 7. 实现时钟组件
  - [x] 7.1 实现 ClockWidget 和数字/时间戳显示模式
    - 在 `ui/clock.go` 中定义 DisplayMode 枚举、ClockWidget 结构体
    - 实现 NewClockWidget()、SetMode()、CreateRenderer()
    - 实现数字时间格式化 ("15:04:05") 和 Unix 时间戳格式化
    - 使用 time.Ticker 每秒刷新
    - _Requirements: 1.1, 1.2, 1.3_

  - [x] 7.2 实现模拟时钟渲染器
    - 在 `ui/analog.go` 中实现 AnalogClockRenderer
    - 使用 Fyne Canvas 绘制表盘圆形、12 个刻度标记、时/分/秒针
    - 根据当前时间计算指针角度
    - _Requirements: 1.2_

  - [ ]* 7.3 编写显示模式状态一致性属性测试
    - **Property 1: Display mode state consistency**
    - **Validates: Requirements 1.2**

- [x] 8. 实现休息提醒模块
  - [x] 8.1 实现 RestTimer
    - 在 `core/rest.go` 中定义 RestTimer 结构体
    - 实现 SetInterval()、SetMaxOpacity()、Start()、Stop()、Reset() 方法
    - 使用 time.Timer 实现定时触发
    - _Requirements: 6.1, 6.4, 6.5_

  - [x] 8.2 实现 RestOverlay 全屏蒙版
    - 在 `ui/overlay.go` 中定义 RestOverlay 结构体
    - 创建全屏无边框绿色窗口
    - 实现渐入渐出动画（通过定时调整背景色 alpha 值）
    - 监听鼠标移动和键盘事件触发消失
    - _Requirements: 6.1, 6.2, 6.3, 6.6_

  - [ ]* 8.3 编写休息计时器属性测试
    - **Property 5: Rest timer triggers at configured interval**
    - **Property 6: Timer reset after overlay dismiss**
    - **Property 7: Rest settings applied immediately**
    - **Validates: Requirements 6.1, 6.4, 6.5, 6.6**

- [x] 9. Checkpoint - 确保所有组件测试通过
  - Ensure all tests pass, ask the user if questions arise.

- [x] 10. 实现右键菜单与应用集成
  - [x] 10.1 实现 ContextMenu
    - 在 `ui/menu.go` 中定义 ContextMenu 结构体
    - 构建完整菜单树：显示模式、主题、窗口层级、位置锁定、休息提醒设置、退出
    - 绑定各菜单项到对应的 Core 模块方法
    - _Requirements: 2.1, 2.2, 2.3_

  - [x] 10.2 实现 MomentApp 控制器并集成所有组件
    - 在 `app.go` 中定义 MomentApp 结构体
    - 初始化所有组件：ConfigStore → ThemeManager → WindowManager → ClockWidget → RestTimer → RestOverlay → ContextMenu
    - 设置无边框悬浮窗口（Fyne Splash Window）
    - 绑定右键事件到 ContextMenu
    - 实现窗口拖动逻辑（受 Position_Lock 控制）
    - 更新 `cmd/moment/main.go` 使用 MomentApp
    - _Requirements: 1.1, 1.4, 2.1, 4.3, 5.3_

- [x] 11. 最终 Checkpoint - 确保所有测试通过
  - Ensure all tests pass, ask the user if questions arise.

## Notes

- 标记 `*` 的任务为可选任务，可跳过以加快 MVP 开发
- 每个任务引用了具体的需求编号以保证可追溯性
- Checkpoint 任务用于阶段性验证
- Property tests 使用 gopter 库，每个 test 至少运行 100 次迭代
- 平台特定代码使用 Go build tags (`//go:build windows` / `//go:build darwin`) 隔离
