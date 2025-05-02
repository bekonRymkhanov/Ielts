FROM golang:1.23.0 as builder
WORKDIR /usr/src/app
COPY . .
WORKDIR /usr/src/app/backend
RUN go mod download

#EXPOSE 8080

CMD ["go", "run", "/usr/src/app/backend/cmd/server", "."]