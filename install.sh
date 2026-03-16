#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${OOPS_BASE_URL:-https://oops-cli.com/releases}"
VERSION="0.1.0"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Colors
R='\033[0;31m'
G='\033[0;32m'
B='\033[1m'
D='\033[0;90m'
N='\033[0m'

info()  { echo -e "  ${D}>${N} $1"; }
ok()    { echo -e "  ${G}✓${N} $1"; }
err()   { echo -e "  ${R}✗${N} $1"; exit 1; }

ask() {
  printf "  %b [Y/n] " "$1"
  if [ -t 0 ]; then
    read -r REPLY
  else
    REPLY="y"
    echo "y"
  fi
  case "$REPLY" in
    [nN]*) return 1 ;;
    *) return 0 ;;
  esac
}

# ── Header ───────────────────────────────────────────

echo ""
echo -e "${R}  oops${N} installer ${D}v${VERSION}${N}"
echo -e "${D}  undo for your terminal${N}"
echo ""

# ── Detect platform ──────────────────────────────────

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
  x86_64)  ARCH="amd64" ;;
  aarch64) ARCH="arm64" ;;
  arm64)   ARCH="arm64" ;;
  *)       err "Unsupported architecture: $ARCH" ;;
esac

case "$OS" in
  darwin|linux) ;;
  *) err "Unsupported OS: $OS" ;;
esac

info "Detected ${B}${OS}/${ARCH}${N}"
echo ""

# ── Step 1: Download and install binary ──────────────

if ask "Install oops to ${B}${INSTALL_DIR}${N}?"; then
  URL="${BASE_URL}/oops_${OS}_${ARCH}.tar.gz"
  TMP=$(mktemp -d)
  trap 'rm -rf "$TMP"' EXIT

  info "Downloading..."
  if command -v curl &>/dev/null; then
    curl -fsSL "$URL" -o "$TMP/oops.tar.gz"
  elif command -v wget &>/dev/null; then
    wget -qO "$TMP/oops.tar.gz" "$URL"
  else
    err "curl or wget required"
  fi

  tar -xzf "$TMP/oops.tar.gz" -C "$TMP"

  if [ -w "$INSTALL_DIR" ]; then
    mv "$TMP/oops" "$INSTALL_DIR/oops"
  else
    info "Requires sudo for ${INSTALL_DIR}"
    sudo mv "$TMP/oops" "$INSTALL_DIR/oops"
  fi
  chmod +x "$INSTALL_DIR/oops"

  ok "Installed to ${INSTALL_DIR}/oops"
else
  echo -e "  ${D}Skipped.${N}"
  echo ""
  exit 0
fi

echo ""

# ── Step 2: Add shell hook ───────────────────────────

SHELL_NAME=$(basename "${SHELL:-zsh}")
HOOK_LINE=""
RC_FILE=""

case "$SHELL_NAME" in
  zsh)
    HOOK_LINE='eval "$(oops init zsh)"'
    RC_FILE="$HOME/.zshrc"
    ;;
  bash)
    HOOK_LINE='eval "$(oops init bash)"'
    if [ -f "$HOME/.bashrc" ]; then
      RC_FILE="$HOME/.bashrc"
    elif [ -f "$HOME/.bash_profile" ]; then
      RC_FILE="$HOME/.bash_profile"
    else
      RC_FILE="$HOME/.bashrc"
    fi
    ;;
  fish)
    HOOK_LINE='oops init fish | source'
    RC_FILE="$HOME/.config/fish/config.fish"
    ;;
esac

if [ -n "$RC_FILE" ]; then
  if [ -f "$RC_FILE" ] && grep -q "oops init" "$RC_FILE" 2>/dev/null; then
    ok "Shell hook already in ${RC_FILE}"
  elif ask "Add shell hook to ${B}${RC_FILE}${N}?"; then
    if [ "$SHELL_NAME" = "fish" ]; then
      mkdir -p "$(dirname "$RC_FILE")"
    fi
    echo "" >> "$RC_FILE"
    echo "$HOOK_LINE" >> "$RC_FILE"
    ok "Added to ${RC_FILE}"
  else
    echo ""
    info "To add it manually later:"
    info "  echo '${HOOK_LINE}' >> ${RC_FILE}"
  fi
fi

echo ""

# ── Step 3: Create oops directory ────────────────────

if [ ! -d "$HOME/.oops" ]; then
  if ask "Create ${B}~/.oops${N} directory for backups?"; then
    mkdir -p "$HOME/.oops/trash"
    ok "Created ~/.oops"
  else
    echo -e "  ${D}Skipped (will be created on first use).${N}"
  fi
else
  ok "~/.oops already exists"
fi

echo ""

# ── Done ─────────────────────────────────────────────

echo -e "  ${G}Done!${N} Open a new terminal tab, then try:"
echo ""
echo -e "  ${D}\$${N} rm something.txt"
echo -e "  ${D}\$${N} ${R}oops${N}"
echo -e "  ${G}✓${N} restored something.txt"
echo ""
