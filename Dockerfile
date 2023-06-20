FROM golang:latest

WORKDIR /app

COPY . .

RUN go build -o app

ENV file_path=db.csv

CMD ["./app"]