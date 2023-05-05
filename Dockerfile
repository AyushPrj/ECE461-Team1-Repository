# syntax=docker/dockerfile:1

# FROM golang:1.18

# WORKDIR /app

# ENV PORT 8080
# ENV HOST 0.0.0.0

# #copy files
# COPY . ./
# RUN go mod tidy && go mod download

# #build
# RUN go build -o maine main.go

# EXPOSE 8080
# # EXPOSE 3000

# #run main
# # CMD HOME=/root go run main/main.go 
# CMD ["./maine"]

# Stage 1: Build the Go application
FROM golang:1.18 AS builder

# Set the working directory
WORKDIR /app

# Copy your Go application's source code
COPY . .

ENV LOG_LEVEL="2"
ENV LOG_FILE="logfile.log"

# Build your Go application
RUN go build -o main . || (echo "Build failed" && exit 1)

# Stage 2: Set up the final image with Go application, Perl, and cloc
FROM perl:5.34

# Set the working directory
WORKDIR /app

# Copy the built Go application from the builder stage
COPY --from=builder /app/main /app/main

# # Copy the built build folder
# COPY ./build /app/build

# Install cloc
RUN curl -L -o cloc https://github.com/AlDanial/cloc/releases/download/v1.90/cloc-1.90.pl \
    && chmod +x cloc \
    && mv cloc /usr/local/bin

# Expose the port your application listens on
EXPOSE 8080

# Start your Go application
CMD ["/app/main"]