#!/usr/bin/env bash
# lib/conf.sh — エンジン共通設定ローダー。source して使う。

# エンジン自身のパス（ランチャーから export されていなければ BASH_SOURCE から導出）
if [ -z "${DOTFILE_ENGINE_LIB:-}" ]; then
    DOTFILE_ENGINE_LIB="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    DOTFILE_ENGINE_DIR="$(dirname "$DOTFILE_ENGINE_LIB")"
fi

resolve_dotfile_dir() {
    if [ -n "${DOTFILE_DIR:-}" ] && [ -d "$DOTFILE_DIR" ]; then
        echo "$DOTFILE_DIR"
        return
    fi

    local git_root
    git_root="$(git rev-parse --show-toplevel 2>/dev/null || true)"
    if [ -n "$git_root" ] && [ -f "$git_root/.infra-version" ]; then
        echo "$git_root"
        return
    fi

    if [ -d "$HOME/dotfile" ] && [ -f "$HOME/dotfile/.infra-version" ]; then
        echo "$HOME/dotfile"
        return
    fi

    echo "[dotfile] Error: データリポジトリが見つかりません。" >&2
    echo "  DOTFILE_DIR 環境変数を設定するか、データリポジトリ内で実行してください。" >&2
    return 1
}

DOTFILE="$(resolve_dotfile_dir)" || exit 1

# データリポジトリの sync.conf を読み込む
if [ -f "$DOTFILE/sync.conf" ]; then
    source "$DOTFILE/sync.conf"
else
    SYNC_AUTO=()
    SYNC_MANUAL=()
    SYNC_IGNORE=()
fi

check_version_compat() {
    local engine_version data_version
    engine_version="$(cat "$DOTFILE_ENGINE_DIR/VERSION" 2>/dev/null || echo "0.0.0")"
    data_version="$(cat "$DOTFILE/.infra-version" 2>/dev/null || echo "0.0.0")"

    local engine_major="${engine_version%%.*}"
    local data_major="${data_version%%.*}"

    if [ "$engine_major" != "$data_major" ]; then
        echo "[dotfile] WARNING: バージョン不整合" >&2
        echo "  エンジン: v${engine_version}" >&2
        echo "  データ:   v${data_version}" >&2
        echo "  メジャーバージョンが異なります。" >&2
    fi
}
