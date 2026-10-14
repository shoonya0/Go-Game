---
name: go-clean-architecture
description: Structure a Go module for maintainability and testability — cmd/ entrypoints, internal/ private packages organized by domain, small interfaces defined in the consuming package, and constructor (dependency) injection. Use when starting a new Go module, laying out packages, deciding where an interface belongs, breaking an import cycle, or reviewing Go code for tight coupling and poor testability.
---

# Go clean architecture & project layout

Keep dependencies pointing inward and keep packages honest about what they expose. The goal
is code that is easy to test in isolation and cheap to change.

## Layout

```
cmd/            # entrypoints (one dir per binary); wires everything together
internal/       # private application code — cannot be imported by other modules
  core/         # domain: types + logic (physics, world, player)
  system/       # adapters to the outside (input, IO)
server/         # a second binary (the static file server)
assets/         # embedded resources
```

- `cmd/` holds `main`. It constructs concrete dependencies and injects them. Keep logic out.
- `internal/` is enforced by the compiler as private to this module — a real boundary, not a
  convention. Put everything that is not a deliberately public API here.
- Organize `internal/` **by domain, not by technical layer**. Prefer `internal/order/`
  (its handler, service, store together) over `internal/handlers/` + `internal/services/`.
  Avoid deep nesting like `internal/services/user/handlers/http/v1/`.

## Interfaces: define them where they are used

Define an interface in the **consumer** package, listing only the methods that consumer
needs (Interface Segregation). The producer returns a concrete type.

```go
// in the package that DOES collision — it only needs bounds.
type Collider interface {
    GetBounds() AABB
}
```

Tiles and entities both satisfy `Collider` without knowing the interface exists, so one
collision path serves both. This is how you make handlers, stores, and clients swappable and
mockable behind a boundary.

## Dependency injection via constructors

Pass dependencies in; do not reach for globals or construct them inside business logic.

```go
func UpdatePlayer(p *PlayerRuntime, in *InputState, qt *DynamicQuadtree) { ... }
```

`main` builds the quadtree, player, and images and threads them through. For a service:

```go
type OrderService struct { store OrderStore; clock Clock }
func NewOrderService(s OrderStore, c Clock) *OrderService { return &OrderService{s, c} }
```

Now `OrderService` is unit-testable with a fake `OrderStore` and a fixed `Clock`.

## Checklist

- [ ] `main` is thin: it wires dependencies and starts the app; no business rules.
- [ ] Private code lives under `internal/`; only genuinely public APIs sit outside it.
- [ ] Packages are grouped by domain/feature, not by `handlers/services/repositories`.
- [ ] Interfaces are small and declared in the consumer, not the implementer.
- [ ] Dependencies arrive through constructors/parameters, not package-level globals.
- [ ] Import graph is acyclic; break cycles with a shared types package or a consumer interface.

## Pitfalls

- Defining fat interfaces in the implementing package "just in case" — it couples consumers to
  methods they never call and invites import cycles.
- A `utils`/`common` package that everything imports and that slowly becomes a dependency
  magnet.
- Globals and `init()`-time singletons that make tests order-dependent.

## Backend relevance

This is the difference between a service you can grow and one that ossifies: clear package
boundaries, injected dependencies, and small interfaces are what keep handlers, business
logic, and storage independently testable and replaceable.

## References

- [Go project structure & clean architecture](https://reintech.io/blog/go-project-structure-2026-clean-architecture-best-practices)
- [Dependency injection in Go: patterns & best practices](https://www.glukhov.org/post/2025/12/dependency-injection-in-go/)
