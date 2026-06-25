# My Dotfiles

[dotfiles-bridge](https://github.com/kokukaityo/dotfiles-bridge) で管理する個人設定リポジトリ。

## 使い方

- 設定ファイルを追加: カテゴリディレクトリにファイルを置き、`link.toml` に symlink 定義を追加
- 同期: `dotfiles push` / `dotfiles pull`
- symlink 再配置: `dotfiles link`

## 構成

| ディレクトリ | 用途                |
| ------------ | ------------------- |
| `ai-agent/`  | AI エージェント設定 |
| `editor/`    | エディタ設定        |
| `shell/`     | シェル設定          |

| ファイル         | 用途               |
| ---------------- | ------------------ |
| `sync.toml`      | 同期モード定義     |
| `.infra-version` | 互換本体バージョン |

詳細は [dotfiles-bridge](https://github.com/kokukaityo/dotfiles-bridge) を参照。
