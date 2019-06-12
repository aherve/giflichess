FROM golang:alpine

RUN apk add --no-cache inkscape git

ADD . .

#RUN go get -u github.com/notnil/chessimg github.com/notnil/chess && \
RUN go get -u \
      github.com/notnil/chessimg \
      github.com/notnil/chess \
      github.com/urfave/cli && \
      go build -o giflichess && \
      apk del git

ENTRYPOINT ["./giflichess"]
