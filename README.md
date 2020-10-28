# Akri: HTTP protocol

## HTTP

 ```bash
 protoc \
 --proto_path=./samples/brokers/http \
 --go_out=plugins=grpc,module=github.com/DazWilkin/akri-http/protos:./golang/http/protos \
 ./samples/brokers/http/proto/http.proto
 ```

 > **NOTE** for temporary convenience the `http.proto` exists outside of the repo

## Run

### Devices

Creates an arbitrary number of pseudo-devices (that return random numbers) and a discovery service.

```bash
go run ./cmd/devices
```

Or equivalently:

```bash
go run ./cmd/devices --discovery_port=9999 --starting_port=8000 --num_devices=10
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
