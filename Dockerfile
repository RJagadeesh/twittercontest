FROM golang:1.15-alpine

ENV GO111MODULE=on
 
RUN mkdir -p /app
 
WORKDIR /app
 
COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

EXPOSE 8080

CMD ["/app/main"]
