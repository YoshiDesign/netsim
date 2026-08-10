### Netsim-0 - Version 0.1
A sandbox environment to create a simple network topology.

**Conceptually:**

```text
api/
    HTTP concerns

topology/
    What is a topology?
    What is a node?
    What is a link?

simulation/
    What does the topology DO?

store/
    Where does our current state live?
```

**Implementation:**

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

The application simulates its domain in-memory in order to focus on fundamental networking concepts and API design.

** Edge Traversal Constraints:**
- Node must be RUNNING
- Link must be UP
- Interface must be UP

**Project_0 TODO's and Omissions:**
- Synchronization + Backpressure for concurrent design
- Storage infra (mysql, redis, etc.)
- Request Validation
- Cancellation
- Operation IDs
- Timeouts
- Contexts
- Concurrency
