FROM golang:1.21-alpine

RUN mkdir -p /app

WORKDIR /app

COPY . .

RUN apk update \
    && apk add --no-cache git

ARG GIT_TOKEN
ENV GIT_TOKEN=${GIT_TOKEN}

# Install project dependencies
RUN go env -w GOPRIVATE=github.com/YOLOTECHNOLOGY \
    && go env -w GONOSUMDB=github.com/YOLOTECHNOLOGY && go env -w GONOPROXY=github.com/ \
    && git config --global url."https://${GIT_TOKEN}@github.com".insteadOf "https://github.com/" \
    && go mod tidy && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./main.go

EXPOSE 8000

CMD ["./main"]