FROM golang:1.13 as builder
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o httpmultiplexor ./cmd

FROM alpine
RUN apk add ca-certificates
COPY --from=builder /app/httpmultiplexor /usr/bin/httpmultiplexor
ENTRYPOINT [ "/usr/bin/httpmultiplexor" ]