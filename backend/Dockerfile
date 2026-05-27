FROM golang:1.22-alpine AS build

WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /app/api ./cmd/api

FROM alpine:3.20
RUN adduser -D -H app
USER app
COPY --from=build /app/api /api
EXPOSE 8080
ENTRYPOINT ["/api"]
