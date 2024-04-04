FROM golang:1.21 as builder

COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build .

FROM alpine

COPY --from=builder /app/cloudserver-vpn /app/cloudserver-vpn

CMD ["/app/cloudserver-vpn"]
