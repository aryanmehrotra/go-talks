FROM golang:1.25 AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY order-service ./order-service
RUN CGO_ENABLED=0 go build -o /out/main ./order-service

FROM alpine:latest
RUN apk add --no-cache tzdata ca-certificates
COPY --from=build /out/main /main
COPY --from=build /src/order-service/configs /configs
EXPOSE 8001 2122
CMD ["/main"]
