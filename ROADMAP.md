# kq Roadmap

Current status: **v0.1.0** - Pod logs functionality complete.

## Current Features

- [x] Interactive namespace selection
- [x] Interactive pod selection with fuzzy filtering
- [x] Pod logs streaming (follow, timestamps, previous)
- [x] Container selection for multi-container pods
- [x] Dry-run mode to show equivalent kubectl commands
- [x] Works as standalone CLI or kubectl plugin

## Planned Features

### Pod Actions
- [ ] `exec` - Interactive shell into a container
- [ ] `describe` - View detailed pod information  
- [ ] `port-forward` - Forward local ports to a pod
- [ ] `delete` - Delete a pod (with confirmation)

### Additional Resources
- [ ] Deployments - Scale, restart, view status
- [ ] Services - View endpoints, port-forward
- [ ] ConfigMaps/Secrets - View and edit
- [ ] Nodes - View status, cordon/drain

### UX Improvements
- [ ] Multi-select for batch operations
- [ ] Saved favorites / recent selections
- [ ] Custom themes
- [ ] Config file for defaults

---

Contributions welcome! Focus is on polished, well-tested features over quantity.
