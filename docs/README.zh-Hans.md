# 1Panel MCP Server（扩展版 Fork）

> 🍴 本仓库是 [1Panel-dev/mcp-1panel](https://github.com/1Panel-dev/mcp-1panel) 的扩展 fork——新增 4 个工具，并修复了适配 1Panel v2 的 API 基础路径。采用与上游相同的 [GPL-3.0](LICENSE) 许可证。

**1Panel MCP Server** 是 [1Panel](https://github.com/1Panel-dev/1Panel) 的 Model Context Protocol (MCP) 协议服务端实现。

---

## 与上游的差异

| 变更 | 说明 |
|---|---|
| ➕ 新增 4 个工具 | `list_runtimes`、`list_containers`、`get_container_logs`、`call_panel_api`，源码位于 [`operations/extended/`](operations/extended/) |
| 🌐 覆盖全部 API | `call_panel_api` 可直接调用**任意** 1Panel API 端点（GET/POST/PUT/DELETE），凡是没有专属工具的面板功能都能通过它访问 |
| 🔧 API 基址修复 | 上游 `main` 分支的 `ApiBase = "/api/v1"`，而 1Panel v2.x 只提供 `/api/v2`（其余路径全部 404）。本 fork 已改为 `/api/v2` |

> 本 fork 的 Linux 二进制发布在本仓库 [Releases 页面](https://github.com/dtyg123/mcp-1panel-extended/releases)（amd64 / arm64 / armv7 / s390x / ppc64le）。官方 `1panel/1panel-mcp-server` Docker 镜像只包含上游的工具。

---

## 安装方式

> 注意：方式四（Docker 镜像）安装的是**上游**官方镜像，不含扩展工具。方式一（Release 下载）与方式二（源码构建）均包含扩展工具。

### ✅ 方式一：从 Release 页面下载（推荐）

1. 访问本仓库 [Releases 页面](https://github.com/dtyg123/mcp-1panel-extended/releases)，下载对应系统的可执行文件。

2. 安装示例（以 `amd64` 为例）：

```bash
chmod +x mcp-1panel-linux-amd64
mv mcp-1panel-linux-amd64 /usr/local/bin/mcp-1panel
```

### 🛠️ 方式二：通过源码构建

确保本地已安装 Go 1.23 或更高版本，执行以下命令：

1. 克隆代码仓库：

```bash
git clone https://github.com/dtyg123/mcp-1panel-extended.git
cd mcp-1panel-extended
```

2. 构建可执行文件：

```bash
make build
```

3. 可执行文件生成路径为：`./build/mcp-1panel`，建议移动到系统 PATH 目录中。

### 🚀 方式三：通过 `go install` 安装

确保本地已安装 Go 1.23 或更高版本：

```bash
go install github.com/dtyg123/mcp-1panel-extended@latest
```

---

## 使用方式

1Panel MCP Server 支持两种运行模式：**stdio** 和 **sse**。

### 模式一：stdio（默认）

在 Cursor、Windsurf 或其他 MCP 客户端的配置文件中添加：

```json
{
  "mcpServers": {
    "mcp-1panel": {
      "command": "mcp-1panel",
      "env": {
        "PANEL_ACCESS_TOKEN": "<your 1Panel access token>",
        "PANEL_HOST": "such as http://localhost:8080"
      }
    }
  }
}
```

### 模式二：sse

启动 MCP Server：

```bash
mcp-1panel -host http://localhost:8080 -token <your 1Panel access token> -transport sse -addr http://localhost:8000
```

客户端配置：

```json
{
  "mcpServers": {
    "mcp-1panel": {
      "url": "http://localhost:8000/sse"
    }
  }
}
```

### 🔧 命令行参数

- `-token`: 1Panel 的访问令牌
- `-host`: 1Panel 的地址，如：http://localhost:8080
- `-transport`: 传输方式：`stdio` 或 `sse`，默认是 `stdio`
- `-addr`: SSE 服务监听地址，默认是 `http://localhost:8000`

### 环境变量

- `PANEL_ACCESS_TOKEN`: 1Panel 访问令牌（面板设置中生成）
- `PANEL_HOST`: 1Panel 访问地址

---

## 🧰 可用工具（Tools）

以下是 MCP Server 提供的工具列表（最后 4 个带 🆕 标记的是本 fork 新增）：

| 工具名称                | 分类        | 描述                             |
|-------------------------|-------------|----------------------------------|
| `get_dashboard_info`    | System      | 获取仪表盘状态                   |
| `get_system_info`       | System      | 获取系统信息                     |
| `list_websites`         | Website     | 列出所有网站                     |
| `create_website`        | Website     | 创建新网站                       |
| `list_ssls`             | Certificate | 列出所有证书                     |
| `create_ssl`            | Certificate | 创建新证书                       |
| `list_installed_apps`   | Application | 列出已安装应用                   |
| `install_openresty`     | Application | 安装 OpenResty                   |
| `install_mysql`         | Application | 安装 MySQL                       |
| `list_databases`        | Database    | 列出所有数据库                   |
| `create_database`       | Database    | 创建新数据库                     |
| `list_runtimes` 🆕      | Runtime     | 列出运行环境（PHP、Go、Node、Python、Java、.NET）及版本、状态 |
| `list_containers` 🆕    | Container   | 列出 Docker 容器（状态、镜像、端口、网络），支持过滤 |
| `get_container_logs` 🆕 | Container   | 获取容器日志文本（tail / follow / since / timestamp） |
| `call_panel_api` 🆕     | General     | 直接调用任意 1Panel API 端点（GET/POST/PUT/DELETE） |

---

## 致谢

本 fork 基于 [1Panel-dev/mcp-1panel](https://github.com/1Panel-dev/mcp-1panel) 构建，感谢原作者与 1Panel 社区。
