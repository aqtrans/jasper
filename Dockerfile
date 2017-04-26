FROM aqtrans/golang-npm:latest

RUN mkdir -p /go/src/jasper
WORKDIR /go/src/jasper

ADD . /go/src/jasper/
RUN go get github.com/kardianos/govendor && govendor sync
RUN go get -d
RUN go build -o ./jasper

# Expose the application on port 3000
#EXPOSE 3000

# Set the entry point of the container to the bee command that runs the
# application and watches for changes
CMD ["./jasper"]
