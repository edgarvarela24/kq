# kq - Kubernetes Quick Actions

An interactive CLI companion for kubectl that reduces friction between "I want to do something" and actually doing it.

Instead of memorizing kubectl flags and resource names, you fuzzy-find your way to the right resource and action.

## Features

- **Interactive Selection**: Fuzzy-find namespaces, pods, and containers
- **Pod Logs**: Stream logs with follow, timestamps, and previous container support
- **Dry Run Mode**: See the equivalent kubectl command before executing
- **kubectl Plugin**: Use as `kubectl kq` or standalone `kq`

## Installation

### From Source

```bash
go install github.com/edgarvarela24/kq@latest
```

### As kubectl Plugin

```bash
# Clone and build
git clone https://github.com/edgarvarela24/kq.git
cd kq
make build-plugin

# Copy to your PATH
cp kubectl-kq ~/bin/  # or /usr/local/bin/

# Verify
kubectl kq --version
```

### Krew (coming soon)

```bash
kubectl krew install kq
```

## Usage

### Interactive Mode

```bash
# Full interactive flow: select namespace → pod → action → options
kq pods

# Interactive log viewing: select namespace → pod → container → options
kq logs
```

### Direct Mode

```bash
# Logs with flags (skips interactive prompts for provided values)
kq logs -n default my-pod -f              # Follow logs
kq logs -n default my-pod --timestamps    # With timestamps
kq logs -n default my-pod -p              # Previous container
kq logs -n default my-pod -c nginx        # Specific container

# Specify namespace, interactive for the rest
kq pods -n kube-system
kq logs -n default
```

### Dry Run

See the equivalent kubectl command without executing:

```bash
kq pods --dry-run
kq logs --dry-run -n default my-pod -f
```

## Commands

| Command | Description |
|---------|-------------|
| `kq pods` | Interactively select a pod and perform actions |
| `kq logs [pod]` | View pod logs (interactive or direct) |

## Flags

### Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--namespace` | `-n` | Kubernetes namespace |
| `--dry-run` | | Print kubectl command instead of executing |
| `--version` | `-v` | Print version |
| `--help` | `-h` | Help for kq |

### Logs Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--follow` | `-f` | Follow log output |
| `--timestamps` | | Show timestamps |
| `--previous` | `-p` | Show logs from previous container instance |
| `--container` | `-c` | Container name (for multi-container pods) |

## Development

```bash
# Clone
git clone https://github.com/edgarvarela24/kq.git
cd kq

# Build
make build

# Run tests
make test

# Create a test cluster with sample workloads
make cluster-create

# Clean up
make cluster-delete
```

## Future Plans

See [ROADMAP.md](ROADMAP.md) for planned features including:
- `exec` - Interactive shell into pods
- `describe` - View pod details
- `port-forward` - Forward local ports to pods
- Additional resource types (deployments, services)

## License

MIT
