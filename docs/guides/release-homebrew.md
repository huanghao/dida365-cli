# Homebrew 打包与发布指南

更新时间：2026-02-23

结论：
- 对 macOS 用户，Homebrew 是最合适的一线分发方式。
- 推荐模式：`dida365-cli`（源码+二进制发布） + `homebrew-tap`（formula 仓库）。

## 1. 先回答你的问题

Q1：`homebrew-tap` 是单独仓库吗？像多个包的索引目录？  
A：是。它通常是单独仓库，里面放 `Formula/*.rb`，可以放多个工具的 formula，相当于你自己的包索引。

Q2：二进制放哪里？是放 tap 仓库还是 dida 仓库？  
A：标准做法是二进制放在 `dida365-cli` 的 GitHub Release 附件里；tap 仓库只放 formula（URL+sha256），不存二进制。

Q3：tag 打在 `dida365-cli`，tap 怎么知道版本？  
A：tap 里的 formula 手动或自动更新 `url` 和 `sha256`。`url` 里会包含版本 tag（如 `v0.1.0`），所以 tap 通过 formula 指向具体版本。

Q4：CI 是什么？要自己做还是有现成工具？  
A：CI 是持续集成流水线（通常用 GitHub Actions）。你可以自己写 shell，也可以用现成工具 `goreleaser`。推荐 `goreleaser`，它可以自动：
- cross build
- 打包压缩
- 生成/上传 GitHub Release
- 更新 Homebrew tap formula

## 2. 发布架构（你会操作的两个仓库）

1. `dida365-cli`
- 放源码
- 打 tag
- 产出 release 二进制

2. `homebrew-tap`
- 放 formula，例如 `Formula/dida365-cli.rb`
- 用户通过 `brew tap <owner>/homebrew-tap && brew install dida365-cli` 安装

## 3. 手工发布流程（先跑通一次）

按顺序执行：

1. 准备版本号  
- 例如 `v0.1.0`
- 确认本地测试通过：`go test ./...`

2. 打 tag 并推送
```bash
git tag v0.1.0
git push origin v0.1.0
```

3. 本地构建二进制并打包
```bash
mkdir -p dist
mkdir -p dist/darwin_arm64 dist/darwin_amd64
GOOS=darwin GOARCH=arm64 go build -o dist/darwin_arm64/dida365-cli ./cmd/dida
GOOS=darwin GOARCH=amd64 go build -o dist/darwin_amd64/dida365-cli ./cmd/dida

tar -C dist/darwin_arm64 -czf dist/dida365-cli_darwin_arm64.tar.gz dida365-cli
tar -C dist/darwin_amd64 -czf dist/dida365-cli_darwin_amd64.tar.gz dida365-cli
```

4. 计算 sha256
```bash
shasum -a 256 dist/dida365-cli_darwin_arm64.tar.gz
shasum -a 256 dist/dida365-cli_darwin_amd64.tar.gz
```

5. 在 `dida365-cli` 创建 GitHub Release（tag=`v0.1.0`）  
- 上传上面两个 tar.gz

6. 更新 `homebrew-tap/Formula/dida365-cli.rb`  
- 把 `url` 改成 release 附件 URL
- 把 `sha256` 改成对应值
- 提交到 tap 仓库

7. 本地验证安装
```bash
brew tap <your-org>/homebrew-tap
brew install dida365-cli
dida365-cli --help
```

## 4. CI 自动构建流程（推荐）

推荐技术栈：
- GitHub Actions
- GoReleaser

最小流程：
1. push tag（如 `v0.1.1`）
2. GitHub Actions 触发
3. GoReleaser 构建和打包
4. 自动创建 GitHub Release 并上传产物
5. 自动更新 `homebrew-tap` 的 formula 并提交

### 4.1 你需要的凭证

1. `GITHUB_TOKEN`（Actions 自带）
- 用于当前仓库 release 发布

2. `HOMEBREW_TAP_TOKEN`（你创建的 PAT）
- 用于写入另一个 tap 仓库
- 需要对 tap 仓库有写权限

### 4.2 `.goreleaser.yaml`（核心示例）

```yaml
project_name: dida365-cli

before:
  hooks:
    - go mod tidy
    - go test ./...

builds:
  - id: dida365-cli
    main: ./cmd/dida
    binary: dida365-cli
    goos: [darwin, linux]
    goarch: [amd64, arm64]

archives:
  - id: default
    builds: [dida365-cli]
    name_template: "dida365-cli_{{ .Version }}_{{ .Os }}_{{ .Arch }}"
    format: tar.gz

brews:
  - name: dida365-cli
    repository:
      owner: <your-org>
      name: homebrew-tap
      token: "{{ .Env.HOMEBREW_TAP_TOKEN }}"
    folder: Formula
    homepage: "https://github.com/<your-org>/dida365-cli"
    description: "Dida365 CLI"
    test: |
      system "#{bin}/dida365-cli", "--help"

release:
  github:
    owner: <your-org>
    name: dida365-cli
```

### 4.3 GitHub Actions 工作流示例

文件：`.github/workflows/release.yml`

```yaml
name: release

on:
  push:
    tags:
      - "v*"

jobs:
  goreleaser:
    runs-on: ubuntu-latest
    permissions:
      contents: write
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod

      - uses: goreleaser/goreleaser-action@v6
        with:
          version: latest
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          HOMEBREW_TAP_TOKEN: ${{ secrets.HOMEBREW_TAP_TOKEN }}
```

## 5. 建议推进顺序

1. 先手工发一个版本（`v0.1.0`）验证链路  
2. 再接入 GoReleaser + GitHub Actions  
3. 稳定后再扩展更多分发渠道（如 Scoop）
