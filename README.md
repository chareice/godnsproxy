# Go DNS Proxy

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

DNS 分流代理服务，基于 Go 语言开发，支持 OpenWRT 平台。

## 功能特性

- 支持 DNS over HTTPS(DoH)
- 根据域名智能分流
- 高性能并发处理
- 轻量级资源占用
- OpenWRT 集成支持

## 安装

### 标准安装

```bash
go install github.com/chareice/godnsproxy@latest
```

### OpenWRT 一键安装

```bash
sh -c "$(curl -kfsSL https://raw.githubusercontent.com/chareice/godnsproxy/main/scripts/openwrt-install.sh)"
```

## 使用

### 命令行参数

```
-f string
    域名文件路径 (必须)
-p int
    DNS服务端口 (默认 5300)
-c string
    国内DNS服务器 (默认 "223.5.5.5")
-t string
    可信DNS服务器 (默认 "https://1.1.1.1/dns-query")
```

### OpenWRT 服务管理

```bash
/etc/init.d/godnsproxy [start|stop|restart|status]
```

### 日志查看

#### 命令行查看

```bash
# 查看最近日志
logread -e godnsproxy

# 查看完整日志
logread | grep godnsproxy

# 实时监控日志
logread -f | grep godnsproxy
```

#### 调整日志级别

编辑配置文件 `/etc/config/godnsproxy`，修改 `log_level` 选项：

- `info` (默认)
- `debug` (详细调试信息)
- `warn` (仅警告和错误)

修改后重启服务生效：

```bash
/etc/init.d/godnsproxy restart
```

#### 网页界面查看

如果安装了 LuCI 界面：

1. 进入 系统 → 系统日志
2. 筛选 `godnsproxy` 相关条目

## 开发

### 构建二进制

```bash
# 构建当前平台
go build -o godnsproxy .

# 交叉编译(示例: Linux ARM64)
GOOS=linux GOARCH=arm64 go build -o godnsproxy-linux-arm64 .
```

### 测试

```bash
# 运行单元测试
go test ./...

# 带覆盖率测试
go test -cover ./...
```

### 发布准备

```bash
# 构建多平台二进制(示例)
GOOS=linux GOARCH=amd64 go build -o releases/godnsproxy-linux-amd64 .
GOOS=linux GOARCH=arm64 go build -o releases/godnsproxy-linux-arm64 .
GOOS=linux GOARCH=arm go build -o releases/godnsproxy-linux-armv7 .
```

## 许可证

MIT License
