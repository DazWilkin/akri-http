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
    -o /bin/client \
    ${MODULE}/cmd/client


FROM gcr.io/distroless/base-debian10

COPY --from=build /bin/client /

USER 999

EXPOSE 50051

ENTRYPOINT ["/client"]
CMD ["--grpc_endpoint=:50051"]
