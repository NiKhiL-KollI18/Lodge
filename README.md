# LODGE

A lightweight, payload-agnostic, room-based server written in Go.

LODGE is a simple and scalable server built around a room-oriented architecture. It is designed to handle arbitrary data without enforcing structure, making it flexible for signaling, real-time systems, and other client-server applications.

---

## Features

### Room-Based Server Abstraction

LODGE manages rooms and their lifecycle internally, allowing clients to interact with the system through simple HTTP requests.

```json
{
  "type": "create",
  "guest_count": 4
}
```

The server handles room creation, membership, and constraints, so you can focus on building your application logic.

---

### Payload-Agnostic Messaging

LODGE does not enforce any payload structure.

It only interprets a minimal set of system-level message types:

* `create`
* `join`
* `leave`

All other data is passed directly between peers within a room, making the system highly flexible and adaptable to different use cases.

---

### Built-in Room Lifecycle Management ("Room Service")

LODGE includes a configurable cleanup system that periodically scans and manages rooms.

This "room service" can:

* Remove inactive or empty rooms
* Handle dangling connections
* Enforce stay limits and timeouts

---

### Scalable by Design

LODGE is built with scalability in mind, leveraging Go's concurrency model (goroutines and channels) to efficiently handle large numbers of connections and rooms.

---

## Getting Started

### Prerequisites

* Go installed on your system

---

### Option 1: Use as a Library

```bash
git clone https://github.com/NikhiL-Kolli18/Lodge/lib
```

Import in your project:

```go
import (
    "LODGE/lib"
)
```

---

### Option 2: Run the Full Server

```bash
git clone https://github.com/NiKhiL-Kolli18/LODGE
cd LODGE
go build -o lodge cmd/server/main.go
./lodge
```

---

### Default Configuration (Standalone Server)

If you use the provided server (`main.go`), LODGE comes with sensible defaults:

* **Port**: `8080` (used if not assigned by the environment)

#### Room Service Configuration

* **Inactive Timeout**: 3 minutes
* **Waiting Timeout**: 10 minutes
* **Max Lifecycle**: Unlimited
* **Cleanup Interval**: Every 30 seconds
* **Max Room Capacity**: 100 rooms

You can modify these values directly in the server code.

Example:

```go
UpdateGlobalMaxCapacity(200)
```

---

### Notes

The standalone server is intentionally minimal and opinionated.

If your use case requires deeper control or custom behavior, it is recommended to use LODGE as a library.

---

## Usage

All interactions with LODGE are performed via JSON requests.

### Request Fields

* `type` (string) — Defines the action or message type (required)
* `room_id` (string) — Target room identifier
* `guest_count` (int) — Maximum room capacity (optional)
* `data` (string) — Payload forwarded to other peers

---

## Basic Flow

### 1. Create a Room

```json
{
  "type": "create",
  "guest_count": 4
}
```

Response:

```json
{
  "type": "room created",
  "room_id": "XfRs12"
}
```

> If `guest_count` is not provided, the global maximum capacity is used.

---

### 2. Join a Room

```json
{
  "type": "join",
  "room_id": "XfRs12"
}
```

**Possible responses:**

Room full:

```json
{
  "type": "error",
  "data": "room is full"
}
```

Room not found:

```json
{
  "type": "error",
  "data": "room not found"
}
```

Success:

* HTTP Status: `200 OK`

---

### 3. Leave a Room

```json
{
  "type": "leave"
}
```

Response:

* HTTP Status: `200 OK`

---

### 4. Broadcast Data

```json
{
  "type": "myusecase",
  "data": "my data"
}
```

Behavior:

* Non-system `type` values are treated as application-level messages
* Payload is not interpreted or modified
* Messages are broadcast to all peers in the room

---

### Notes

* `type` is mandatory
* Payload is ignored for `create`, `join`, `leave`
* `room_id` is a randomly generated 6-character string (A–Z, a–z, 0–9)
* LODGE does not enforce any payload structure

---

## Using LODGE as a Library

### Waitress (Request Handler)

Primary HTTP handler that processes and routes requests.

```go
func Waitress(w http.ResponseWriter, r *http.Request)
```

Example:

```go
http.HandleFunc("/request", lib.Waitress)
```

---

### Room Service (Lifecycle Management)

Deploy and start the room cleanup system:

```go
roomService := lib.DeployRoomService(
    60 * time.Second,
    10 * time.Minute,
    0,
    30 * time.Second,
)
roomService.Start()
```

---

#### Behavior

* 0 guests → removed
* 1 guest → removed after waiting timeout
* > 1 guests → removed after inactivity timeout
* Any room → removed after max lifecycle

---

#### Notes

* `0` duration = infinite
* Must be started manually
* Independent of request handling

---

## Use Cases

LODGE is designed for flexibility and can be used in:

* WebRTC signaling
* Multiplayer lobbies
* Real-time chat systems
* Peer-to-peer coordination
* Event-driven systems

Because it is payload-agnostic, it can adapt to custom protocols and workflows.

---

## Philosophy

LODGE was originally built as a WebRTC signaling server, and later adapted into a payload-agnostic system for simplicity and flexibility.

The goal throughout development was straightforward: give the user (me) control where it matters, while abstracting away parts that usually don’t need to be touched.

The naming follows the same idea.
It’s meant to feel like a small, friendly lodge — something that accepts requests without asking too many questions, but still has its own rules (and will politely show you the door when your time is up).

---

Is LODGE perfect? No — not even close.
Is it good? That depends on what you need.

If you want something easy to set up and run, it fits.
If you need something lightweight that can handle multiple peers and scale reasonably well, it fits there too.

---

It’s not built for any one kind of developer.

* Students building projects
* Individual developers integrating a backend
* Freelancers building systems

If it works for your use case, it works.

---

That’s really all there is to it.

If you find it useful, feel free to use it.

---

## License

MIT License
