FROM golang:alpine
RUN apk update && apk add --no-cache git && apk add --no-cach bash && apk add build-base
RUN mkdir /app
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o ./bin/currency-conversion-service .
EXPOSE 8080
CMD ./bin/currency-conversion-service