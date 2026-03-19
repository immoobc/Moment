# Requirements Document

## Introduction

"此刻 (Moment)" 是一款基于 Go 语言开发的跨平台桌面悬浮时钟工具，支持 Windows 和 macOS。它以简约风格在桌面上悬浮显示时间信息，并提供定时护眼提醒功能。

## Glossary

- **Moment_App**: 此刻桌面时钟应用程序主体
- **Clock_Widget**: 悬浮时钟窗口组件，用于在桌面上显示时间信息
- **Context_Menu**: 右键弹出的上下文菜单，用于访问设置和操作
- **Theme_Engine**: 主题引擎，负责管理和切换应用的视觉主题
- **Rest_Overlay**: 护眼休息蒙版层，定时全屏显示绿色半透明蒙版提醒用户休息
- **Display_Mode**: 时间显示模式，包括数字时间、模拟时钟、时间戳三种形式
- **Window_Level**: 窗口层级设置，控制悬浮窗口在桌面的前后层级
- **Position_Lock**: 位置锁定功能，防止意外拖动改变窗口位置

## Requirements

### Requirement 1: 时间显示

**User Story:** 作为用户，我希望在桌面上看到一个悬浮的时间窗口，以便随时了解当前时间。

#### Acceptance Criteria

1. WHEN the Moment_App starts, THE Clock_Widget SHALL display a floating window on the desktop showing the current time
2. WHEN a user selects a Display_Mode from the Context_Menu, THE Clock_Widget SHALL switch to the selected mode (digital time, analog clock, or Unix timestamp)
3. WHILE the Clock_Widget is visible, THE Clock_Widget SHALL update the displayed time every second
4. THE Clock_Widget SHALL render in a minimal and clean visual style

### Requirement 2: 右键菜单

**User Story:** 作为用户，我希望通过右键菜单快速访问所有设置和操作，以便方便地控制应用行为。

#### Acceptance Criteria

1. WHEN a user right-clicks on the Clock_Widget, THE Moment_App SHALL display a Context_Menu with all available settings and actions
2. WHEN a user selects an option from the Context_Menu, THE Moment_App SHALL execute the corresponding action and close the menu
3. WHEN a user clicks outside the Context_Menu, THE Moment_App SHALL close the menu without performing any action

### Requirement 3: 主题切换

**User Story:** 作为用户，我希望切换应用的主题和背景色，以便适应不同的使用环境和个人偏好。

#### Acceptance Criteria

1. WHEN a user selects "Light" theme from the Context_Menu, THE Theme_Engine SHALL apply a light color scheme to the Clock_Widget
2. WHEN a user selects "Dark" theme from the Context_Menu, THE Theme_Engine SHALL apply a dark color scheme to the Clock_Widget
3. WHEN a user selects "Follow System" theme from the Context_Menu, THE Theme_Engine SHALL detect the operating system theme and apply the matching color scheme
4. WHILE "Follow System" theme is active, THE Theme_Engine SHALL monitor the operating system theme and update the Clock_Widget appearance when the system theme changes

### Requirement 4: 窗口层级

**User Story:** 作为用户，我希望设置悬浮窗口的层级，以便控制时钟窗口在桌面上的前后位置。

#### Acceptance Criteria

1. WHEN a user selects "Always on Top" from the Context_Menu, THE Clock_Widget SHALL remain above all other windows on the desktop
2. WHEN a user selects "Normal Level" from the Context_Menu, THE Clock_Widget SHALL behave as a standard window that can be covered by other windows
3. WHEN the Window_Level setting changes, THE Clock_Widget SHALL apply the new level immediately without restarting

### Requirement 5: 位置锁定

**User Story:** 作为用户，我希望锁定或解锁悬浮窗口的位置，以便防止意外拖动或方便重新定位。

#### Acceptance Criteria

1. WHEN a user activates Position_Lock from the Context_Menu, THE Clock_Widget SHALL prevent any drag movement of the window
2. WHEN a user deactivates Position_Lock from the Context_Menu, THE Clock_Widget SHALL allow the user to drag and reposition the window freely
3. WHILE Position_Lock is inactive, THE Clock_Widget SHALL save the new position when the user finishes dragging

### Requirement 6: 护眼休息提醒

**User Story:** 作为用户，我希望每隔一段时间收到休息提醒，以便保护视力和保持健康的工作习惯。

#### Acceptance Criteria

1. WHEN the configured rest interval elapses (default 45 minutes), THE Rest_Overlay SHALL gradually increase its opacity to display a full-screen green overlay on the desktop
2. WHEN the Rest_Overlay is fully visible, THE Rest_Overlay SHALL remain displayed for 5 minutes
3. WHEN the user moves the mouse or presses a key while the Rest_Overlay is displayed, THE Rest_Overlay SHALL gradually decrease its opacity and disappear
4. WHEN the Rest_Overlay disappears (either after timeout or user interaction), THE Moment_App SHALL reset the rest interval timer
5. WHEN a user adjusts the rest interval from the Context_Menu, THE Moment_App SHALL apply the new interval immediately
6. WHEN a user adjusts the overlay opacity level from the Context_Menu, THE Rest_Overlay SHALL use the configured maximum opacity for future displays

### Requirement 7: 跨平台支持

**User Story:** 作为用户，我希望在 Windows 和 macOS 上都能使用此工具，以便在不同设备间保持一致的体验。

#### Acceptance Criteria

1. THE Moment_App SHALL compile and run on Windows (amd64, arm64)
2. THE Moment_App SHALL compile and run on macOS (amd64, arm64)
3. THE Moment_App SHALL use platform-native APIs for system theme detection on each supported operating system

### Requirement 8: 设置持久化

**User Story:** 作为用户，我希望应用记住我的设置，以便下次启动时无需重新配置。

#### Acceptance Criteria

1. WHEN a user changes any setting, THE Moment_App SHALL persist the setting to a local configuration file
2. WHEN the Moment_App starts, THE Moment_App SHALL load previously saved settings from the configuration file
3. IF the configuration file is missing or corrupted, THEN THE Moment_App SHALL use default settings and create a new configuration file
