FROM golang:1.25 AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY notification-service ./notification-service
RUN CGO_ENABLED=0 go build -o /out/main ./notification-service

FROM alpine:latest
RUN apk add --no-cache tzdata ca-certificates
COPY --from=build /out/main /main
COPY --from=build /src/notification-service/configs /configs
EXPOSE 8002 2123
CMD ["/main"]
