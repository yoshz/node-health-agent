FROM golang:alpine as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN go build

FROM alpine
RUN adduser -S -D -H agent
COPY --from=builder /build/node-health-agent /usr/local/bin
USER agent
CMD ["node-health-agent", "--incluster"]
