# Homebrew 发布最小流程

更新时间：2026-02-23

目标：回答“命令行是不是用 Homebrew 发布最好，以及需要准备什么”。

结论（当前项目）：
- 对 macOS 用户，Homebrew 是最合适的一线分发方式。
- 建议采用“GitHub Release + Homebrew Tap”模式。

## 1. 你需要先准备的内容

1. GitHub 仓库
- 源码仓库：`dida365-cli`
- 一个 tap 仓库（建议）：`homebrew-tap`

2. 可复现构建产物
- 至少提供：
  - `dida_darwin_arm64.tar.gz`
  - `dida_darwin_amd64.tar.gz`
  - （可选）`dida_linux_amd64.tar.gz`
- 每个压缩包里只放可执行文件 `dida`

3. 版本规范
- 使用 git tag（例如 `v0.1.0`）
- 每次 release 固定对应一个 tag

4. 校验值
- 每个产物要有 SHA256（供 brew formula 校验）

## 2. 发布步骤（最小可执行）

1. 本地打 tag 并推送
```bash
git tag v0.1.0
git push origin v0.1.0
```

2. 生成多平台二进制并上传到 GitHub Release
- 方式 A：手工 `go build` 后上传
- 方式 B：用 CI 自动构建上传（推荐）

3. 在 tap 仓库新增 formula（`Formula/dida.rb`）
- 指向 release 的 tar.gz URL
- 填写 SHA256

4. 用户安装
```bash
brew tap <your-org>/tap
brew install dida
```

## 3. Formula 示例（darwin arm64）

```ruby
class Dida < Formula
  desc "Dida365 CLI"
  homepage "https://github.com/<your-org>/dida365-cli"
  url "https://github.com/<your-org>/dida365-cli/releases/download/v0.1.0/dida_darwin_arm64.tar.gz"
  sha256 "<fill-sha256>"
  version "0.1.0"

  def install
    bin.install "dida"
  end

  test do
    assert_match "Dida365 CLI", shell_output("#{bin}/dida --help")
  end
end
```

## 4. 当前仓库还缺什么

- 缺 `LICENSE`（本轮已补）
- 缺正式版本 tag 与 release 产物
- 缺 tap 仓库和 formula
- 缺自动化发布流水线（可后续补 GitHub Actions）

## 5. 建议执行顺序

1. 先发 `v0.1.0` 手工 release（验证流程）
2. 再补 CI 自动构建与自动更新 formula
3. 稳定后再考虑 Scoop / apt / npm 等多渠道分发
