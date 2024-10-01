FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /tools cmd/web/main.go

# Run the tests in the container
FROM builder AS tester
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM alpine:latest AS releaser

WORKDIR /

COPY --from=builder /tools /tools

EXPOSE 8080

RUN addgroup -S nonroot && adduser -S nonroot -G nonroot

USER nonroot:nonroot

ENTRYPOINT ["/tools"]
