FROM golang:1.13

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 9000

CMD ["prestige-api", "--hostname", "0.0.0.0:9000"]
