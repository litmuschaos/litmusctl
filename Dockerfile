# auto-build: {"image_name": "chaos-as-service/litmusctl-sumo, "tags": ["latest"], "platform": "linux/amd64,linux/arm64"}

FROM golang:1.23.4 AS builder
WORKDIR /litmusctl
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o litmusctl main.go
RUN chmod +x litmusctl
RUN cp litmusctl /usr/local/bin

FROM alpine:latest
WORKDIR /litmusctl
COPY --from=builder /litmusctl/ ./
RUN cp ./litmusctl /usr/local/bin/
