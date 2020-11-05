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
    -o /bin/device \
    ${MODULE}/cmd/device


FROM gcr.io/distroless/base-debian10

COPY --from=build /bin/device /

USER 999



ENTRYPOINT ["/device"]
CMD ["--path=/","--path=/sensor"]
