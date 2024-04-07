FROM golang:1.20.6
WORKDIR /test
COPY . /test
RUN go build /test
EXPOSE 8000

ENTRYPOINT ["./FampayAssignment"]
