FROM golang:1.25 AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY api-gateway ./api-gateway
RUN CGO_ENABLED=0 go build -o /out/main ./api-gateway

FROM alpine:latest
RUN apk add --no-cache tzdata ca-certificates
COPY --from=build /out/main /main
COPY --from=build /src/api-gateway/configs /configs
EXPOSE 8000 2121
CMD ["/main"]
