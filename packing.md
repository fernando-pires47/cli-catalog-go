# Packaging `command-catalog-cli` (`cs`) for `apt`

This is a project-specific guide for this repository (`cli-catalog-go`) so you can publish `cs` through both:

1. Ubuntu PPA (Launchpad)
2. Self-hosted APT repo (for Debian + Ubuntu)

Suggested Debian package name: `command-catalog-cli`  
Installed command: `cs`

---

## 0) One-time setup

Install packaging tools:

```bash
sudo apt update
sudo apt install -y build-essential devscripts debhelper dh-make lintian \
  fakeroot gnupg dput-ng golang-any
```

Create/import a GPG key (if you do not have one):

```bash
gpg --full-generate-key
gpg --list-secret-keys --keyid-format LONG
```

Set your key id for later commands:

```bash
export GPG_KEY_ID="YOUR_LONG_KEY_ID"
```

---

## 1) Add Debian packaging files to this repo

From repo root:

```bash
dh_make --createorig -s -p command-catalog-cli_0.1.0 -y
```

That creates `debian/`. Replace generated content with the files below.

### `debian/control`

```debcontrol
Source: command-catalog-cli
Section: utils
Priority: optional
Maintainer: Your Name <you@example.com>
Build-Depends: debhelper-compat (= 13), golang-any
Standards-Version: 4.6.2
Rules-Requires-Root: no
Homepage: https://github.com/<you>/cli-catalog-go

Package: command-catalog-cli
Architecture: any
Depends: ${misc:Depends}
Description: Local CLI catalog for reusable shell commands
 A local command catalog CLI that lets users create, list, delete,
 and execute reusable shell command templates.
```

### `debian/changelog`

```bash
dch --create --package command-catalog-cli -v 0.1.0-1 "Initial Debian package release"
```

### `debian/rules`

```make
#!/usr/bin/make -f

%:
	dh $@

override_dh_auto_build:
	go build -trimpath -ldflags="-s -w" -o cs ./cmd/cs

override_dh_auto_install:
	install -D -m 0755 cs debian/command-catalog-cli/usr/bin/cs
```

Make it executable:

```bash
chmod +x debian/rules
```

### `debian/source/format`

```text
3.0 (native)
```

### `debian/copyright`

Fill with your real copyright/license metadata.

---

## 2) Build and test locally

Build package:

```bash
debuild -us -uc
```

Lint it:

```bash
lintian ../command-catalog-cli_*.changes
```

Install and smoke test:

```bash
sudo apt install ../command-catalog-cli_*_amd64.deb
cs list
```

---

## 3) Publish Path A: Ubuntu PPA (Launchpad)

### 3.1 Launchpad setup

1. Create a Launchpad account.
2. Upload your GPG public key to Launchpad.
3. Create PPA, for example: `ppa:<launchpad-user>/command-catalog`.

### 3.2 Build source package

PPA uploads must be source packages:

```bash
debuild -S
```

### 3.3 Upload

```bash
dput ppa:<launchpad-user>/command-catalog ../command-catalog-cli_0.1.0-1_source.changes
```

### 3.4 User install (Ubuntu)

```bash
sudo add-apt-repository ppa:<launchpad-user>/command-catalog
sudo apt update
sudo apt install command-catalog-cli
```

---

## 4) Publish Path B: Self-hosted APT (Debian + Ubuntu)

### 4.1 Create repo with aptly

```bash
sudo apt install -y aptly
aptly repo create -distribution=stable -component=main command-catalog
aptly repo add command-catalog ../command-catalog-cli_*_amd64.deb
```

### 4.2 Publish and sign metadata

```bash
aptly publish repo -gpg-key="$GPG_KEY_ID" command-catalog
```

Published repo root: `~/.aptly/public/`

### 4.3 Host repo files over HTTPS

Host `~/.aptly/public/` at a domain you control, for example:

- `https://repo.yourdomain.com/`

### 4.4 Export signing key for clients

```bash
gpg --armor --export "$GPG_KEY_ID" > repo-public.asc
gpg --dearmor < repo-public.asc > repo-public.gpg
```

Upload `repo-public.gpg` to:

- `https://repo.yourdomain.com/repo-public.gpg`

### 4.5 User install (Debian + Ubuntu)

```bash
curl -fsSL https://repo.yourdomain.com/repo-public.gpg | \
  sudo gpg --dearmor -o /etc/apt/keyrings/command-catalog.gpg

echo "deb [signed-by=/etc/apt/keyrings/command-catalog.gpg] https://repo.yourdomain.com stable main" | \
  sudo tee /etc/apt/sources.list.d/command-catalog.list

sudo apt update
sudo apt install command-catalog-cli
```

---

## 5) Release process for this repo

For each release (`0.1.1`, `0.2.0`, etc):

1. Update package changelog:

   ```bash
   dch -v 0.1.1-1 "Release 0.1.1"
   ```

2. Build and lint:

   ```bash
   debuild -us -uc
   lintian ../command-catalog-cli_*.changes
   ```

3. Publish to one or both:
   - PPA: `debuild -S` then `dput ..._source.changes`
   - Self-hosted: `aptly repo add ...deb` then `aptly publish update stable`

4. Validate in clean environments:

   ```bash
   docker run --rm -it ubuntu:24.04 bash
   docker run --rm -it debian:12 bash
   ```

---

## 6) CI automation (recommended)

- Trigger on git tags.
- Build `.deb` artifact.
- Upload source package to Launchpad PPA.
- Add `.deb` to aptly repo and republish signed metadata.
- Run smoke install test that executes `cs list`.

---

## 7) Common gotchas

- PPA rejects binary-only uploads; always upload source package.
- Missing keyring setup on client systems causes signature errors.
- `apt-key` is deprecated; use `/etc/apt/keyrings`.
- Forgetting `debian/changelog` updates blocks new uploads.

---

## 8) Quick checklist

- [ ] `debian/control`, `debian/rules`, `debian/changelog` are valid
- [ ] `debuild` and `lintian` pass
- [ ] PPA upload and build succeed
- [ ] Self-hosted repo metadata is signed and reachable over HTTPS
- [ ] Install works with `apt install command-catalog-cli` on Ubuntu and Debian
