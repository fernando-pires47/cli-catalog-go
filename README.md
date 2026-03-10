# command-catalog-cli

Local CLI (`cs`) to create, list, delete, and execute reusable shell command templates.

## Quick start

```bash
make build
make run ARGS="list"
make test
```

## Build

```bash
make build
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
./cs create "kill port" 'sudo kill -9 $(sudo lsof -t -i:$port)'
./cs create 'kp $port' 'sudo kill -9 $(sudo lsof -t -i:$port)'
./cs create 'logs $ns $lines' 'kubectl logs deployment/api -n $ns --tail=$lines'
./cs list
./cs kill port 3040
./cs kp 3040
./cs logs prod 200
./cs delete <id>
```

Catalog path defaults to `$HOME/.cs/catalog.json` and can be overridden with `CS_CATALOG_PATH`.

Danger patterns can be extended with `CS_DANGER_PATTERNS` as a comma-separated list.
Example: `CS_DANGER_PATTERNS="terraform destroy,helm uninstall"`.

Enable local debug hooks with `CS_DEBUG=1`.
You can also enable debug logs for a single invocation with `--debug` (example: `./cs --debug list`).
When enabled, the CLI emits debug events to stderr (`catalog_loaded`, `command_created`, `command_deleted`, `match_resolved`, `danger_confirmation_prompted`, `command_executed`).

## Dangerous command confirmation

Commands matching baseline dangerous patterns (for example `rm -rf`) require confirmation via `[y/N]`.
In non-interactive mode, dangerous commands fail safely without execution.
