FROM golang:1.21-alpine

WORKDIR /app

# Копируем файлы модулей
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем исходный код
COPY *.go ./

# Собираем приложение ВНУТРИ контейнера
RUN go build -o server_1 .

# Даем права на выполнение
RUN chmod +x server_1

EXPOSE 8080

CMD ["./server_1"]