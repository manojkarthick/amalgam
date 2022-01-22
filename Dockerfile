FROM golang:1.17-alpine3.15 AS builder

RUN mkdir -pv /build
ADD . /build
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -o amalgam

FROM alpine:3.15
RUN apk add --no-cache bash
COPY --from=builder /build/amalgam .

ENTRYPOINT [ "./amalgam" ]
