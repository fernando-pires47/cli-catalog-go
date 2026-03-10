# cs Runbook (Linux)

This runbook shows how to build `cs`, point it on linux, and run daily commands.

## 1) Build and install `cs`
```bash
make build
mkdir -p "$HOME/.local/bin"
rm -r "$HOME/.local/bin/cs"
cp ./cs "$HOME/.local/bin/cs"
chmod +x "$HOME/.local/bin/cs"
```

Ensure `~/.local/bin` is in your `PATH`:

```bash
echo 'export PATH="$HOME/.local/bin:$PATH"' >> "$HOME/.bashrc"
source "$HOME/.bashrc"
```

Verify:

```bash
cs --help
```

## 2) Point `cs` to your project catalog

Choose a catalog location inside your project (example below):

```bash
mkdir -p "$HOME/my-project/.cs"
export CS_CATALOG_PATH="$HOME/my-project/.cs/catalog.json"
```

Persist for new terminals:

```bash
echo 'export CS_CATALOG_PATH="$HOME/my-project/.cs/catalog.json"' >> "$HOME/.bashrc"
source "$HOME/.bashrc"
```

Tip: You can also set `CS_CATALOG_PATH` inside a project-specific shell script and source it when entering the project.

## 3) Core commands you can execute

Create command templates:

```bash
cs create 'kill port' 'sudo kill -9 $(sudo lsof -t -i:$port)' dangerous=yes
cs create 'logs api' 'kubectl logs deployment/api -n $ns --tail=$lines'
cs create 'kp $port' 'sudo kill -9 $(sudo lsof -t -i:$port)'
cs create 'logs $ns $lines' 'kubectl logs deployment/api -n $ns --tail=$lines'
```

List saved commands:

```bash
cs list
```

`cs list` shows a `dangerous` column with `yes` or `no`.

Run commands by key + runtime args:

```bash
cs kill port 3040
cs logs api prod 200
cs kp 3040
cs logs prod 200
```

Delete by id:

```bash
cs delete <id>
```

## 4) Safety and debug options

Enable debug logs (stderr):

```bash
cs --debug list
```

Or via env var:

```bash
export CS_DEBUG=1
cs list
```

Mark commands as dangerous at creation time:

```bash
cs create '<key>' '<value>' dangerous=yes
cs create '<key>' '<value>' dangerous=no
```

## 5) Quick troubleshooting

- `command not found: cs`
  - Confirm `~/.local/bin` is in `PATH` and re-open terminal.
- `invalid catalog json`
  - Fix or remove file at `CS_CATALOG_PATH`.
- Command requires confirmation but shell/script is non-interactive
  - Run interactively in terminal to confirm commands created with `dangerous=yes`.
