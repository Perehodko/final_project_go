FROM golang:1.23

WORKDIR /app

COPY . .
RUN go mod tidy
RUN go mod download
RUN go build -o scheduler main.go

# Переменные окружения по умолчанию
ENV TODO_PORT=8080
ENV TODO_DBFILE=/data/scheduler.db
ENV TODO_PASSWORD=my_password

ENTRYPOINT ["./scheduler"]
