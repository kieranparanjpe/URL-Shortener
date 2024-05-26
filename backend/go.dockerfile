# syntax=docker/dockerfile:1.7-labs

FROM golang:1.22.3-alpine3.20

WORKDIR /app

COPY --exclude=.env.make ./ ./

# Download and install the dependencies:
#RUN go get -d -v ./...
RUN go mod download

# Build the go app
RUN go build -v -o url-shortener .

EXPOSE 8080

CMD ["./url-shortener"]
#CMD ["echo", "hello world"]