FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-w -s -X main.version=${VERSION}" -o server .

FROM alpine:3.20
WORKDIR /app
RUN adduser -D -u 1000 appuser && apk --no-cache add ca-certificates
COPY --from=builder /app/server /app/server
RUN chown -R appuser:appuser /app
USER appuser
EXPOSE 8080
CMD ["/app/server"]