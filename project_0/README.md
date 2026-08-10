### Netsim-0
A sandbox environment to create a simple network topology.

The general idea is to simulate a domain like so:

```text
Topology
 ├── Nodes
 │    ├── Node
 │    │    └── Interfaces
 │    └── Node
 │         └── Interfaces
 │
 └── Links
      └── interface <----> interface
```

The application simulates the entire domain in-memory in order to focus on the fundamental concepts.