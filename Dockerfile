FROM golang:1.20.6
WORKDIR /test
COPY . /test
RUN go build /test
EXPOSE 8000

ENV DB_HOST=localhost
ENV DB_PORT=5432
ENV DB_USER=postgres
ENV DB_PASSWORD=postgres
ENV DB_NAME=Fampay
ENV  YOUTUBE_API_KEY=AIzaSyBXzjV6FlnSxbBUKKsImVzzlGvL13aWLPA


ENTRYPOINT ["./FampayAssignment"]
