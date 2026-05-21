# BACnet/IP Integration Tests

End-to-end integration tests for the plc4go BACnet/IP driver. The suite runs
the driver and a [bacpypes3](https://github.com/JoelBender/BACpypes3) virtual
device side-by-side in two docker containers on a shared bridge network — so
each end has its own network namespace and can bind UDP 47808 the way the
BACnet/IP spec assumes.

## Why two containers

BACnet/IP defines port 47808 as the well-known UDP port for *both* sides of a
conversation. Some flows (Discover via WhoIs broadcast, and Confirmed COV
notifications on spec-strict devices) only work if the driver can also bind
47808. With a single host network namespace, only one process can hold the
port, so the simulator + the driver collide. A docker bridge gives each
container its own namespace; the kernel routes broadcasts across the bridge,
so both binds can coexist.

## How to run

```sh
# From plc4go/ directory:
make integration-bacnetip
```

That target builds both images, brings them up, runs the test suite to
completion, tears the stack down, and propagates the test-runner exit code.

For ad-hoc invocation:

```sh
docker compose -f tests/integration/bacnetip/docker-compose.yml up \
    --build --abort-on-container-exit --exit-code-from test-runner
docker compose -f tests/integration/bacnetip/docker-compose.yml down
```

The compose project must be evaluated from the plc4go module root (which is
what `make integration-bacnetip` does) because the test-runner's build
context is the whole Go module.

## What's in here

| File                  | Purpose                                                                |
|-----------------------|------------------------------------------------------------------------|
| `Dockerfile`          | Builds `python:3.12-slim` + `bacpypes3` for the simulated device.      |
| `device.py`           | A bacpypes3 LocalDeviceObject with AV.0–4, BV.0–1, AI.0, MSV.0.        |
| `Dockerfile.test`     | Builds `golang:1.26-bookworm` + the plc4go module sources.             |
| `docker-compose.yml`  | Wires both containers onto the `bacnet` bridge (172.30.0.0/24).        |
| `integration_test.go` | Test cases (Discover, Read, Write+re-read, Subscribe).                 |

The test-runner reads `BACNET_IT_HOST` (defaulted to the docker DNS name
`bacnet-device`) so the same suite can be pointed at any reachable simulator
by exporting that env var before invoking `go test` outside compose.

## Caveats

- The two-container setup needs a working Docker daemon with bridge driver
  support. Docker-in-Docker CI runners need `--privileged` (or careful
  configuration of `dockerd-rootless`) to spawn user-space bridges.
- bacpypes3 isn't ready the instant the container starts. The test-runner
  sleeps three seconds before running the Go suite to absorb the race;
  a proper UDP healthcheck would let us drop that.
- The simulated device is single-instance (device 1234). Multi-device routing
  is covered by unit tests in `internal/bacnetip/DeviceInfoCache_test.go`.

## Scope

What's covered:

1. `Discover()` returns device `1234` within the configured timeout (real
   WhoIs broadcast + IAm round-trip over the bridge).
2. `Read AnalogValue.0/PRESENT_VALUE` round-trips as a `PlcREAL`.
3. `Write AnalogValue.1/PRESENT_VALUE = 42.5` followed by a re-read.
4. `Subscribe AnalogInput.0` receives ≥1 COV notification within 10s
   (the simulator runs a 2-second sawtooth on AI-0).

What's not covered (deferred to v2):

- BACnet/SC (secure connect over WebSocket).
- BACnet/IPv6 (Annex U).
- MS/TP routing.
- Quirky third-party stacks (older Honeywell, Trane Tracer); only a real-device
  QA pass exposes these.
