FROM golang:alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o /transactions ./cmd/app

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /transactions /transactions
USER nonroot
EXPOSE 8080
ENTRYPOINT ["/transactions"]
