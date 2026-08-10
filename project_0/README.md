### Netsim-0 - Version 0.1
A sandbox environment to create a simple network topology. Most of the code will be built from Go's standard library, instead of using packages like `fiber` or `gin`.

#### Project 0 Checklist:
[] CRUD-ish topology API — create a topology with POST, retrieve it, list topologies, and delete one. This teaches request-body decoding, API models, status codes, path parameters, validation, and our in-memory store.
[] Nodes — add/remove simulated routers, switches, hosts, etc. Nodes get IDs and lifecycle states such as Stopped, Starting, and Running.
[] Links and interfaces — connect nodes together through interfaces and validate nonsense such as linking to nonexistent nodes or reusing an occupied interface.
[] Topology lifecycle/orchestration — POST /topologies/{id}/start and /stop, initially synchronously and then probably with some asynchronous lifecycle behavior. This is where context, goroutines, synchronization, and backend orchestration concepts start entering.
[] Actual simulation behavior — treat the topology as a graph and implement something like GET /path?from=A&to=B. A BFS determines reachability while respecting node state, link state, and eventually interface state.
[] Failure manipulation — bring a link down, stop a node, query reachability again, and watch the simulated network respond appropriately.
[] Clean backend architecture and tests — by the end, the HTTP handlers shouldn't own the domain logic. We'll have something closer to API → Service/Orchestrator → Simulation → Store, with unit tests around the important pieces.

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
- Storage infra (mysql, redis, etc.)
- Request Authorization
- Cancellation
- Operation IDs
- Timeouts
- Contexts
- Concurrency
