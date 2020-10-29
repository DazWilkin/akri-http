# Akri: HTTP protocol

## Protoc

Requires [`protoc`](https://github.com/protocolbuffers/protobuf/releases) in the path

 ```bash
MODULE="github.com/DazWilkin/akri-http"
protoc \
 --proto_path=./protos \
 --go_out=plugins=grpc,module=${MODULE}:. \
 ./protos/http.proto
 ```

### Container Build

```bash
USER="dazwilkin" # Or your GitHub username
REPO="akri-http" # Or your preferred GHCR repo
TAGS="$(git rev-parse HEAD)"

docker build \
--tag=ghcr.io/${USER}/${REPO}:${TAGS} \
--file=./deployment/Dockerfile.devices \
.
```

## Run

### Devices

Creates an arbitrary number of pseudo-devices (that return random numbers) and a discovery service.

Either:

```bash
go run ./cmd/devices
```

Or equivalently:

```bash
DISCO="9999"
START="8000"
COUNT="10"
go run ./cmd/devices --discovery_port=${DISCO} --starting_port=${START} --num_devices=${COUNT}
```

Or:

```bash
USER="dazwilkin" # Or your GitHub username
REPO="akri-http" # Or your preferred GHCR repo
TAGS="$(git rev-parse HEAD)"

docker run \
--rm --interactive --tty \
--publish=0.0.0.0:8000-8100:8000-8100/tcp \
--publish=9999:9999 \
ghcr.io/${USER}/${REPO}:${TAGS}
```

Or the following. This is over-engineered but it shows how you may dynamically revise the discovery, starting ports and the number of devices.

```bash
USER="dazwilkin" # Or your GitHub username
REPO="akri-http" # Or your preferred GHCR repo
TAGS="$(git rev-parse HEAD)"

DISCO="9999"
START="8000"
COUNT="10"

docker run \
--rm --interactive --tty \
--publish=0.0.0.0:8000-$(echo "8000"+${COUNT}|bc):${START}-$(echo ${START}+${COUNT}|bc)/tcp \
--publish=9999:${DISCO} \
ghcr.io/${USER}/${REPO}:${TAGS} \
  --discovery_port=${DISCO} \
  --starting_port=${START} \
  --num_devices=${COUNT}
```

> **NOTE** the flags are all integers (not strings)

yields:

```console
go run ./cmd/devices
2020/10/28 13:22:28 [main] Creating Device: 0.0.0.0:8000
2020/10/28 13:22:28 [main] Register Device
2020/10/28 13:22:28 [main] Creating Device: 0.0.0.0:8001
2020/10/28 13:22:28 [main] Register Device
2020/10/28 13:22:28 [main] Creating Device: 0.0.0.0:8002
2020/10/28 13:22:28 [main] Register Device
2020/10/28 13:22:28 [main] Creating Device: 0.0.0.0:8003
2020/10/28 13:22:28 [main] Register Device
2020/10/28 13:22:28 [main] Creating Device: 0.0.0.0:8004
2020/10/28 13:22:28 [main] Register Device
2020/10/28 13:22:28 [main] Creating Device: 0.0.0.0:8005
2020/10/28 13:22:28 [main] Register Device
2020/10/28 13:22:28 [main] Creating Device: 0.0.0.0:8006
2020/10/28 13:22:28 [main] Register Device
2020/10/28 13:22:28 [main] Creating Device: 0.0.0.0:8007
2020/10/28 13:22:28 [main] Register Device
2020/10/28 13:22:28 [main] Creating Device: 0.0.0.0:8008
2020/10/28 13:22:28 [main] Register Device
2020/10/28 13:22:28 [main] Creating Device: 0.0.0.0:8009
2020/10/28 13:22:28 [main] Register Device
2020/10/28 13:22:28 [main:go] Starting Discovery Service: 0.0.0.0:9999
2020/10/28 13:22:28 [main:go] Starting Device: 0.0.0.0:8009
2020/10/28 13:22:28 [main:go] Starting Device: 0.0.0.0:8007
2020/10/28 13:22:28 [main:go] Starting Device: 0.0.0.0:8003
2020/10/28 13:22:28 [main:go] Starting Device: 0.0.0.0:8008
2020/10/28 13:22:28 [main:go] Starting Device: 0.0.0.0:8004
2020/10/28 13:22:28 [main:go] Starting Device: 0.0.0.0:8005
2020/10/28 13:22:28 [main:go] Starting Device: 0.0.0.0:8001
2020/10/28 13:22:28 [main:go] Starting Device: 0.0.0.0:8006
2020/10/28 13:22:28 [main:go] Starting Device: 0.0.0.0:8002
2020/10/28 13:22:28 [main:go] Starting Device: 0.0.0.0:8000
```

Then:

+ http://0.0.0.0:8004/ returns a random float
+ http://0.0.0.0:8004/health returns `ok`
+ http://0.0.0.0:9999/health returns `ok`
+ http://0.0.0.0:9999/ returns a list of the `0.0.0.0:8000` --> `0.0.0.0:8009` devices


### gRPC Server

```bash
go run ./cmd/server --grpc_endpoint=:50051
```

### gRPC Client

```bash
go run ./cmd/client --grpc_endpoint=:50051
```
