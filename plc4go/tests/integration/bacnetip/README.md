# BACnet/IP Integration Tests

End-to-end integration tests for the plc4go BACnet/IP driver against a
dockerized [bacpypes3](https://github.com/JoelBender/BACpypes3) virtual device.

## Why this exists

The unit tests in `internal/bacnetip/*_test.go` exercise the pure decoding /
encoding paths. They do not catch:

- Real UDP socket behavior (port reuse, multicast, broadcast).
- Live transaction-manager retry loops.
- Real-device quirks (segmentation window negotiation, IAm timing).

This suite spins a real BACnet device in a container and drives the plc4go
driver against it. It is **opt-in**: the default `make test` does not run it.

## How to run

```sh
# From plc4go/ directory:
docker compose -f tests/integration/bacnetip/docker-compose.yml up -d
BACNET_IT=1 go test -tags integration ./tests/integration/bacnetip/... -v -count=1
docker compose -f tests/integration/bacnetip/docker-compose.yml down
```

Or, once the Makefile target lands:

```sh
make integration-bacnetip
```

## What's in here

| File                  | Purpose                                                          |
|-----------------------|------------------------------------------------------------------|
| `Dockerfile`          | Builds a `python:3.12-slim` image with `bacpypes3` preinstalled. |
| `device.py`           | A bacpypes3 LocalDeviceObject with AV.0–4, BV.0–1, AI.0, MSV.0.  |
| `docker-compose.yml`  | Exposes UDP 47808 on the host.                                   |
| `integration_test.go` | Test cases (Discover, Read, Write, COV, segmented response).     |

## Caveats

- UDP 47808 must be free on the host. If you already run BACnet locally
  (Niagara, Cimplicity, etc.), pick a different host port and override
  `BACNET_PORT` in the test command.
- `docker compose up -d` returns before bacpypes3 is fully listening. The
  test wraps the first WhoIs in a short retry loop to absorb that.
- bacpypes3 simulates a single device; the multi-device routing tests in
  `internal/bacnetip/DeviceInfoCache_test.go` cover that path via unit tests.

## Scope

What's covered:

1. `Discover()` returns device `1234` within the configured timeout.
2. `Read AnalogValue.0/PRESENT_VALUE` round-trips as a `PlcREAL`.
3. `Write AnalogValue.0/PRESENT_VALUE = 42.5` followed by a re-read.
4. `WritePropertyMultiple` of 3 properties succeeds.
5. `Subscribe AnalogInput.0` receives ≥1 COV notification within 5s.
6. Large `ReadPropertyMultiple` (10 properties) exercises 2-segment
   reassembly.

What's not covered (deferred to v2):

- BACnet/SC (secure-connect over WebSocket).
- BACnet/IPv6 (Annex U).
- MS/TP routing.
- Real-device interop with quirky stacks (Tracer, older Honeywell).
