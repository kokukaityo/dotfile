# My Dotfile

[dotfile](https://github.com/kokukaityo/dotfile) エンジンで管理する個人設定リポジトリ。

## セットアップ

### 1. エンジンをインストール

```bash
git clone https://github.com/kokukaityo/dotfile.git ~/.local/share/dotfile
export PATH="$HOME/.local/share/dotfile/bin:$PATH"
```

### 2. 初期設定

```bash
dotfile setup
```

### 3. シェル起動時の自動同期（任意）

`~/.bashrc` or `~/.zshrc` に追加:

```bash
export DOTFILE_DIR="$HOME/dotfile"
export PATH="$HOME/.local/share/dotfile/bin:$PATH"
command -v dotfile >/dev/null && dotfile pull
command -v dotfile >/dev/null && dotfile status
```

## 使い方

- 設定ファイルを追加: カテゴリディレクトリにファイルを置き、`link.yaml` に symlink 定義を追加
- 同期: `dotfile push` / `dotfile pull`
- symlink 再配置: `dotfile link`

## 構成

| ディレクトリ | 用途 |
|---|---|
| `ai-agent/` | AI エージェント設定 |
| `editor/` | エディタ設定 |
| `shell/` | シェル設定 |

| ファイル | 用途 |
|---|---|
| `sync.conf` | 同期モード定義 |
| `.infra-version` | 互換エンジンバージョン |

詳細は [dotfile エンジン](https://github.com/kokukaityo/dotfile) を参照。
