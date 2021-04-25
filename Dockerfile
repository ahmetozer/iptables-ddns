FROM golang
WORKDIR /app/iptables-ddns
COPY . .
RUN go get -v .
RUN CGO_ENABLED=0 go build -o /app/iptables-ddns/iptables-ddns
FROM alpine
RUN apk add iptables ip6tables
ENV runningenv="container"
COPY --from=0  /app/iptables-ddns/iptables-ddns /app/iptables-ddns/iptables-ddns
COPY  iptables.list config.json /config/
LABEL org.opencontainers.image.source="https://github.com/ahmetozer/iptables-ddns"
ENTRYPOINT [ "/app/iptables-ddns/iptables-ddns" ]
