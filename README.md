# Akri: HTTP protocol

![containers-devices](https://github.com/DazWilkin/akri-http/workflows/containers-devices/badge.svg?branch=master)
![containers-device-discovery](https://github.com/DazWilkin/akri-http/workflows/containers-device-discovery/badge.svg)
![containers-grpc](https://github.com/DazWilkin/akri-http/workflows/containers-grpc/badge.svg)


## v2

The former `devices` service has been split into two:

+ `device`
+ `discovery`

This overcomes an unanticipated problem with `devices` in that its discovery endpoint was bound to a single host address at runtime.

The new services can be run:

### Golang

```bash
go run ./cmd/device \
--path="/sensor1" \
--path="/sensor2" \
...
```

And:

```bash
go run ./cmd/discovery \
--device="device:8000" \
--device="device:8001" \
...
```

### Docker

Which is more useful when running under Docker:

```bash
USER="..."
TAG=$(git rev-parse HEAD)

# Create devices on ports 8000:8009
DISCOVERY=()
for PORT in {8000..8009}
do
  # Create the device on ${PORT}
  # For Docker only: name each device: device-${PORT}
  docker run \
  --rm --detach=true \
  --name=device-${PORT} \
  --publish=${PORT}:8080 \
  ghcr.io/${USER}/akri-http-device:${TAG}
    --path="/sensor1" \
    --path="/sensor2"
  # Add the device to the discovery document
  DISCOVERY+=("--device=device:${PORT} ")
done

# Create a discovery server for these devices
docker run \
  --rm --detach=true \
  --name=discovery \
  --publish=9999:9999 \
  ghcr.io/${USER}/akri-http-discovery:${TAG} ${DISCOVERY[@]}
```

Test:

```bash
curl localhost:9999/
device:8000
device:8001
device:8002
device:8003
device:8004
device:8005
device:8006
device:8007
device:8008
device:8009

curl localhost:8006/sensor
```

To stop:

```bash
USER="..."
TAG=$(git rev-parse HEAD)

# Delete devices on ports 8000:8009
for PORT in {8000..8009}
do
  docker stop  device-${PORT}
done

# Delete discovery server
docker stop discovery
```

### Kubernetes

And most useful on Kubernetes because one (!) or more devices can be created and then discovery can be created with correct DNS names.

Ensure the `image` references are updated in `./kubernetes/v2/device.yaml` and `./kubernetes/v2/discovery.yaml`

Then:

```bash

# Create one device to rule them all
kubectl apply --filename=device.yaml
deployment.apps/device created

# But multiple services against the singular device
for NUM in {1..9}
do
  # Services are uniquely named
  # The service uses the Pods port: 8080
  kubectl expose deployment/device \
  --name=device-${NUM} \
  --port=8080 \
  --target-port=8080
done
service/device-1 exposed
service/device-2 exposed
service/device-3 exposed
service/device-4 exposed
service/device-5 exposed
service/device-6 exposed
service/device-7 exposed
service/device-8 exposed
service/device-9 exposed

# Expose Discovery as a service on its default port: 9999
# The Discovery service spec is statically configured for devices 1-9
kubectl expose deployment/discovery \
--name=discovery \
--port=9999 \
--target-port=9999

kubectl run curl --image=radial/busyboxplus:curl --stdin --tty --rm
curl discovery:9999
device-1:8080
device-2:8080
device-3:8080
device-4:8080
device-5:8080
device-6:8080
device-7:8080
device-8:8080
device-9:8080
```

Delete:

```bash
kubectl delete deployment/discovery
kubectl delete deployment/device

kubectl delete service/discovery

for NUM in {1..9}
do
  kubectl delete service/device-${NUM}
done
```

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
