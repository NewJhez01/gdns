# GDNS — Block All Unwanted Traffic

A DNS resolver built from scratch in Go. Parses raw DNS messages per [RFC 1035](https://datatracker.ietf.org/doc/html/rfc1035), implements UDP socket networking, and blocks ads, trackers, and malware via a local SQLite blocklist with Redis-backed response caching. Designed to run on a Raspberry Pi as a Pi-hole alternative.

This project will parse incoming DNS requests and run them against a preselected list provided by [StevenBlack](https://github.com/StevenBlack/hosts); unwanted entries return an error and the request will not complete; valid requests are resolved upstream by Cloudflare at 1.1.1.1.

## Description

- **Binary DNS protocol parser** — manual bit-level parsing of headers, flags, and length-prefixed QNAME labels
- **UDP server** — concurrent request handling with `net.ListenUDP` and per-packet goroutines
- **Response construction** — NXDOMAIN and sinkhole response generation without external libraries
- **Caching layer** — Redis with JSON-serialized values and configurable TTL
- **Blocklist management** — SQLite with `WITHOUT ROWID` optimization and automatic migration on startup

## Getting Started

## Prerequisites

- Docker & Docker Compose
- Port 53/UDP available on the host (or configure an alternative in docker-compose.yml)

## Quick Start

```bash
git clone git@github.com:NewJhez01/gdns.git
cd gdns
docker compose up -d

```

This starts:

- Redis for response caching
- GDNS UDP server on port 53 (host), mapped to 5555 (container), with automatic SQLite migration on first boot

## Configure Your Network

Point your router's DNS to the host IP running GDNS:

```
Router DNS: 192.168.1.100 ← your host machine's IP. No port needed, 53 is implicit with DNS request
```

Or test locally:

```bash
dig @localhost example.com ← no port needed, 53 is implicit
```

Note: The container maps host port 53 → container port 5555. External clients (including your router) must use port 53, not 5555.

## Project Structure

```text
cmd/
  main.go                # Wire dependencies
  server/
    main.go              # Start UDP server, handle signals

internal/
  dns/
    resolver.go          # Orchestrate cache → blocklist → upstream
    upstream.go          # Forward to Cloudflare (1.1.1.1) with timeout orchestrate answer parsing
    parser/
      answer_parser.go   # Parses the answer recursively with the pointer
      header_parser.go   # 12-byte DNS header parser
      marshall.go        # Marshals the DNS response
      parser.go          # orchestrates the parsing of the incoming request
      question_parser.go # QNAME length-prefixed label parser
      response.go        # Builds the response to DNS request

  cache/
    cache.go             # Cache interface
    cache_repo.go        # Cache implementation with Redis implementation with JSON serialization
  blocklist/
    blocklist.go         # Blocklist interface
    blocklist_repo.go    # Blocklist implementation with SQLite implementation with migration

test/
  unit/
    answer_parser_test.go
    header_parser_test.go
    nxdomain_test.go
    question_parser_test.go
  integration/
    integration_test.go # Tests the interaction of redis sqlite and upstream resolver
data/
  blocklist.db          # SQLite database (mounted volume)
  blocked_domains.txt   # The seed to be dumped in the data base containing malicious domains

docker-compose.yml
Dockerfile
```

## Design Decisions

| Decision                  | Rationale                                                  |
| ------------------------- | ---------------------------------------------------------- |
| Built parser from scratch | Deep understanding of DNS wire format; no magic            |
| Feature-based packages    | Clear boundaries, easy to test, swap implementations       |
| Redis + SQLite            | Cache for speed, SQLite for persistence and blocklist size |

## Why Build a DNS Server?

This project explores the fundamentals of network protocols, database optimization, and concurrent systems design. Built from scratch without DNS libraries to deeply understand the wire format.
For more info I have also written a blog post [here](https://dev.to/newjhez01/regaining-privacy-by-parsing-dns-requests-bit-by-bit-462h)

## Contributing

Found a bug? Have ideas? Check out our [open issues](https://github.com/NewJhez01/gdns/issues) or feel free to open a new one. We welcome all contributions!
