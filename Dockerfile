FROM golang:latest AS builder
LABEL stage=gobuilder
ENV CGO_ENABLED 0
ENV GOOS linux
WORKDIR /build
COPY ./ /build
RUN go mod download
RUN go build -ldflags="-s -w" -o analytic /build/cmd/analytic/main.go

FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /build/analytic .
COPY --from=builder /build/config.yaml .
CMD ["./analytic migrate up", "./analytic"]