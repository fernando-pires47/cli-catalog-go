#!/bin/sh
set -eu

REPO_OWNER="fernando-pires47"
REPO_NAME="cli-catalog-go"
BIN_NAME="cs"
SYSTEM_BIN_DIR="/usr/local/bin"
LOCAL_BIN_DIR="${HOME}/.local/bin"
requested_version=""
resolved_version=""
checksums_tmp_file=""

info() {
  printf '==> %s\n' "$1"
}

warn() {
  printf 'warning: %s\n' "$1" >&2
}

fail() {
  printf 'error: %s\n' "$1" >&2
  exit 1
}

need_cmd() {
  command -v "$1" >/dev/null 2>&1 || fail "required command not found: $1"
}

usage() {
  cat <<'EOF'
Usage:
  sh install.sh [--version <tag>]

Options:
  --version <tag>  Install a specific release tag (example: 0.1.0)
  --help           Show this help message
EOF
}

parse_args() {
  while [ "$#" -gt 0 ]; do
    case "$1" in
      --version)
        [ "$#" -ge 2 ] || fail "missing value for --version"
        requested_version="$2"
        shift 2
        ;;
      --help)
        usage
        exit 0
        ;;
      *)
        usage
        fail "unknown argument: $1"
        ;;
    esac
  done
}

detect_os() {
  [ "$(uname -s)" = "Linux" ] || fail "unsupported platform: Linux only"
  [ -f /etc/os-release ] || fail "cannot detect operating system: /etc/os-release not found"

  # shellcheck disable=SC1091
  . /etc/os-release

  os_id=$(printf '%s' "${ID:-}" | tr '[:upper:]' '[:lower:]')
  os_like=$(printf '%s' "${ID_LIKE:-}" | tr '[:upper:]' '[:lower:]')

  case "$os_id" in
    ubuntu)
      return 0
      ;;
  esac

  case " $os_like " in
    *" ubuntu "*|*ubuntu*)
      return 0
      ;;
  esac

  fail "unsupported OS: this installer only supports Ubuntu and Ubuntu-based distributions (detected ID='${ID:-unknown}', ID_LIKE='${ID_LIKE:-unknown}')"
}

detect_arch() {
  arch_raw=$(uname -m)
  case "$arch_raw" in
    x86_64|amd64)
      arch="amd64"
      ;;
    aarch64|arm64)
      arch="arm64"
      ;;
    *)
      fail "unsupported architecture: $arch_raw (supported: x86_64/amd64, aarch64/arm64)"
      ;;
  esac
}

download_binary() {
  asset="${BIN_NAME}-linux-${arch}"
  version="$requested_version"

  if [ -n "$version" ]; then
    case "$version" in
      v*)
        resolved_version="$version"
        ;;
      *)
        resolved_version="v$version"
        ;;
    esac

    download_url="https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/download/${resolved_version}/${asset}"
    info "downloading ${BIN_NAME} ${resolved_version} (${arch})"
  else
    download_url="https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/latest/download/${asset}"
    info "downloading ${BIN_NAME} latest release (${arch})"
  fi

  tmp_file=$(mktemp)
  checksums_tmp_file=$(mktemp)
  trap 'rm -f "$tmp_file" "$checksums_tmp_file"' EXIT INT TERM

  if ! curl -fL "$download_url" -o "$tmp_file"; then
    fail "failed to download release asset '${asset}'. Ensure a GitHub release contains this binary."
  fi

  if [ -n "$version" ]; then
    checksums_url="https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/download/${resolved_version}/checksums.txt"
  else
    checksums_url="https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/latest/download/checksums.txt"
  fi

  info "downloading checksums.txt"
  if ! curl -fL "$checksums_url" -o "$checksums_tmp_file"; then
    fail "failed to download checksums.txt. Ensure the release publishes checksums.txt."
  fi
}

verify_checksum() {
  need_cmd sha256sum

  asset="${BIN_NAME}-linux-${arch}"
  expected=$(grep "  ${asset}$" "$checksums_tmp_file" | awk '{print $1}')
  [ -n "$expected" ] || fail "checksum entry not found for ${asset} in checksums.txt"

  actual=$(sha256sum "$tmp_file" | awk '{print $1}')
  [ "$actual" = "$expected" ] || fail "checksum verification failed for ${asset}"

  info "checksum verified for ${asset}"
}

install_binary() {
  installed_path=""

  if [ -w "$SYSTEM_BIN_DIR" ] || [ "$(id -u)" -eq 0 ]; then
    install -m 0755 "$tmp_file" "${SYSTEM_BIN_DIR}/${BIN_NAME}"
    installed_path="${SYSTEM_BIN_DIR}/${BIN_NAME}"
    return 0
  fi

  if command -v sudo >/dev/null 2>&1; then
    info "trying to install to ${SYSTEM_BIN_DIR} using sudo"
    if sudo install -m 0755 "$tmp_file" "${SYSTEM_BIN_DIR}/${BIN_NAME}"; then
      installed_path="${SYSTEM_BIN_DIR}/${BIN_NAME}"
      return 0
    fi
    warn "sudo install failed, falling back to user install"
  fi

  mkdir -p "$LOCAL_BIN_DIR"
  install -m 0755 "$tmp_file" "${LOCAL_BIN_DIR}/${BIN_NAME}"
  installed_path="${LOCAL_BIN_DIR}/${BIN_NAME}"
}

print_success() {
  "$installed_path" path >/dev/null 2>&1 || fail "installation completed but smoke check failed"

  info "installed ${BIN_NAME} at ${installed_path}"

  case ":$PATH:" in
    *":$LOCAL_BIN_DIR:"*)
      ;;
    *)
      if [ "$installed_path" = "${LOCAL_BIN_DIR}/${BIN_NAME}" ]; then
        printf '\nAdd this to your shell profile if needed:\n'
        printf '  export PATH="$HOME/.local/bin:$PATH"\n'
      fi
      ;;
  esac

  printf '\nRun this now:\n'
  printf '  cs list\n'
}

main() {
  parse_args "$@"

  need_cmd uname
  need_cmd curl
  need_cmd mktemp
  need_cmd install
  need_cmd grep
  need_cmd awk

  detect_os
  need_cmd apt-get
  need_cmd dpkg
  detect_arch
  download_binary
  verify_checksum
  install_binary
  print_success
}

main "$@"
