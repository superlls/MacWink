# 🍎 MacWink 🪟

[English](#english) | [中文](#chinese)

> **Break the ecosystem wall.** 🔨
> 打破 Mac 与 Windows 的生态隔阂，让剪贴板在局域网内自由流转。

<a id="chinese"></a>
## 📖 背景

你是否厌倦了：
*   为了发一段文字，想用微信QQ中转，两台设备却不支持同时在线？
*   为了跨设备复制，还要去顾虑云端的隐私风险？

**MacWink** 的诞生就是为了解决这个问题。它是一个极简的局域网 P2P 同步工具，**无需登录、不走云端**，让 macOS、Windows 和 Linux 设备间的剪贴板实现毫秒级同步。

## 🛠️ 快速开始

确保两端在同一局域网，且已安装 Go 1.25+。

### 1️⃣ 核心逻辑
*   **Peer-to-Peer**: 无中心服务器，两端直连。
*   **配置**: 只需告诉程序对方的 IP (`-peer`)。

### 2️⃣ 启动命令

**在 A 电脑上 (告诉它 B 的 IP):**
```bash
# 假如 B 的 IP 是 192.168.1.100
go run main.go -peer 192.168.1.100
```

**在 B 电脑上 (告诉它 A 的 IP):**
```bash
# 假如 A 的 IP 是 192.168.1.101
go run main.go -peer 192.168.1.101
```

*(默认使用 TCP 端口 `9999`，Windows 用户请记得允许防火墙通过)*

### 🛑 停止/杀死进程
前台运行时直接 `Ctrl+C` 即可退出；如果你是后台运行或端口被占用，可以按端口杀掉进程：
```bash
# macOS / Linux：杀死占用 9999 端口的进程（优雅退出：SIGTERM）
kill -TERM $(lsof -ti tcp:9999)
```

### ⚙️ 进阶配置
如果端口冲突，可以自定义：
```bash
# 本地监听 8888，并连接对方的 8888
go run main.go -port 8888 -peer 192.168.1.100:8888
```

## ⚡️ 生产力 Combo

> ### [🎙️ CodeWhisper](https://github.com/superlls/CodeWhisper) + [🚀 MacWink](https://github.com/your-username/MacWink)
1.  使用 **CodeWhisper** 将你的语音灵感实时转写为文字，自动存入剪贴板。
2.  **MacWink** 立即接力，将这段文字无缝同步到你桌面的另一台电脑上。

👉 **场景**：对着 Mac 说话，文字直接出现在 Windows上

## 📜 License

本项目基于 [MIT License](LICENSE) 开源。依赖 `atotto/clipboard` 库。

<a id="english"></a>

## 📖 Background

Are you tired of:
*   Copying a piece of text but having to relay it through WeChat/QQ, only to find that the two devices can’t be online at the same time?
*   Copying across devices while still worrying about privacy risks in the cloud?

**MacWink** is built to solve exactly that. It’s a minimalist LAN P2P sync tool: **no accounts, no cloud**. It enables millisecond-level clipboard sync across macOS, Windows, and Linux devices.

## 🛠️ Quick Start

Make sure both machines are on the same LAN and have Go 1.25+ installed.

### 1️⃣ Core Logic
*   **Peer-to-Peer**: No central server; the two peers connect directly.
*   **Config**: You only need to tell the program the other peer’s IP (`-peer`).

### 2️⃣ Run

**On machine A (tell it B’s IP):**
```bash
# 假如 B 的 IP 是 192.168.1.100
go run main.go -peer 192.168.1.100
```

**On machine B (tell it A’s IP):**
```bash
# 假如 A 的 IP 是 192.168.1.101
go run main.go -peer 192.168.1.101
```

*(Uses TCP port `9999` by default. On Windows, remember to allow it through the firewall.)*

### 🛑 Stop / Kill the Process
If you’re running it in the foreground, press `Ctrl+C` to exit. If it’s running in the background or the port is in use, you can kill the process by port:
```bash
# macOS / Linux：杀死占用 9999 端口的进程（优雅退出：SIGTERM）
kill -TERM $(lsof -ti tcp:9999)
```

### ⚙️ Advanced Configuration

如果端口冲突，可以自定义：

If the port is already in use, you can customize it:
```bash
# 本地监听 8888，并连接对方的 8888
go run main.go -port 8888 -peer 192.168.1.100:8888
```

## ⚡️ Productivity Combo

> ### [🎙️ CodeWhisper](https://github.com/superlls/CodeWhisper) + [🚀 MacWink](https://github.com/your-username/MacWink)
1.  使用 **CodeWhisper** 将你的语音灵感实时转写为文字，自动存入剪贴板。
2.  **MacWink** 立即接力，将这段文字无缝同步到你桌面的另一台电脑上。

1.  Use **CodeWhisper** to transcribe your voice ideas into text in real time and automatically write them to the clipboard.
2.  **MacWink** picks it up immediately and seamlessly syncs that clipboard content to your other desktop machine.

👉 **场景**：对着 Mac 说话，文字直接出现在 Windows上

👉 **Scenario**: Speak to your Mac, and the text shows up on Windows instantly.

## 📜 License

本项目基于 [MIT License](LICENSE) 开源。依赖 `atotto/clipboard` 库。

This project is open-sourced under the [MIT License](LICENSE) and depends on the `atotto/clipboard` library.
