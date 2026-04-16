FROM golang:1.25rc2

WORKDIR /app

RUN go install github.com/air-verse/air@latest

ENV CGO_ENABLED=0
ENV GOTOOLCHAIN=auto

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD ["air", "-c", ".air.toml"]