FROM golang:latest

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /mediabalancer cmd/balancer/main.go

EXPOSE 80
CMD [ "/mediabalancer", "-d", "database:6379" ]