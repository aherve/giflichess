FROM golang:1.16.5-stretch
RUN mkdir /app
WORKDIR /app

RUN apt update && apt install inkscape imagemagick git -y

ADD . .

RUN go build -o giflichess

EXPOSE 8080
ENTRYPOINT ["./giflichess"]
