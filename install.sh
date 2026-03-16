#!/usr/bin/env bash
set -euo pipefail

# oops installer
# Usage: curl -fsSL https://oops-cli.com/install.sh | bash

BASE_URL="${OOPS_BASE_URL:-https://oops-cli.com/releases}"
VERSION="0.1.0"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
  x86_64)  ARCH="amd64" ;;
  aarch64) ARCH="arm64" ;;
  arm64)   ARCH="arm64" ;;
  *)       echo "Error: unsupported architecture: $ARCH"; exit 1 ;;
esac

case "$OS" in
  darwin|linux) ;;
  *) echo "Error: unsupported OS: $OS"; exit 1 ;;
esac

URL="${BASE_URL}/oops_${OS}_${ARCH}.tar.gz"
TMP=$(mktemp -d)
trap 'rm -rf "$TMP"' EXIT

echo "Downloading oops v${VERSION} (${OS}/${ARCH})..."
if command -v curl &>/dev/null; then
  curl -fsSL "$URL" -o "$TMP/oops.tar.gz"
elif command -v wget &>/dev/null; then
  wget -qO "$TMP/oops.tar.gz" "$URL"
else
  echo "Error: curl or wget required"
  exit 1
fi

tar -xzf "$TMP/oops.tar.gz" -C "$TMP"

if [ -w "$INSTALL_DIR" ]; then
  mv "$TMP/oops" "$INSTALL_DIR/oops"
else
  echo "Installing to $INSTALL_DIR (requires sudo)..."
  sudo mv "$TMP/oops" "$INSTALL_DIR/oops"
fi
chmod +x "$INSTALL_DIR/oops"

echo ""
echo "oops v${VERSION} installed to $INSTALL_DIR/oops"
echo ""

# Detect shell and offer to add the hook automatically
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
  # Check if hook is already installed
  if [ -f "$RC_FILE" ] && grep -q "oops init" "$RC_FILE" 2>/dev/null; then
    echo "Shell hook already in $RC_FILE"
  else
    printf "Add shell hook to %s? [Y/n] " "$RC_FILE"
    # Handle piped input (curl | bash) — default to yes
    if [ -t 0 ]; then
      read -r REPLY
    else
      REPLY="y"
      echo "y"
    fi
    case "$REPLY" in
      [nN]*)
        echo ""
        echo "To add it manually:"
        echo "  echo '$HOOK_LINE' >> $RC_FILE"
        ;;
      *)
        if [ "$SHELL_NAME" = "fish" ]; then
          mkdir -p "$(dirname "$RC_FILE")"
        fi
        echo "" >> "$RC_FILE"
        echo "$HOOK_LINE" >> "$RC_FILE"
        echo "Added to $RC_FILE"
        ;;
    esac
  fi
fi

echo ""
echo "Restart your shell to activate. Run 'oops doctor' to verify."
