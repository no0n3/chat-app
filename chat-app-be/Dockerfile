FROM golang:1.15.7-buster
WORKDIR /code
COPY . .
ARG CHAT_APP_ENDPOINT=http://localhost:81
ENV CHAT_APP_ENDPOINT=$CHAT_APP_ENDPOINT
RUN go mod download
RUN go build
