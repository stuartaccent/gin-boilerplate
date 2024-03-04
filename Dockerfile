FROM golang:bookworm AS builder
WORKDIR /app
COPY . .
ENV CGO_ENABLED=0
RUN go mod tidy
RUN go build -o . ./...

FROM alpine
WORKDIR /app
COPY --from=builder /app/cli /app/cli
COPY --from=builder /app/server /app/server
COPY --from=builder /app/config.toml /app/config.toml
EXPOSE 80
ENTRYPOINT ["/app/server"]