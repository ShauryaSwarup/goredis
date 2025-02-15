---
# GoRedis - A Redis Clone in Go 🐭

GoRedis is a lightweight, in-memory key-value store inspired by Redis, implemented in Go. It supports a subset of Redis commands and the Redis Serialization Protocol (RESP).
---

## Features

- **In-Memory Storage**: Fast key-value storage with support for strings.
- **RESP Protocol**: Fully compatible with Redis Serialization Protocol (RESP).
- **Command Support**:
  - `PING`
  - `SET key value`
  - `GET key`
  - `DEL key`
- **Concurrency**: Handles multiple clients concurrently using goroutines.
- **Graceful Shutdown**: Properly closes connections and cleans up resources.

---

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/goredis.git
   cd goredis
   ```

2. Build and run the server:

   ```bash
   go run main.go
   ```

3. Connect using `redis-cli` or `netcat`:
   ```bash
   redis-cli -p 5001
   ```

---

## Usage

### Start the Server

```bash
make
```

### Connect with `redis-cli`

```bash
redis-cli -p 6379
```

### Example Commands

```redis
127.0.0.1:6379> PING
PONG

127.0.0.1:6379> SET foo bar
OK

127.0.0.1:6379> GET foo
"bar"

127.0.0.1:6379> DEL foo
(integer) 1
```

### Connect with `netcat`

```bash
echo -ne "*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n" | nc localhost 5001
```

---

## Architecture

### Components

1. **Server**: Handles client connections and command processing.
2. **Peer**: Manages individual client connections.
3. **RESP Parser**: Parses Redis Serialization Protocol (RESP) commands.

### Concurrency Model

- Each client connection is handled in a separate goroutine.
- Commands are processed sequentially in a single-threaded event loop.

---

## RESP Protocol Support

GoRedis implements the following RESP types:

- **Simple Strings**: `+OK\r\n`
- **Errors**: `-ERR message\r\n`
- **Integers**: `:1000\r\n`
- **Bulk Strings**: `$5\r\nhello\r\n`
- **Arrays**: `*2\r\n$3\r\nGET\r\n$3\r\nkey\r\n` (Non-homogeneous and Nested arrays also supported)

---

## Development

### Running RESP Test

```bash
cd resp
go test -v
```

---

## Improvements to be made

- [ ] Add support for more Redis commands (`EXISTS`, `INCR`, etc.).
- [ ] Implement TTL (Time-to-Live) for keys.
- [ ] Add persistence (AOF/RDB).
- [ ] Support Redis replication.

---

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

---

## Acknowledgments

- Inspired by [Redis](https://redis.io/).
- Built using Go.

---
