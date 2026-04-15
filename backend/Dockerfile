FROM golang:1.25rc2

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.mod ./
RUN go mod download

COPY . .

CMD ["air", "-c", ".air.toml"]