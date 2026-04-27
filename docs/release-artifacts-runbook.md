# Runbook: Publish GitHub Release Artifacts

This runbook describes how to publish release artifacts required by `install.sh`.

## Automated release (recommended)

This repository includes a workflow at `.github/workflows/release.yml` that automatically builds and uploads release assets when you push a tag like `v0.1.0`.

Minimal automated flow:

```bash
VERSION="0.1.0"
TAG="v${VERSION}"
git tag "${TAG}"
git push origin "${TAG}"
```

The workflow publishes:

- `cs-linux-amd64`
- `cs-linux-arm64`
- `checksums.txt`

Use the manual steps below only if you need to publish without CI.

## Required assets

Each release must include these files with exact names:

- `cs-linux-amd64`
- `cs-linux-arm64`
- `checksums.txt`

## 1) Prepare version

```bash
VERSION="0.1.0"
TAG="v${VERSION}"
```

## 2) Build Linux binaries

```bash
GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o cs-linux-amd64 ./cmd/cs
GOOS=linux GOARCH=arm64 go build -trimpath -ldflags="-s -w" -o cs-linux-arm64 ./cmd/cs
```

## 3) Generate checksums file

```bash
sha256sum cs-linux-amd64 cs-linux-arm64 > checksums.txt
```

## 4) Verify artifacts

```bash
ls -lh cs-linux-amd64 cs-linux-arm64 checksums.txt
```

## 5) Tag and push

```bash
git tag "${TAG}"
git push origin "${TAG}"
```

## 6) Create release and upload artifacts

```bash
gh release create "${TAG}" \
  --title "${TAG}" \
  --notes "Release ${TAG}" \
  ./cs-linux-amd64 \
  ./cs-linux-arm64 \
  ./checksums.txt
```

## 7) Verify uploaded files

```bash
gh release view "${TAG}" --json assets --jq '.assets[].name'
```

Expected output includes:

- `cs-linux-amd64`
- `cs-linux-arm64`
- `checksums.txt`

## 8) Test installer

Latest:

```bash
curl -fsSL https://raw.githubusercontent.com/fernando-pires47/cli-catalog-go/main/install.sh | sh
```

Pinned version:

```bash
curl -fsSL https://raw.githubusercontent.com/fernando-pires47/cli-catalog-go/main/install.sh | sh -s -- --version "${VERSION}"
```

## Security notes

- Do not publish tokens or secrets in docs, scripts, commit history, or remote URLs.
- Rotate any token immediately if it appears in a command output or config.
- Keep release assets immutable after publishing.
