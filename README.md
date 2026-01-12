# 🍎 MacWink 🪟

> **Break the ecosystem wall.** 🔨
> 打破 Mac 与 Windows 的生态隔阂，让剪贴板在局域网内自由流转。

## 📖 背景 (Background)

你是否厌倦了：
*   为了发一段文字，想用微信QQ中转，两台设备却不支持同时在线？
*   为了跨设备复制，还要去顾虑云端的隐私风险？

**MacWink** 的诞生就是为了解决这个问题。它是一个极简的局域网 P2P 同步工具，**无需登录、不走云端**，让 macOS、Windows 和 Linux 设备间的剪贴板实现毫秒级同步。


## 🛠️ 快速开始 (Quick Start)

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

### ⚙️ 进阶配置
如果端口冲突，可以自定义：
```bash
# 本地监听 8888，并连接对方的 8888
go run main.go -port 8888 -peer 192.168.1.100:8888
```
## ⚡️生产力Combo
<a href="https://github.com/superlls/CodeWhisper">
  <img align="center" src="https://github-readme-stats.vercel.app/api/pin/?username=superlls&repo=CodeWhisper&theme=transparent" />
</a>
**🎙️ CodeWhisper + 🚀 MacWink**

生产力组合：
1.  使用 **CodeWhisper** 将你的语音灵感实时转写为文字，自动存入剪贴板。
2.  **MacWink** 立即接力，将这段文字无缝同步到你桌面的另一台电脑上。

👉 **场景**：对着 Mac 说话，文字直接出现在 Windows上 **Talk here, Paste there.**

## 📜 License

本项目基于 [MIT License](LICENSE) 开源。依赖 `atotto/clipboard` 库。
