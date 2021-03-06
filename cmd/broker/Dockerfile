FROM golang:1.15 as build

ARG PROJECT="akri-http"
ARG MODULE="github.com/DazWilkin/${PROJECT}"

WORKDIR /${PROJECT}

# Copy go.mod first and install dependencies
COPY go.mod .
RUN go mod download

# Copy sources
COPY . .

# Installs protoc and plugins: protoc-gen-go
ARG VERS="3.13.0"
ARG ARCH="linux-x86_64"
ARG NAME="protoc-${VERS}-${ARCH}"
RUN wget https://github.com/protocolbuffers/protobuf/releases/download/v${VERS}/${NAME}.zip --output-document=./${NAME}.zip && \
    apt update && apt install -y unzip && \
    unzip -o ${NAME}.zip -d ${NAME} && \
    mv ${NAME}/bin/* /usr/local/bin && \
    mv ${NAME}/include/* /usr/local/include && \
    go get -u github.com/golang/protobuf/protoc-gen-go

# Generates the Golang protobuf files
RUN protoc \
    --proto_path=./protos \
    --go_out=plugins=grpc,module=${MODULE}:. \
    ./protos/*.proto

# Compile Go binary
RUN GOOS=linux \
    go build -a -installsuffix cgo \
    -o /bin/broker \
    ${MODULE}/cmd/broker

# Install gRPC healthcheck
RUN GRPC_HEALTH_PROBE_VERSION=v0.3.1 && \
    wget \
    --quiet \
    --output-document=/bin/grpc_health_probe \
    https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

FROM gcr.io/distroless/base-debian10

COPY --from=build /bin/broker /
COPY --from=build /bin/grpc_health_probe /

USER 999

EXPOSE 50051

ENTRYPOINT ["/broker"]
CMD ["--grpc_endpoint=:50051"]
