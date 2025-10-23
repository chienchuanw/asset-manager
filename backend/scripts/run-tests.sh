#!/bin/zsh

# 載入 zsh 環境（確保 go 在 PATH 中）
source ~/.zshrc 2>/dev/null || true

# 執行 gotestsum
~/go/bin/gotestsum "$@"

