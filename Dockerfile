# stage for development, which contains tools for code generation and debugging.
FROM golang:1.18-bullseye as dev

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        wget \
        make \
        unzip \
        git \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# install protoc
RUN wget https://github.com/protocolbuffers/protobuf/releases/download/v3.20.1/protoc-3.20.1-linux-x86_64.zip \
    && unzip -d /usr/local protoc-3.20.1-linux-x86_64.zip \
    && rm protoc-3.20.1-linux-x86_64.zip

# install evans, which is a command line tool that works as a grpc client.
#   https://github.com/ktr0731/evans
RUN wget https://github.com/ktr0731/evans/releases/download/v0.10.6/evans_linux_amd64.tar.gz \
    && tar xzf evans_linux_amd64.tar.gz \
    && mv evans /usr/local/bin \
    && rm evans_linux_amd64.tar.gz

# install protoc plugins
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.10.3
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.10.3

WORKDIR /work

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . ./

# builder stage
FROM golang:1.18-bullseye as builder

WORKDIR /work

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . ./

RUN go build -o /server ./cmd/server
RUN go build -o /gateway ./cmd/gateway


# runner stage for server
FROM debian:bullseye-slim as server-runner

COPY --from=builder /server /server

CMD ["/server"]

EXPOSE 50051



# runner stage for gateway
FROM debian:bullseye-slim as gateway-runner

COPY --from=builder /gateway /gateway
CMD ["/gateway"]

EXPOSE 8080
