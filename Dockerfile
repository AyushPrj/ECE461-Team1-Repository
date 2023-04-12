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
CMD HOME=/root go run main/main.go & cd assets && npm install && npm start

# docker run --env-file=.env alpine env
# docker build --tag webservice --build-arg 
# docker run --publish 5500:5500 webservice
# docker run --publish 0.0.0.0:3000:3000 webservice
# docker run --publish 0.0.0.0:3000:3000 --publish 0.0.0.0:8080:8080 webservice
# docker run --publish 35.209.87.90:3000:3000 webservice
