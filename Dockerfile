FROM golang:alpine
RUN mkdir /app
WORKDIR /app

RUN apk add --no-cache inkscape git

ADD . .

RUN go build -o giflichess && \
      apk del git

ENTRYPOINT ["./giflichess"]
