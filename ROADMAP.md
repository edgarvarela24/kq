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

### Phase 3: Pod Selection (current)
- List pods in a namespace
- Fuzzy-find UI with promptui
- Wire into `kq pods` command

### Phase 4: Pod Actions
- Logs (stream pod logs)
- Exec (shell into pod)

### Phase 5: Polish
- `--dry-run` flag (show kubectl equivalent)
- `--namespace` flag (skip namespace prompt)
- Better error messages
- Helpful `--help` text

### Stretch Goals
- Additional resource types (deployments, services)
- Port-forward action
- Describe action
- Config file support

---

*This roadmap is intentionally light. We'll reassess as we go.*
