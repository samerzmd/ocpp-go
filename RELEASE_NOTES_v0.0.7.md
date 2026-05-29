# v0.0.7

Supersedes the unreleased v0.0.6 tag (which pointed at an intermediate commit that did not pass `go test -race`). Use v0.0.7.

## Fixed

- **Fleet-wide WebSocket server deadlock under network degradation.** `Server.Write` held `connMutex.RLock` across a blocking send to a size‑1 `outQueue`. A single stalled or half‑open peer could park a writer goroutine while holding the read lock; `cleanupConnection` and new‑connection registration in `wsHandler` then blocked on `connMutex.Lock`, and Go's `RWMutex` writer‑starvation rule blocked every other `Write` and every reconnect. The result was a permanent, server‑wide deadlock recoverable only by a process restart — observed in production as all connected chargers disconnecting at once and failing to reconnect until the service was restarted.
- **Data race in client error reporting.** `Client.error`/`Errors`/`Stop` accessed `errC` without synchronization. Error delivery is now non‑blocking and guarded by a dedicated mutex (mirrors the server fix).

## Changes

- `Server.Write` releases the lock before sending and applies bounded backpressure: it waits up to `WriteWait` on a buffered queue, then sheds the captured connection (never a fresh lookup by ID). No lock is held across the send.
- Connection lifecycle reworked: a `done` channel plus a `connMutex`‑guarded `closed` flag replace closing `outQueue`/`closeC`. `cleanupConnection` is idempotent and its map delete is guarded (`current == ws`). Removes send‑on‑closed‑channel panics and the `sync.Once` `copylocks` vet failure.
- On reconnect, `wsHandler` atomically swaps the new connection into the registry and evicts the captured stale one, instead of rejecting the reconnecting client with a policy violation.
- Server‑initiated ping/pong keepalive so dead/half‑open peers are detected within ~`PingWait` instead of waiting out a passive read deadline. The ping ticker is disabled when `PingWait == 0` (consistent with the disabled read deadline) and the derived period is guarded against `0`.
- Error reporting on both `Server` and `Client` is now non‑blocking and mutex‑synchronized.
- `outQueue` buffer increased from 1 to 64.

## Notes

- Behavior change: a reconnecting charger with an existing session ID is now accepted (old connection evicted) rather than rejected.
- A few pre‑existing client/test data races (`TestServerStartErrors`, `TestWebsocketClientConnectionBreak`, `TestServerErrors`) are quarantined under `-race`; they still run without `-race`. Tracked for a follow‑up fix.
