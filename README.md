# pm

A superlight process manager for running named layouts of services. Define your processes once, run any combination by name.

## Installation

```sh
brew install rusintez/tap/pm
```

Or build from source (requires Go 1.22+):

```sh
git clone https://github.com/rusintez/pm
cd pm && go build -o pm . && mv pm /usr/local/bin/pm
```

## Configuration

Create a `pm.toml` in your project root:

```toml
[processes]
api    = { cmd = "bun run src/api/index.ts" }
db     = { cmd = "docker compose up db" }
ui     = { cmd = "bun run --hot src/ui/index.ts" }
worker = { cmd = "bun run src/worker.ts", env = { QUEUE = "default" } }

[layouts]
backend = { processes = ["api", "db"] }
dev     = { processes = ["api", "db", "ui"], env = { NODE_ENV = "development" } }
full    = { processes = ["api", "db", "ui", "worker"], env = { NODE_ENV = "production" } }
```

Layout `env` is merged with process `env`. Process `env` wins on conflicts.

## Usage

```sh
# Run a named layout
pm run dev
pm run backend

# Run an ad-hoc combination
pm run api db
pm run api db ui

# Run in tmux (one window per process, session named after the layout)
pm run dev --tmux

# Attach to a running tmux session
pm attach dev

# Stop a tmux session
pm stop dev

# List defined processes and layouts
pm list
```

### Default mode

Processes run in the foreground. Output is prefixed and color-coded per process:

```
[api]    Listening on http://localhost:3000
[db]     database system is ready to accept connections
[ui]     ready in 212ms
```

Ctrl-C sends SIGTERM to all processes and waits for clean exit.

### Tmux mode

Each process gets its own tmux window. The session is named after the layout (or `pm` for ad-hoc runs). Sessions survive terminal closes — reattach any time with `pm attach <layout>`.

## License

MIT
