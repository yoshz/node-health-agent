FROM golang:1.27-alpine AS builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN go build

FROM alpine:3.21
RUN adduser -S -D -H agent
COPY --from=builder /build/node-health-agent /usr/local/bin
USER agent
CMD ["node-health-agent", "--incluster"]
