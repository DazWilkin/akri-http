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
    -o /bin/discovery \
    ${MODULE}/cmd/discovery


FROM gcr.io/distroless/base-debian10

COPY --from=build /bin/discovery /

USER 999
EXPOSE 9999


ENTRYPOINT ["/discovery"]
CMD ["--device=device:8000","--device=device:8001"]
