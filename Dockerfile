FROM golang:bookworm AS builder
WORKDIR /app
COPY . .
ENV CGO_ENABLED=0
RUN go mod tidy
RUN go build -o app .

FROM alpine
WORKDIR /app
COPY --from=builder /app/app /app/app
COPY --from=builder /app/config.toml /app/config.toml
EXPOSE 80
ENTRYPOINT ["/app/app"]
CMD ["server", "--config", "config.toml"]