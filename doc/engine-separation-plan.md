# dotfile エンジン分離 実装プラン

## Context

現在の `kokukaityo/dotfile` は `.infra/`（同期・リンクエンジン）と個人設定データ（ai-agent/ 等）が1つの private repo に同居している。二人目のユーザーが現れたことで、エンジンとデータの分離トリガーが満たされた。

**ゴール**: エンジンを public OSS として配布可能にし、`dotfile init` コマンドでユーザーが自分のデータリポジトリを即座に作れる構造に移行する。

**全体ロードマップ**:
1. **今回のスコープ**: shell のままエンジンとデータを分離し、新構造を確立する
2. **次のフェーズ（別作業）**: Go でエンジンを書き直す（シングルバイナリ配布）

Go で書き直す前提なので、shell 段階では installer や self-update を作り込まない。構造と設計の確立に集中する。

**分離後の構成**:
| リポジトリ | 公開 | 用途 |
|---|---|---|
| `kokukaityo/dotfile` | public (MIT) | エンジン + テンプレート同梱 |
| 各ユーザーの private repo | private | 個人設定データ（`dotfile init` で生成） |

テンプレートリポジトリは独立して作らず、エンジンリポジトリ内の `template/` から `dotfile init` でローカル生成する。

---

## Phase 0: 作業環境の準備

1. リポジトリを `D:\work` に移動
2. `feature/engine-separation` ブランチを作成
3. このプランファイルを `doc/engine-separation-plan.md` として保存
4. 既存の設計ドキュメントを `doc/archive/` に移動（`doc/dotfile-engine-separation-plan.md` 等）

---

## Phase 1: エンジンリポジトリの構造に変換

### 目標構造
```
dotfile/
├── bin/
│   └── dotfile              # ランチャー（エントリポイント）
├── lib/
│   ├── conf.sh              # パス解決 + sync.conf ローダー + バージョンチェック
│   ├── link.sh              # symlink エンジン
│   ├── sync.sh              # 同期エンジン
│   └── hook/
│       ├── pre-push
│       └── post-merge
├── template/                # dotfile init 用の雛形
│   ├── ai-agent/
│   │   └── link.yaml
│   ├── editor/
│   │   └── link.yaml
│   ├── shell/
│   │   └── link.yaml
│   ├── sync.conf
│   ├── .infra-version
│   ├── .gitignore
│   └── README.md
├── doc/
│   ├── engine-separation-plan.md  # この実装計画
│   └── archive/                    # 旧設計ドキュメントのバックアップ
├── VERSION                  # "1.0.0"
├── LICENSE                  # MIT
└── README.md
```

**注**: `installer/install.sh` と `self-update` は Go 版で実装するため、shell 段階では作らない。当面のインストールは `git clone` + PATH 追加で行う。

### 変更詳細

#### 新規作成

**`bin/dotfile`** — エントリポイント。
- `DOTFILE_ENGINE_DIR` / `DOTFILE_ENGINE_LIB` を export
- `lib/conf.sh` を source してデータリポジトリのパスを解決
- サブコマンドルーティング:
  - `init [path]` — template/ を指定パス（デフォルト `$HOME/dotfile`）にコピー、git init、initial commit、setup 実行。既にディレクトリが存在したらエラー
  - `setup` — hook 設定・gitattributes・gitignore 生成・symlink 配置（現 setup.sh の内容）
  - `link` — `lib/link.sh` に委譲
  - `pull` / `push` / `delete-category` / `gitignore` / `status` — `lib/sync.sh` に委譲
  - `version` — VERSION ファイルの内容とパスを表示
  - `help` — 使い方表示
- macOS の `readlink -f` 非対応を考慮した自前シンボリックリンク解決関数
- `init` はデータリポジトリ不要で動く必要があるため、`conf.sh` の source 前にハンドリング

**`VERSION`** — `1.0.0`

**`LICENSE`** — MIT

**`template/`** — `dotfile init` で使うデータリポジトリの雛形:
- `sync.conf`: 同期モード定義（現 conf.sh から移動）
  ```bash
  SYNC_AUTO=(ai-agent editor shell)
  SYNC_MANUAL=()
  SYNC_IGNORE=(backup raw)
  ```
- `.infra-version`: `1.0.0`
- 各カテゴリの `link.yaml`: コメントのみのサンプル
- `.gitignore`: セキュリティ除外のみ
- `README.md`: セットアップ手順・使い方

#### 移植・改変（.infra/ → lib/）

**`lib/conf.sh`**（元 `.infra/conf.sh`）— 全面書き換え:
- エンジンパス: `DOTFILE_ENGINE_DIR` / `DOTFILE_ENGINE_LIB`
- データリポジトリパス `DOTFILE` の解決（優先順位）:
  1. `DOTFILE_DIR` 環境変数
  2. カレント git root に `.infra-version` が存在すればそこ
  3. `$HOME/dotfile`
  4. エラー
- `sync.conf` を source して `SYNC_AUTO` / `SYNC_MANUAL` / `SYNC_IGNORE` を取得
- `check_version_compat()`: メジャーバージョン比較、不一致時は WARNING のみ

**`lib/link.sh`**（元 `.infra/link.sh`）— source 行のみ。`lib/` 内の相対パスで `conf.sh` を見つけるので現行構造と同じ

**`lib/sync.sh`**（元 `.infra/sync.sh`）— 以下を変更:
- `cmd_gitignore()` コメント: `conf.sh` → `sync.conf`
- `cmd_push()` エラーメッセージ: `bash .infra/sync.sh delete-category` → `dotfile delete-category`
- `cmd_delete_category()`:
  - 変更検知対象: `.infra/conf.sh` → `sync.conf`
  - awk 書き換え対象: `.infra/conf.sh` → `sync.conf`

**`lib/hook/pre-push`**（元 `.infra/hook/pre-push`）— source 行は現行構造と同じで動く

**`lib/hook/post-merge`**（元 `.infra/hook/post-merge`）— `bash "$DOTFILE/.infra/link.sh"` → `dotfile link` に変更

#### バックアップ（doc/archive/ に移動）
- `doc/dotfile-engine-separation-plan.md` → `doc/archive/`
- その他 `doc/` 内の旧ドキュメント → `doc/archive/`

#### 削除するもの
- `.infra/`（`lib/` に移行済み）
- `ai-agent/`, `editor/`, `shell/`（個人データ。template/ にサンプルとして残る）
- `.devcontainer/`, `.vscode/`

---

## Phase 2: README

**エンジン README.md**:
- 概要（何をするツールか）
- インストール手順（当面は `git clone` + PATH 追加）
- `dotfile init` でデータリポジトリを作る手順
- サブコマンド一覧
- シェル起動時の自動同期設定例
- ライセンス

---

## 検証方法

1. `dotfile version` — エンジンのバージョンとパスが正しく表示されるか
2. `dotfile init ~/test-dotfile` — template/ がコピーされ git init されるか
3. `dotfile setup`（データリポジトリ内で）— hook 設定・gitignore・symlink が動くか
4. `dotfile push` / `dotfile pull` — 同期が正常に動作するか
5. `.infra-version` を意図的にずらして WARNING が出るか

---

## 実装順序

| # | 作業 | 主要ファイル |
|---|------|-------------|
| 0 | リポジトリを `D:\work` に移動、ブランチ作成 | — |
| 1 | プランファイルを doc/ に保存、旧ドキュメントを doc/archive/ にバックアップ | `doc/` |
| 2 | `lib/conf.sh` — パス解決・sync.conf ローダー・バージョンチェック | `lib/conf.sh` |
| 3 | `lib/link.sh` 移植 | `lib/link.sh` |
| 4 | `lib/sync.sh` 移植（参照パス変更） | `lib/sync.sh` |
| 5 | `lib/hook/` 移植 | `lib/hook/pre-push`, `lib/hook/post-merge` |
| 6 | `bin/dotfile` ランチャー（init, setup, ルーティング） | `bin/dotfile` |
| 7 | `template/` 作成（sync.conf, .infra-version, link.yaml, README, .gitignore） | `template/*` |
| 8 | `VERSION`, `LICENSE` | `VERSION`, `LICENSE` |
| 9 | エンジン `README.md` | `README.md` |
| 10 | 個人データ・旧ファイル削除、.infra/ 削除 | — |

---

## 今回スコープ外（Go 化フェーズで実装）

- `installer/install.sh`（`curl | bash` インストーラ）
- `dotfile self-update` サブコマンド
- シングルバイナリ配布（GitHub Releases）
- Homebrew tap
- Docker E2E テスト
- bats テスト・GitHub Actions CI
