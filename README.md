# command-catalog-cli

Local CLI (`cs`) to create, list, delete, and execute reusable shell command templates.

## Quick start

```bash
make build
make run ARGS="list"
```

## Makefile shortcuts

```bash
make help
make build
make run ARGS="list"
make test
make install
make clean
```

Notes:
- `make run` always builds first, then runs `./cs` with `ARGS`.
- `make install` copies the binary to `~/.local/bin/cs`.

## Usage

```bash
./cs create 'kill port' 'sudo kill -9 $(sudo lsof -t -i:$port)' dangerous=yes
./cs create 'kp $port' 'sudo kill -9 $(sudo lsof -t -i:$port)'
./cs create 'logs $ns $lines' 'kubectl logs deployment/api -n $ns --tail=$lines'
./cs list
./cs kill port 3040
./cs kp 3040
./cs logs prod 200
./cs path
./cs delete <id>
```

### Catalog file resource: `CS_CATALOG_PATH`

By default, `cs` stores commands in `$HOME/.cs/catalog.json`.
Set `CS_CATALOG_PATH` to point to a different catalog file.

For one command only:

```bash
CS_CATALOG_PATH="$HOME/my-project/.cs/catalog.json" ./cs list
```

To export for your current shell session:

```bash
export CS_CATALOG_PATH="$HOME/my-project/.cs/catalog.json"
```

To check which catalog file `cs` is currently using:

```bash
./cs path
```

Details:
- If `CS_CATALOG_PATH` is set and non-empty, `cs` uses it as-is.
- If the file does not exist yet, `cs` treats the catalog as empty until you create your first command.
- Parent directories are created automatically when saving.
- If the file contains invalid JSON, `cs` fails with `invalid catalog json` and includes the file path in the error message.

Enable local debug hooks with `CS_DEBUG=1`.
You can also enable debug logs for a single invocation with `--debug` (example: `./cs --debug list`).
When enabled, the CLI emits debug events to stderr (`catalog_loaded`, `command_created`, `command_deleted`, `match_resolved`, `danger_confirmation_prompted`, `command_executed`).

## Dangerous command confirmation

Commands created with `dangerous=yes` require confirmation via `[y/N]` before execution.
In non-interactive mode, dangerous commands fail safely without execution.

Create syntax supports optional explicit safety flag:

```bash
./cs create '<key>' '<value>' dangerous=yes
./cs create '<key>' '<value>' dangerous=no
```

`./cs list` includes a `dangerous` column with `yes`/`no` values.
