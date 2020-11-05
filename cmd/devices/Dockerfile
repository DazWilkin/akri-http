FROM golang:1.15 as build

ARG PROJECT="akri-http"
ARG MODULE="github.com/DazWilkin/${PROJECT}"

WORKDIR /${PROJECT}

# Copy go.mod first and install dependencies
COPY go.mod .
RUN go mod download

# Copy all sources
COPY . .


















# Compile Go binary
RUN GOOS=linux \
    go build -a -installsuffix cgo \
    -o /bin/devices \
    ${MODULE}/cmd/devices


FROM gcr.io/distroless/base-debian10

COPY --from=build /bin/devices /

USER 999



ENTRYPOINT ["/devices"]
CMD ["--discovery_port=9999","--starting_port=8000","--num_devices=10"]
