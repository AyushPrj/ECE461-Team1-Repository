# syntax=docker/dockerfile:1

FROM golang:1.18

WORKDIR /app

#copy files
COPY . ./
RUN go mod tidy && go mod download

#build
WORKDIR /main/
RUN go build -o main main.go

#set env vars
ARG MONGOURI
ENV MONGOURI $MONGOURI
ENV DANGEROUSLY_DISABLE_HOST_CHECK=true

# RUN apt-get update && apt-get upgrade -y && apt-get install -y nodejs npm   
# RUN npm install

EXPOSE 8080
# EXPOSE 3000

#run main
# CMD HOME=/root go run main/main.go 
CMD ["./main"]

