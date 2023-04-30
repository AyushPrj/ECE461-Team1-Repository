# syntax=docker/dockerfile:1

FROM golang:1.18

WORKDIR /app

ENV PORT 8080
ENV HOST 0.0.0.0

#copy files
COPY . ./
RUN go mod tidy && go mod download

#build
RUN go build -o maine main.go

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
CMD ["./maine"]

