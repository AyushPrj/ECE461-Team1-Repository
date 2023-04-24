# syntax=docker/dockerfile:1

FROM golang:1.18

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod tidy && go mod download

COPY . ./


#run install
ARG MONGOURI
ENV MONGOURI $MONGOURI

ENV DANGEROUSLY_DISABLE_HOST_CHECK=true

RUN apt-get update && apt-get upgrade -y && apt-get install -y nodejs npm   
RUN npm install

EXPOSE 8080
EXPOSE 3000

#run main
CMD HOME=/root go run main/main.go 
