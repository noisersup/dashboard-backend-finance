FROM golang:1.16.3 AS build
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o app .
CMD ["/app/app"]