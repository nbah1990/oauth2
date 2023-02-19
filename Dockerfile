##
## Build
FROM golang:1.19.3-buster AS build

ENV CGO_ENABLED=0
WORKDIR /app

COPY . /app
RUN go mod tidy && \
    go get -d ./... && \
    go test ./... && \
    go build -v ./cmd/oauth-proxy

##
## Deploy
FROM gcr.io/distroless/base-debian10 as deploy

WORKDIR /

COPY --from=build /app/oauth-proxy /app/oauth-proxy

EXPOSE 8092

USER nonroot:nonroot

ENTRYPOINT ["/app/oauth-proxy"]