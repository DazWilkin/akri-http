# Akri: HTTP protocol

![devices](https://github.com/DazWilkin/akri-http/workflows/containers-devices/badge.svg?branch=master)
![device-discovery](https://github.com/DazWilkin/akri-http/workflows/containers-device-discovery/badge.svg)
![grpc](https://github.com/DazWilkin/akri-http/workflows/containers-grpc/badge.svg)

[![Go Report Card](https://goreportcard.com/badge/github.com/DazWilkin/akri-http)](https://goreportcard.com/report/github.com/DazWilkin/akri-http)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/DazWilkin/akri-http)](https://pkg.go.dev/github.com/DazWilkin/akri-http)


## Device|Discovery Services

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
  ghcr.io/${USER}/akri-http-device:${TAG} \
    --path="/sensor1" \
    --path="/sensor2" \
  # Add the device to the discovery document
  DISCOVERY+=("--device=http://localhost:${PORT} ")
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
curl http://localhost:9999/
http://localhost:8000
http://localhost:8001
http://localhost:8002
http://localhost:8003
http://localhost:8004
http://localhost:8005
http://localhost:8006
http://localhost:8007
http://localhost:8008
http://localhost:8009

curl http://localhost:8006/sensor
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

# Create one device deployment
kubectl apply --filename=./device.yaml

# But multiple Services against the single Pod
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

# Create one discovery deployment
kubectl apply --filename=./discovery.yaml

# Expose Discovery as a service on its default port: 9999
# The Discovery service spec is statically configured for devices 1-9
kubectl expose deployment/discovery \
--name=discovery \
--port=9999 \
--target-port=9999

kubectl run curl --image=radial/busyboxplus:curl --stdin --tty --rm
curl http://discovery:9999
http://device-1:8080
http://device-2:8080
http://device-3:8080
http://device-4:8080
http://device-5:8080
http://device-6:8080
http://device-7:8080
http://device-8:8080
http://device-9:8080
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


## gRPC Broker|Client

### Protoc

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
--file=./deployment/Dockerfile \
.
```

### Golang

Then run the gRPC Broker:

```bash
PORT=50051
go run ./cmd/broker --grpc_endpoint=:${PORT}
```

The run the gRPC Client

```bash
go run ./cmd/client --grpc_endpoint=:${PORT}
```


These are containerized too:

```bash
docker run \
--rm --interactive --tty \
--net=host \
--name=grpc-broker-golang \
--env=AKRI_HTTP_DEVICE_ENDPOINT=localhost:8005 \
ghcr.io/dazwilkin/akri-http-grpc-broker-golang:latest \
--grpc_endpoint=:50051
```

And:

```bash
docker run \
--rm --interactive --tty \
--net=host \
--name=grpc-client-golang \
ghcr.io/dazwilkin/akri-http-grpc-client-golang:latest \
--grpc_endpoint=:50051
```


