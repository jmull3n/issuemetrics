FROM golang:alpine as builder

RUN apk add git
RUN mkdir -p $GOPATH/src/github.com/jmull3n/issuemetrics
WORKDIR $GOPATH/src/github.com/jmull3n/issuemetrics
COPY . .
ENV GIT_TERMINAL_PROMPT 1
RUN go build -ldflags "-s -w" -o /bin/issuemetrics main.go

FROM alpine:edge
RUN apk add openssl-dev cyrus-sasl-dev ca-certificates && rm -rf /var/cache/apk/*
WORKDIR /
COPY --from=builder /bin/issuemetrics /bin/issuemetrics

ENTRYPOINT [ "/bin/issuemetrics"]
EXPOSE 8000 8083