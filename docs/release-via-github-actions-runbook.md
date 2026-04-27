# Runbook: Publish Release Artifacts via GitHub Actions

This runbook covers only the CI publishing flow. Release assets are built and uploaded by `.github/workflows/release.yml`.

## Scope

When a tag matching `v*` is pushed (example: `v0.1.0`), GitHub Actions:

- builds Linux binaries
- generates checksums
- publishes release assets

Expected assets:

- `cs-linux-amd64`
- `cs-linux-arm64`
- `checksums.txt`

## 1) Prerequisites

- You have push permission to the repository.
- GitHub Actions is enabled for the repository.
- Workflow file exists at `.github/workflows/release.yml`.

## 2) Publish a release

```bash
VERSION="0.1.0"
TAG="v${VERSION}"
git tag "${TAG}"
git push origin "${TAG}"
```

This push triggers the release workflow automatically.

## 3) Verify workflow success

- Open the repository Actions tab.
- Confirm the `Release` workflow run for your tag completed successfully.

## 4) Verify release assets

```bash
gh release view "${TAG}" --json assets --jq '.assets[].name'
```

Expected output includes:

- `cs-linux-amd64`
- `cs-linux-arm64`
- `checksums.txt`

## 5) Verify installer

Pinned version:

```bash
curl -fsSL https://raw.githubusercontent.com/fernando-pires47/cli-catalog-go/main/install.sh | sh -s -- --version "${VERSION}"
```

Latest:

```bash
curl -fsSL https://raw.githubusercontent.com/fernando-pires47/cli-catalog-go/main/install.sh | sh
```

## Troubleshooting

- Workflow did not start:
  - Confirm tag starts with `v`.
  - Confirm tag was pushed to `origin`.
- Workflow failed to upload assets:
  - Confirm workflow permissions include `contents: write`.
- Installer download fails:
  - Confirm release contains all required assets with exact names.

## Security notes

- Do not store or paste tokens in docs, commands, or remote URLs.
- Rotate credentials immediately if exposed.
