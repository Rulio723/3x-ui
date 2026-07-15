# 3x-ui Rulio 定制与正式版编译记录

记录时间：2026-07-15 23:38
项目目录：`D:\Codex\3XUI`  
目标系统：Debian 13 amd64  
当前分支：`rulio-v3.5.0-custom`
上游正式版基线：GitHub 3x-ui `v3.5.0`，本地基线提交 `4e928a1c v3.5.0`

## 已添加的三个功能

### 1. 自定义 Clash / Mihomo 订阅模板

作用：

- `/clash/:subid` 输出你自己的 Clash/Mihomo YAML 模板。
- 模板里自动注入 3x-ui 生成的 `proxies`。
- 保留你配置里的 DNS、规则组、规则、策略组等模板结构。
- 模板中的 `secret` 已留空，避免把面板或 Clash Dashboard 密码写进订阅。

涉及文件：

- `internal/sub/rulio_clash_template.yaml`
- `internal/sub/clash_service.go`
- `internal/sub/rulio_clash_template_test.go`

相关提交：

```text
33cfbf0b feat: use Rulio clash subscription template
```

实现要点：

- `internal/sub/clash_service.go` 使用 `//go:embed rulio_clash_template.yaml` 嵌入模板。
- 新增 `buildRulioClashConfig(proxies)`，读取模板 YAML 后替换 `proxies` 字段。
- Clash 订阅生成时不再用默认最小模板，而是用 `rulio_clash_template.yaml`。
- 测试 `TestBuildRulioClashConfigInjectsProxies` 用来确认代理节点能正确注入。

### 2. 订阅页添加“导入到 Clash”按钮

作用：

- 订阅页面的 `CLASH` 行新增蓝色 `导入` 按钮。
- 点开后可选择：
  - `FlClash`
  - `Clash Verge`
  - `Clash`
- 这三个选项都调用 Clash 系客户端常用协议：

```text
clash://install-config?url=<URL编码后的Clash订阅地址>&name=<URL编码后的订阅名称>
```

涉及文件：

- `frontend/src/pages/sub/SubPage.tsx`
- `frontend/src/pages/sub/SubPage.css`

相关提交：

```text
f714a6a6 feat: add clash direct import buttons
```

改动摘要：

- 引入 `ImportOutlined` 图标。
- 新增 `openClient()`，用 `window.location.href` 打开本地客户端协议。
- 新增 `clashInstallUrl`，根据 `subClashUrl` 和订阅名称生成导入链接。
- 新增 `clashImportMenuItems`，菜单包含 `FlClash / Clash Verge / Clash`。
- 在 `CLASH` 订阅行的复制、二维码按钮前加入 `导入` 下拉按钮。
- CSS 增加 `.sub-import-btn`，固定按钮最小宽度，并让操作区垂直居中。

### 3. 网页导入 Clash / Mihomo YAML 模板

作用：

- 在面板设置 → 订阅设置 → `Clash / Mihomo` 页签中新增 `YAML 模板` 区域。
- 可点击 `导入 YAML` 选择 `.yaml` / `.yml` 文件，也可以直接粘贴模板内容。
- 保存设置时后端自动清理不适合 3x-ui 订阅输出的内容：
  - 删除顶层 `proxies`
  - 删除顶层 `proxy-providers`
  - 清空 `secret`
  - 把使用 `use: [provider]` 的策略组改成 `include-all: true`
- `/clash/:subid` 生成订阅时优先使用网页保存的模板，再自动注入 3x-ui 生成的节点。
- 模板为空时继续使用内置 `internal/sub/rulio_clash_template.yaml`。

涉及文件：

- `internal/clashtemplate/template.go`
- `internal/web/entity/entity.go`
- `internal/web/service/setting.go`
- `internal/sub/clash_service.go`
- `internal/sub/sub.go`
- `internal/sub/controller.go`
- `frontend/src/pages/settings/SubscriptionGeneralTab.tsx`
- `frontend/src/models/setting.ts`
- `frontend/src/schemas/setting.ts`
- `frontend/src/generated/*`

相关提交：

```text
8f253f92 feat: add Clash YAML template import
```

## 正式版编译说明

这次编译的是 GitHub 正式版 `v3.5.0` 基线上的本地定制版，不是 dev 构建。

关键点：

- `internal/config/version` 仍为 `3.5.0`。
- 不注入 `buildCommit`。
- `go build` 的 `-ldflags` 只使用 `-s -w`。
- 不使用类似下面这种 dev 注入参数：

```text
-X github.com/mhsanaei/3x-ui/v3/internal/config.buildCommit=...
```

原因：

- 3x-ui 的版本逻辑会把注入了 `buildCommit` 的二进制显示成 `dev+<commit>`。
- 正式版/本地稳定版需要保持 `buildCommit` 为空，这样 UI 显示 `v3.5.0`。

## 编译方式

必须使用 CGO 编译。

原因：

- 3x-ui 使用 `go-sqlite3`。
- 如果用 `CGO_ENABLED=0`，Debian 服务器运行会报错：

```text
Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work
```

## 本机编译工具路径

Go：

```text
D:\Codex\3XUI\work\tools\go\bin\go.exe
```

Zig：

```text
D:\Codex\3XUI\work\tools\zig-x86_64-windows-0.15.2\zig.exe
```

Zig 用作 Windows 到 Linux amd64 的 CGO 交叉编译 C 编译器。

## 完整编译命令

先构建前端，把订阅页改动写入 `internal/web/dist`：

```powershell
cd D:\Codex\3XUI\frontend
npm.cmd run build
npm.cmd run typecheck
npx.cmd eslint src/pages/sub/SubPage.tsx
```

然后测试 Clash 模板注入：

```powershell
cd D:\Codex\3XUI
$go = 'D:\Codex\3XUI\work\tools\go\bin\go.exe'
& $go test ./internal/sub -run TestBuildRulioClashConfigInjectsProxies
```

最后正式版 CGO 交叉编译：

```powershell
cd D:\Codex\3XUI

$go = 'D:\Codex\3XUI\work\tools\go\bin\go.exe'
$zig = 'D:\Codex\3XUI\work\tools\zig-x86_64-windows-0.15.2\zig.exe'
$out = 'D:\Codex\3XUI\outputs\x-ui'

$env:GOOS = 'linux'
$env:GOARCH = 'amd64'
$env:CGO_ENABLED = '1'
$env:CC = "$zig cc -target x86_64-linux-gnu"
$env:CXX = "$zig c++ -target x86_64-linux-gnu"

& $go build -trimpath -ldflags "-s -w" -o $out -v main.go
Get-FileHash -Algorithm SHA256 -LiteralPath $out
```

## 本次产物

输出文件：

```text
D:\Codex\3XUI\outputs\x-ui
```

文件大小：

```text
100884664 bytes
```

SHA256：

```text
C3668652994C2DD6D486A066524FAC3016A277D4943FAFF9E65B925511FDD02F
```

## 已执行验证

```text
npm run build
npm run typecheck
npm run lint
npm run test
CGO_ENABLED=1 go test ./internal/clashtemplate ./internal/sub
GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build
```

## Debian 13 amd64 部署命令

把二进制上传到服务器后执行：

```bash
sudo systemctl stop x-ui
sudo systemctl reset-failed x-ui
sudo install -m 755 ./x-ui /usr/local/x-ui/x-ui
sudo systemctl start x-ui
sudo systemctl status x-ui --no-pager
```

确认版本：

```bash
/usr/local/x-ui/x-ui version
```

或者打开面板确认 UI 显示为：

```text
v3.5.0
```

而不是：

```text
dev+...
```

## 新会话继续工作的提示词

如果新开 Codex 会话，可以直接贴这段：

```text
项目在 D:\Codex\3XUI，分支 rulio-v3.5.0-custom，基于 GitHub 3x-ui 正式版 v3.5.0，不要构建 dev 版本。

已有三个定制功能：
1. internal/sub/rulio_clash_template.yaml 作为 Clash/Mihomo 订阅模板，internal/sub/clash_service.go 用 go:embed 读取并注入 proxies。
2. frontend/src/pages/sub/SubPage.tsx 的 CLASH 订阅行有“导入”按钮，下拉 FlClash / Clash Verge / Clash，使用 clash://install-config?url=...&name=...。
3. 面板设置 → 订阅设置 → Clash / Mihomo 页签支持导入 YAML 模板，保存时后端清理 proxies、proxy-providers、secret，并把 provider 策略组改为 include-all。

编译必须是 Debian 13 amd64，CGO_ENABLED=1。不要注入 buildCommit，否则 UI 会显示 dev+...。Go 路径：
D:\Codex\3XUI\work\tools\go\bin\go.exe
Zig 路径：
D:\Codex\3XUI\work\tools\zig-x86_64-windows-0.15.2\zig.exe

正式编译命令使用：
GOOS=linux GOARCH=amd64 CGO_ENABLED=1
CC="$zig cc -target x86_64-linux-gnu"
CXX="$zig c++ -target x86_64-linux-gnu"
go build -trimpath -ldflags "-s -w" -o outputs\x-ui main.go
```
