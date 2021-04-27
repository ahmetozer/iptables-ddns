FROM golang
WORKDIR /app/iptables-ddns
COPY . .
RUN go get -v .
RUN export GIT_COMMIT=$(git rev-list -1 HEAD) && \
    export GIT_TAG=$(git tag | tail -1) && \
    export GIT_URL=$(git config --get remote.origin.url) && \
    CGO_ENABLED=0 go build -v -ldflags="-X 'main.GitUrl=$GIT_URL' -X 'main.GitTag=$GIT_TAG' -X 'main.GitCommit=$GIT_COMMIT' -X 'main.BuildTime=$(date -Isecond)' -X 'main.RunningEnv=container'" -o /app/iptables-ddns/iptables-ddns
FROM alpine
RUN apk add iptables ip6tables
COPY --from=0  /app/iptables-ddns/iptables-ddns /app/iptables-ddns/iptables-ddns
COPY  iptables.list config.json /config/
LABEL org.opencontainers.image.source="https://github.com/ahmetozer/iptables-ddns"
ENTRYPOINT [ "/app/iptables-ddns/iptables-ddns" ]
