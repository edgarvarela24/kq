# kq Roadmap

A flexible guide for building `kq` — an interactive Kubernetes CLI.

## Phases

### Phase 1: Project Skeleton ✅
- Go module setup
- Cobra CLI framework
- Basic `kq` and `kq pods` commands

### Phase 2: Kubernetes Client ✅
- Connect to cluster via kubeconfig
- List namespaces
- Unit tests with fake clientset

### Phase 3: Pod Selection ✅
- List pods in a namespace
- Interactive selection with filtering (huh)
- Wire into `kq pods` command

### Phase 4: Logs Action (current)
- `kq logs` subcommand for power users
- Action menu after pod selection for interactive flow
- Options: follow, timestamps, previous, container selection
- `--dry-run` flag to show kubectl equivalent

### Phase 5: Polish & Ship
- Good error messages
- README with usage examples
- Demo GIF
- Maybe: `exec` action

### Stretch Goals (post-ship)
- Additional actions (exec, describe, port-forward)
- Additional resource types (deployments, services)
- Config file support

---

*Focus: One or two polished actions over many half-baked ones.*
