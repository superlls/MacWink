# MacWink 🖥️ ↔ 💻

一个轻量级的局域网剪贴板同步工具，采用 **Peer 模式**：每个实例启动后既是 Server（监听 `9999`）也是 Watcher（监控本地剪贴板），可在 macOS/Windows/Linux 间双向同步文本。

**English** | [中文](#中文文档)

---

## 💡 Core Workflow

MacWink 通过 Peer 模式实现跨平台剪贴板同步：

```
┌──────────────────┐         TCP Network         ┌──────────────────┐
│      Peer A      │  <────────────────────────> │      Peer B      │
│ (Server+Watcher) │      Clipboard Frames       │ (Server+Watcher) │
│ • Listen :9999   │                             │ • Listen :9999   │
│ • Watch clipboard│                             │ • Watch clipboard│
└──────────────────┘                             └──────────────────┘
```

**工作原理：**
1. 每个实例同时监听 `:9999` 并轮询本地剪贴板变化
2. 本地剪贴板变化后，向 `-peer` 指定的对端发送数据
3. 收到网络内容后写入剪贴板，并更新 `tracking` 以防止回环发送
4. 对端重启/断开后自动重连

---

## 🛠️ Prerequisites (环境准备)

### 系统要求

| 组件 | 要求 | 备注 |
|------|------|------|
| **Go SDK** | 1.25+ | [下载地址](https://golang.org/dl/) |
| **网络** | 同一局域网 | Mac 和 Windows 需在同一网络 |
| **防火墙** | 允许 9999 端口 | 需要配置 Windows 防火墙 |

### 安装 Go SDK

**macOS:**
```bash
# 使用 Homebrew
brew install go

# 验证安装
go version
```

**Windows:**
1. 访问 [golang.org/dl](https://golang.org/dl/)
2. 下载 Windows 版本 (.msi)
3. 双击安装，按默认选项完成
4. 重启电脑后验证：
```powershell
go version
```

### 国内用户加速配置 🚀

如果在国内环境下依赖下载缓慢，配置 Go Proxy：

**macOS/Linux:**
```bash
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
```

**Windows (PowerShell):**
```powershell
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
```

验证配置：
```bash
go env GOPROXY
# 应该输出: https://goproxy.cn,direct
```

---

## 📦 Dependency Installation (依赖安装)

项目依赖 `github.com/atotto/clipboard` 用于跨平台剪贴板访问。

**下载依赖：**

```bash
# 进入项目目录
cd MacWink

# 下载并整理依赖
go mod tidy
```

如果遇到网络问题，可以尝试：
```bash
go get -u github.com/atotto/clipboard@v0.1.4
```

---

## 🚀 Quick Start Guide (快速启动)

### 前置检查

#### 第一步：检查网络连接 🌐

在 **Mac** 上 ping Windows 的 IP 地址，确保网络通畅：

```bash
# 替换 192.168.1.100 为你 Windows 的实际 IP
ping 192.168.1.100
```

如果能收到回复，说明网络连接正常。

**如何查找 Windows IP？**

在 Windows PowerShell 中运行：
```powershell
ipconfig
```

找到 `IPv4 Address` 一行，例如 `192.168.1.100`

---

#### 第二步：配置 Windows 防火墙 🔥

**方法 A：使用 PowerShell (推荐)**

以管理员身份打开 PowerShell，运行：

```powershell
New-NetFirewallRule -DisplayName "MacWink-TCP-In" -Direction Inbound -Action Allow -Protocol TCP -LocalPort 9999
```

**方法 B：手动配置**

1. 打开 **Windows Defender 防火墙**
2. 点击 **允许应用通过防火墙**
3. 点击 **更改设置** (需要管理员权限)
4. 点击 **允许其他应用**
5. 选择 `go.exe` 或 `cmd.exe`，点击 **添加**
6. 确保勾选了 **专用** 和 **公用** 两个选项

---

#### 第三步：启动 Peer (Windows) 🖥️

在 **Windows** 上打开 PowerShell 或 CMD，进入项目目录并启动（`-peer` 传对端 IP）：

```powershell
cd C:\Users\YourUsername\MacWink
go run .\main.go -peer 192.168.1.101
```

---

#### 第四步：启动 Peer (macOS/Linux) 💻

在 **macOS/Linux** 上打开终端，进入项目目录并启动（`-peer` 传对端 IP）：

```bash
cd /Users/YourUsername/GolandProjects/MacWink
go run ./main.go -peer 192.168.1.100
```

---

### 测试同步 ✅

现在尝试在 **Mac** 上复制一些文本：

```bash
# Mac 上复制文本
echo "Hello from MacWink!" | pbcopy
```

你应该看到：

**Mac 端日志：**
```
2026/01/11 20:30:15 clipboard changed: 19 bytes (queued)
```

**Windows 端日志：**
```
2026/01/11 20:30:15 received 19 bytes from 192.168.1.X:XXXXX
```

然后在 Windows 上粘贴（Ctrl+V），就能看到同步过来的内容！🎉

---

## ⚠️ Troubleshooting (常见问题排查)

### 问题 1：i/o timeout 错误

**错误信息：**
```
Failed to sync: connection error: dial tcp 192.168.1.8:9999: i/o timeout
```

**可能原因和解决方案：**

| 原因 | 解决方案 |
|------|--------|
| **IP 地址错误** | 在 Windows 上运行 `ipconfig` 确认真实 IP，重新启动 Sender |
| **防火墙阻止** | 运行防火墙配置命令（见第二步），或检查防火墙设置 |
| **网络不通** | 在 Mac 上 ping Windows IP，确保网络连接正常 |
| **Receiver 未启动** | 确保 Windows 上的 Receiver 程序正在运行 |
| **网络模式** | 确保 Windows 网络设置为"专用"而非"公用" |

**快速诊断：**
```bash
# Mac 上测试连接
nc -zv 192.168.1.100 9999

# 如果显示 "succeeded" 说明连接正常
```

---

### 问题 2：依赖下载失败

**错误信息：**
```
go: github.com/atotto/clipboard@v0.1.4: Get "https://proxy.golang.org/...":
context deadline exceeded
```

**解决方案：**

1. **配置国内代理（推荐）**
```bash
go env -w GOPROXY=https://goproxy.cn,direct
go mod tidy
```

2. **如果在公司网络，配置代理**
```bash
go env -w GOPROXY=http://your-proxy:port
```

3. **离线方案：在有网络的 Mac 上下载**
```bash
go mod download
# 然后将整个项目文件夹复制到 Windows
```

---

### 问题 3：Windows 上找不到 go 命令

**解决方案：**

1. 确认 Go 已安装：在 PowerShell 中运行 `go version`
2. 如果找不到命令，重启 PowerShell 或电脑
3. 检查 Go 安装路径是否在系统 PATH 中：
```powershell
$env:Path -split ';' | Select-String 'Go'
```

---

### 问题 4：Receiver 启动后立即退出

**可能原因：**
- 端口 9999 已被占用
- 权限不足

**解决方案：**
```powershell
# 检查端口占用
netstat -ano | findstr :9999

# 使用其他端口
go run .\cmd\receiver\main.go -port 8888

# 对应地，Mac 上也要改端口
go run ./cmd/sender/main.go -ip 192.168.1.100 -port 8888
```

---

## 📝 Disclaimer (免责声明)

⚠️ **重要提示：**

- **仅限局域网使用** - 本工具设计用于同一局域网内的设备间通信
- **隐私安全** - 剪贴板内容以明文形式在网络上传输，请勿在公网或不信任的网络中使用
- **数据安全** - 不建议在传输敏感信息（密码、密钥等）时使用
- **网络安全** - 确保只在受信任的网络环境中运行此工具

**建议：**
- 仅在家庭网络或公司内网中使用
- 定期检查防火墙规则
- 不需要时及时关闭程序

---

## 📄 License

MIT License - 详见 LICENSE 文件

---

## 🤝 Contributing

欢迎提交 Issue 和 Pull Request！

---

<div id="中文文档"></div>

# MacWink 🖥️ → 💻 (中文文档)

一个轻量级的局域网剪贴板同步工具，支持从 macOS/Linux 实时同步文本到 Windows。

---

## 💡 核心流程

MacWink 通过以下流程实现跨平台剪贴板同步：

```
┌─────────────────┐         TCP 网络         ┌──────────────────┐
│   Mac/Linux     │  ─────────────────────────> │    Windows       │
│   (发送端)      │  剪贴板内容 (9999 端口)    │   (接收端)       │
│                 │  <─────────────────────────  │                  │
│ • 每 500ms 检查 │  连接状态                   │ • 监听 9999 端口 │
│ • 检测变化      │                             │ • 写入剪贴板     │
│ • TCP 发送      │                             │                  │
└─────────────────┘                             └──────────────────┘
```

**工作原理：**
1. **发送端 (Mac/Linux)** 每 500ms 检查一次剪贴板
2. 检测到内容变化时，通过 TCP 连接发送到接收端
3. **接收端 (Windows)** 接收数据并自动写入系统剪贴板
4. 支持自动重连和错误恢复

---

## 🛠️ 环境准备

### 系统要求

| 组件 | 要求 | 备注 |
|------|------|------|
| **Go SDK** | 1.25+ | [下载地址](https://golang.org/dl/) |
| **网络** | 同一局域网 | Mac 和 Windows 需在同一网络 |
| **防火墙** | 允许 9999 端口 | 需要配置 Windows 防火墙 |

### 安装 Go SDK

**macOS:**
```bash
# 使用 Homebrew
brew install go

# 验证安装
go version
```

**Windows:**
1. 访问 [golang.org/dl](https://golang.org/dl/)
2. 下载 Windows 版本 (.msi)
3. 双击安装，按默认选项完成
4. 重启电脑后验证：
```powershell
go version
```

### 国内用户加速配置 🚀

如果在国内环境下依赖下载缓慢，配置 Go Proxy：

**macOS/Linux:**
```bash
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
```

**Windows (PowerShell):**
```powershell
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
```

验证配置：
```bash
go env GOPROXY
# 应该输出: https://goproxy.cn,direct
```

---

## 📦 依赖安装

项目依赖 `github.com/atotto/clipboard` 用于跨平台剪贴板访问。

**下载依赖：**

```bash
# 进入项目目录
cd MacWink

# 下载并整理依赖
go mod tidy
```

如果遇到网络问题，可以尝试：
```bash
go get -u github.com/atotto/clipboard@v0.1.4
```

---

## 🚀 快速启动指南

### 前置检查

#### 第一步：检查网络连接 🌐

在 **Mac** 上 ping Windows 的 IP 地址，确保网络通畅：

```bash
# 替换 192.168.1.100 为你 Windows 的实际 IP
ping 192.168.1.100
```

如果能收到回复，说明网络连接正常。

**如何查找 Windows IP？**

在 Windows PowerShell 中运行：
```powershell
ipconfig
```

找到 `IPv4 Address` 一行，例如 `192.168.1.100`

---

#### 第二步：配置 Windows 防火墙 🔥

**方法 A：使用 PowerShell (推荐)**

以管理员身份打开 PowerShell，运行：

```powershell
New-NetFirewallRule -DisplayName "MacWink-TCP-In" -Direction Inbound -Action Allow -Protocol TCP -LocalPort 9999
```

**方法 B：手动配置**

1. 打开 **Windows Defender 防火墙**
2. 点击 **允许应用通过防火墙**
3. 点击 **更改设置** (需要管理员权限)
4. 点击 **允许其他应用**
5. 选择 `go.exe` 或 `cmd.exe`，点击 **添加**
6. 确保勾选了 **专用** 和 **公用** 两个选项

---

#### 第三步：启动接收端 (Windows) 📥

在 **Windows** 上打开 PowerShell 或 CMD，进入项目目录：

```powershell
cd C:\Users\你的用户名\MacWink
go run .\cmd\receiver\main.go
```

你应该看到类似的输出：
```
2026/01/11 20:30:00 Receiver started. Listening on 0.0.0.0:9999...
2026/01/11 20:30:00 Waiting for clipboard data from Sender...
```

**保持这个窗口打开！** ⚠️

---

#### 第四步：启动发送端 (Mac/Linux) 📤

在 **Mac** 上打开终端，进入项目目录：

```bash
cd /Users/你的用户名/GolandProjects/MacWink

# 替换 192.168.1.100 为你 Windows 的实际 IP
go run ./cmd/sender/main.go -ip 192.168.1.100
```

你应该看到类似的输出：
```
2026/01/11 20:30:05 Sender started. Monitoring clipboard...
2026/01/11 20:30:05 Target Receiver: 192.168.1.100:9999
```

---

### 测试同步 ✅

现在尝试在 **Mac** 上复制一些文本：

```bash
# Mac 上复制文本
echo "Hello from MacWink!" | pbcopy
```

你应该看到：

**Mac 端日志：**
```
2026/01/11 20:30:15 Clipboard changed. Length: 19. Sending...
```

**Windows 端日志：**
```
2026/01/11 20:30:15 Received 19 bytes from 192.168.1.X
```

然后在 Windows 上粘贴（Ctrl+V），就能看到同步过来的内容！🎉

---

## ⚠️ 常见问题排查

### 问题 1：i/o timeout 错误

**错误信息：**
```
Failed to sync: connection error: dial tcp 192.168.1.8:9999: i/o timeout
```

**可能原因和解决方案：**

| 原因 | 解决方案 |
|------|--------|
| **IP 地址错误** | 在 Windows 上运行 `ipconfig` 确认真实 IP，重新启动发送端 |
| **防火墙阻止** | 运行防火墙配置命令（见第二步），或检查防火墙设置 |
| **网络不通** | 在 Mac 上 ping Windows IP，确保网络连接正常 |
| **接收端未启动** | 确保 Windows 上的接收端程序正在运行 |
| **网络模式** | 确保 Windows 网络设置为"专用"而非"公用" |

**快速诊断：**
```bash
# Mac 上测试连接
nc -zv 192.168.1.100 9999

# 如果显示 "succeeded" 说明连接正常
```

---

### 问题 2：依赖下载失败

**错误信息：**
```
go: github.com/atotto/clipboard@v0.1.4: Get "https://proxy.golang.org/...":
context deadline exceeded
```

**解决方案：**

1. **配置国内代理（推荐）**
```bash
go env -w GOPROXY=https://goproxy.cn,direct
go mod tidy
```

2. **如果在公司网络，配置代理**
```bash
go env -w GOPROXY=http://your-proxy:port
```

3. **离线方案：在有网络的 Mac 上下载**
```bash
go mod download
# 然后将整个项目文件夹复制到 Windows
```

---

### 问题 3：Windows 上找不到 go 命令

**解决方案：**

1. 确认 Go 已安装：在 PowerShell 中运行 `go version`
2. 如果找不到命令，重启 PowerShell 或电脑
3. 检查 Go 安装路径是否在系统 PATH 中：
```powershell
$env:Path -split ';' | Select-String 'Go'
```

---

### 问题 4：接收端启动后立即退出

**可能原因：**
- 端口 9999 已被占用
- 权限不足

**解决方案：**
```powershell
# 检查端口占用
netstat -ano | findstr :9999

# 使用其他端口
go run .\cmd\receiver\main.go -port 8888

# 对应地，Mac 上也要改端口
go run ./cmd/sender/main.go -ip 192.168.1.100 -port 8888
```

---

## 📝 免责声明

⚠️ **重要提示：**

- **仅限局域网使用** - 本工具设计用于同一局域网内的设备间通信
- **隐私安全** - 剪贴板内容以明文形式在网络上传输，请勿在公网或不信任的网络中使用
- **数据安全** - 不建议在传输敏感信息（密码、密钥等）时使用
- **网络安全** - 确保只在受信任的网络环境中运行此工具

**建议：**
- 仅在家庭网络或公司内网中使用
- 定期检查防火墙规则
- 不需要时及时关闭程序

---

## 📄 许可证

MIT License - 详见 LICENSE 文件

---

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！
