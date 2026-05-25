# weflow

WeFlow 本机 HTTP / SSE API 的命令行客户端，完整覆盖 `HTTP-API.md` 中 §1–§8 全部接口。

## 安装

### 从源码

```bash
git clone <repo>
cd weflow
go build -o weflow .
sudo mv weflow /usr/local/bin/
```

### 从 release（GoReleaser 多平台二进制）

下载对应平台 `weflow_<os>_<arch>.tar.gz`（Windows 为 `.zip`），解压后把可执行文件放到 `PATH`。

## 快速开始

```bash
weflow config init                    # 在 ~/.config/weflow/config.yaml 写入模板
$EDITOR ~/.config/weflow/config.yaml  # 把 token 改成你的 Access Token
weflow health                         # 验证连通
weflow sessions --limit 5             # 看最近 5 个会话
weflow messages --talker xxx@chatroom --limit 10
weflow watch                          # 订阅实时消息（Ctrl-C 退出）
```

## 配置

按以下顺序找第一个存在的配置文件：

1. `--config <path>` 指定的路径
2. `./weflow.yaml`
3. `$XDG_CONFIG_HOME/weflow/config.yaml`（未设则 `~/.config/weflow/config.yaml`）
4. `~/.weflow.yaml`

```yaml
host: http://127.0.0.1:5031
token: ""             # 在 WeFlow 应用设置 → API 服务里复制
timeout: 30s
output:
  color: true
```

优先级（高 → 低）：**flag > 环境变量（`WEFLOW_HOST` / `WEFLOW_TOKEN` / `WEFLOW_TIMEOUT`） > 配置文件 > 内置默认**。

## 命令一览

| 命令 | 对应接口 |
|---|---|
| `weflow health` | §1 GET /health |
| `weflow watch [--reconnect]` | §2 SSE /api/v1/push/messages |
| `weflow messages --talker <id> [--chatlab] [--media …]` | §3 /api/v1/messages |
| `weflow sessions [--format chatlab]` | §4 / §4.1 /api/v1/sessions |
| `weflow pull <sessionId> [--follow]` | §4.2 /api/v1/sessions/:id/messages |
| `weflow contacts` | §5 /api/v1/contacts |
| `weflow members <chatroomId> [--with-counts]` | §6 /api/v1/group-members |
| `weflow sns timeline` | §7.1 |
| `weflow sns usernames` | §7.2 |
| `weflow sns stats` | §7.3 |
| `weflow sns proxy --url <url> [--key <k>]` | §7.4 |
| `weflow sns export --output-dir <dir>` | §7.5 |
| `weflow sns block-delete <status\|install\|uninstall>` | §7.6 |
| `weflow sns delete <postId>` | §7.7 |
| `weflow media get <relativePath> [-o file]` | §8 /api/v1/media/* |
| `weflow config {init\|show\|path}` | 本地配置管理 |
| `weflow version` | 版本信息 |

任何子命令都可以切换输出格式：
- 默认：圆框 pretty 表格（终端友好）
- `--text`：每行一条记录，`[时间] 发送者: 内容` 格式，适合直接重定向到 `.txt`
- `--json`：后端原样 JSON，适合 `jq` / 入库

例：导出一个群一天的聊天记录到 txt（推荐 `--chatlab` 让发送者显示为人名而非 wxid）：
```bash
weflow messages \
  --talker xxx@chatroom \
  --start 20260525 --end 20260525 \
  --limit 10000 \
  --chatlab --text > chat-20260525.txt
```

## Shell 补全

```bash
# zsh
weflow completion zsh > ~/.zsh/completion/_weflow

# bash
weflow completion bash > /usr/local/etc/bash_completion.d/weflow

# fish
weflow completion fish > ~/.config/fish/completions/weflow.fish
```

## 开发

```bash
go test ./...                   # 跑单元测试
make build                      # 本机构建
make snapshot                   # GoReleaser 多平台快照
make release                    # 真正发布（需 GITHUB_TOKEN + git tag）
```
