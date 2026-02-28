FROM golang:1.26-alpine AS build
RUN apk add --no-cache ca-certificates
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app ./cmd/app

FROM scratch
COPY config/config.docker.yml /app/config/
COPY --from=build /app /app/app
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 50052
CMD ["/app/app", "-config", "/app/config/config.docker.yml"]
