---
name: go-realtime-websocket-service
description: Design a horizontally scalable real-time backend in Go — a gorilla/websocket connection hub, Redis pub/sub to fan messages across multiple server instances, JWT auth with Redis-backed token revocation and presence, and a layered router → middleware → controller → service → store request flow. Use when building chat, notifications, or live features, scaling WebSockets beyond a single instance, adding presence or typing indicators, or reviewing a real-time Go service for correctness and scale.
---

# Go real-time WebSocket service

Modeled on the **ECHO** chat backend (Go, Gin, gorilla/websocket, Redis, MongoDB, JWT).
A single-instance hub is easy; the real work is staying correct when you run several
instances behind a load balancer.

## Architecture

```
REST:   router → middleware (JWT, CORS, logging) → controller → service → store
Live:   client ⇄ WebSocket hub ⇄ Redis pub/sub ⇄ other instances
```

- **HTTP layer** (Gin): auth, profiles, contacts, chats, groups, messages.
- **WebSocket hub**: owns the set of live connections and routes messages between them
  (typing indicators, read receipts, reactions).
- **Redis pub/sub**: the fan-out bus so instances share rooms.
- **Store**: MongoDB for durable data; Redis for tokens/presence.

## The hub (single instance)

Centralize connection lifecycle in one goroutine to avoid races. Each connection gets a
buffered send channel and a write pump; the hub never writes to a socket directly.

```go
type Hub struct {
    register   chan *Client
    unregister chan *Client
    rooms      map[string]map[*Client]bool
}

func (h *Hub) run() {
    for {
        select {
        case c := <-h.register:   h.rooms[c.room][c] = true
        case c := <-h.unregister: delete(h.rooms[c.room], c); close(c.send)
        }
    }
}
```

- One reader goroutine and one writer goroutine per connection.
- Never call `conn.WriteMessage` from two goroutines — gorilla/websocket forbids concurrent
  writes. Serialize through the client's `send` channel.
- Heartbeat with ping/pong and a read deadline to drop dead connections.

## Scaling across instances with Redis pub/sub

With N instances, the sender and recipient may be on different instances. Publish every
outbound message to Redis; each instance subscribes and delivers to its own local clients.

```go
// on send: persist, then publish to the room's channel
rdb.Publish(ctx, "room:"+roomID, payload)

// each instance runs one subscriber that fans in to its local hub
sub := rdb.Subscribe(ctx, "room:*")
for msg := range sub.Channel() {
    hub.deliverLocal(msg.Payload) // only to clients connected HERE
}
```

Tag each message with an origin instance ID so an instance can skip re-delivering what it
just published locally.

## Auth, revocation, presence

- **JWT** on the WebSocket upgrade request (query param or `Authorization`), validated in
  middleware before the hub accepts the connection.
- **Revocation**: JWTs are stateless, so keep a Redis denylist of revoked token IDs (`jti`)
  and check it in middleware — logout/ban takes effect immediately.
- **Presence**: a Redis key per online user with a TTL refreshed by heartbeat; publish
  join/leave so other instances update presence.

## Checklist

- [ ] Exactly one writer per connection; all writes go through a buffered `send` channel.
- [ ] Read deadline + ping/pong so dead peers are reaped.
- [ ] Every outbound message is published to Redis, not just sent locally.
- [ ] Messages carry an origin instance ID to prevent double delivery.
- [ ] JWT verified on upgrade; a Redis denylist enforces revocation.
- [ ] Slow-consumer policy: bounded `send` buffer; drop or disconnect when it overflows.
- [ ] Contract tests cover REST endpoints and the WebSocket message protocol.

## Pitfalls

- Concurrent writes to one gorilla/websocket connection → corrupted frames or panics.
- Broadcasting only to local clients → messages vanish the moment you scale past one instance.
- Unbounded per-client buffers → one slow client balloons memory; bound and shed load.
- Trusting a stateless JWT with no denylist → revoked/rotated tokens keep working until expiry.

## Backend relevance

Real-time fan-out, horizontal scale via a message bus, stateful-connection management, and
token security are core backend concerns. A service that stays correct across multiple
instances is a strong, concrete signal of backend capability.

## References

- [Using Redis for Go WebSocket scaling](https://oneuptime.com/blog/post/2026-03-31-redis-use-redis-for-go-websocket-scaling/view)
- [Scaling pub/sub with WebSockets and Redis](https://ably.com/blog/scaling-pub-sub-with-websockets-and-redis)
